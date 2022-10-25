package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"

	"rin/sishop/controllers"

)

func main(){
	store:= session.New()

	engine:=html.New("./views", ".html")
	
	app:= fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/public", "./public")

	authController := controllers.InitAuthController(store)
	prodController:=controllers.InitProductController()
	transController:=controllers.InitTransactionController(store)

	app.Get("/login", authController.Login)
	app.Post("/login", authController.LoginPosted)
	app.Get("/register", authController.Register)
	app.Post("/register", authController.RegisterPosted)
	app.Get("/profile",authController.Profile)
	app.Get("/logout", authController.Logout)

	prod := app.Group("/products")
	prod.Get("/", prodController.IndexProduct)//semua produk
	prod.Get("/create", prodController.AddProduct)
	prod.Post("/create", prodController.AddProductPosted)
	prod.Get("/productdetail/", prodController.DetailProduct)
	// prod.Get("/detail/:id", prodController.DetailProduct2)
	prod.Get("/editproduct/:id", prodController.EditProduct)
	prod.Post("/editproduct/:id", prodController.EditProductPosted)
	prod.Get("/deleteproduct/:id", prodController.DeleteProduct)

	app.Get("/profile", func(c *fiber.Ctx)error{
		sess,_:=store.Get(c)
		val:=sess.Get("username")
		if val != nil {
			return c.Next()
		} 
		return c.Redirect("./login")
	}, authController.Profile)

	trans := app.Group("/transactions")
	trans.Get("/", func(c *fiber.Ctx) error {
		sess, _ := store.Get(c)
		val := sess.Get("username")
		if val != nil {
			return c.Next()
		}

		return c.Redirect("/login")

	}, transController.IndexTransaction)

	trans.Post("/create", func(c *fiber.Ctx) error {
		sess, _ := store.Get(c)
		val := sess.Get("username")
		if val != nil {
			return c.Next()
		}

		return c.Redirect("/login")

	}, transController.AddTransactionPosted)

	trans.Get("/delete/:id", func(c *fiber.Ctx) error {
		sess, _ := store.Get(c)
		val := sess.Get("username")
		if val != nil {
			return c.Next()
		}

		return c.Redirect("/login")

	}, transController.DeleteTransactionById)

	// trans.Post("/bayar/:id", func(c *fiber.Ctx) error {
	// 	sess, _ := store.Get(c)
	// 	val := sess.Get("username")
	// 	if val != nil {
	// 		return c.Next()
	// 	}

	// 	return c.Redirect("/login")

	// }, transController.BayarTransactionById)


app.Listen(":3000")

}