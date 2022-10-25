package controllers

import (
	"strconv"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm" 

	"rin/sishop/database"
	"rin/sishop/models"
)

type ProductController struct{
	Db *gorm.DB
}

func InitProductController() *ProductController {
	db := database.InitDb()
	db.AutoMigrate(&models.Product{})

	return &ProductController{Db: db}
}

//GET product index
func (controller *ProductController) IndexProduct(c *fiber.Ctx) error {
	// var products []models.Product
	// err := models.ReadProducts(controller.Db, &products)

		// var products = []*ProductController{
		// {Id:1, Name: "Product 1", Price: 2.5},
		// {Id:2, Name: "Product 2", Price: 2.5},
		// {Id:3, Name: "Product 3", Price: 2.5},
		// {Id:4, Name: "Product 4", Price: 2.5},
		// {Id:5, Name: "Product 5", Price: 2.5},
	// } 
	// if err != nil {
	// 	return c.SendStatus(500)
	// }

	var products []models.Product
	err := models.ReadProducts(controller.Db, &products)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	return c.Render("products", fiber.Map{
		"Title": "SiShop",
		"Content": "Daftar Produk",
		"Products": products,
	})
}

//GET Product/create
func (controller *ProductController) AddProduct(c *fiber.Ctx) error {
	return c.Render("addproduct", fiber.Map{
		"Title": "SiShop",
		"Content": "Tambah Produk",
	})
}

// POST products/create
func (controller *ProductController) AddProductPosted(c *fiber.Ctx) error {
	var myform models.Product

	file, erfile := c.FormFile("image")
	if erfile != nil {
		fmt.Println("Error File =", erfile)
	}
	var filename string = file.Filename
	if file != nil {

		ersavefile := c.SaveFile(file, fmt.Sprintf("./public/images/%s", filename))
		if ersavefile != nil {
			fmt.Println("gagal menyimpan gambar.")
		}
	} else {
		fmt.Println("Nothing file to uploading.")
	}

	if err := c.BodyParser(&myform); err != nil {
		return c.Redirect("/products")
	}

	myform.Image = filename
	// save product
	errr := models.CreateProduct(controller.Db, &myform)
	if errr != nil {
		return c.Redirect("/products")
	}
	// if succeed
	return c.Redirect("/products")
}

//GET detailproduct
func (controller *ProductController) DetailProduct(c *fiber.Ctx) error {
	id := c.Query("id")
	idn, _ := strconv.Atoi(id)

	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	return c.Render("productdetail", fiber.Map{
		"Title": "SiShop",
		"Content": "Detail Produk",
		"Product": product,

	})
}

//GET products/deleteproduct/xx
func (controller *ProductController) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn,_ := strconv.Atoi(id)

	var product models.Product
	models.DeleteProductById(controller.Db, &product, idn)
	return c.Redirect("/products")	
}

func (controller *ProductController) EditProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)
	// fmt.Print("a")

	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	// fmt.Print("b")

	if err != nil {
		// fmt.Print("c")
		return c.SendStatus(500) // http 500 internal server error
	}

	fmt.Print("d")
	return c.Render("editproduct", fiber.Map{
		"Title": "SiShop",
		"Content": "Edit Produk",
		"Product": product,
	})
}

func (controller *ProductController) EditProductPosted(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	var myform models.Product

	if err := c.BodyParser(&myform); err != nil {
		return c.Redirect("/products")
	}

	file, errFile := c.FormFile("image")
	if errFile != nil {
		fmt.Println("Error File =", errFile)
	}
	var filename string = file.Filename
	if file != nil {

		errSaveFile := c.SaveFile(file, fmt.Sprintf("./public/images/%s", filename))
		if errSaveFile != nil {
			fmt.Println("Fail to store file into public/ikmages directory.")
		}
	} else {
		fmt.Println("Nothing file to uploading.")
	}
	myform.Image = filename
	product.Name = myform.Name
	product.Image = myform.Image
	product.Price = myform.Price
	// save product
	models.UpdateProduct(controller.Db, &product)

	return c.Redirect("/products")

}