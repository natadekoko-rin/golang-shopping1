package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"rin/sishop/database"
	"rin/sishop/models"
	"fmt"
)

type LoginController struct{
	Db *gorm.DB
	store *session.Store
}

type LoginForm struct{//buat login dan regis
	Username string `form:"username" validation:"required"`
	Password string `form:"password" validation:"required"`
}

func InitAuthController(s *session.Store) *LoginController{
	db := database.InitDb()

	db.AutoMigrate(&models.User{})
	return &LoginController{Db: db, store: s}
}

//GET Login
// func (controller *LoginController) Login(c *fiber.Ctx) error {
// 	return c.Render("login", fiber.Map{
// 		"Title": "SiShop",
// 		"Content": "Login",
// 	})
// }

//POST Login
func (controller *LoginController) LoginPosted(c *fiber.Ctx) error {

	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}

	var user models.User
	var myform LoginForm
	if err := c.BodyParser(&myform); err != nil {
		
		fmt.Println(err)
		return c.JSON(fiber.Map{
			"status":  400,
			"message": "Bad Request, Pastikan form terisi komplit",
		})
	}

	er := models.FindByUsername(controller.Db, &user, myform.Username)
	if er != nil {
			return c.JSON(fiber.Map{
				"message": "User tidak ditemukan",
			})
	}

	// hardcode auth
	mycompare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(myform.Password))
	if mycompare != nil {
		sess.Set("username", user.Username)
		sess.Set("userID", user.Id)
		sess.Save()

		return c.JSON(fiber.Map{
			"message": "Login sukses",
		})
	}
	return c.JSON(fiber.Map{
		"status":  401,
		"message": "Register dahulu",
	})

}

//GET Register
// func (controller *LoginController) Register(c *fiber.Ctx) error {
// 	return c.Render("register", fiber.Map{
// 		"Title": "Register",
// 	})
// }

//POST Register
func (controller *LoginController) RegisterPosted(c *fiber.Ctx) error {
	var myform models.User
	var encpass LoginForm

	if err := c.BodyParser(&myform); err != nil {
		return c.JSON(fiber.Map{
			"status":  400,
			"message": "Bad Request, Pastikan data terisi komplit",
		})
	}

	// fmt.Println(myform)
	// fmt.Println(encpass.Password)

	encpassword, _ := bcrypt.GenerateFromPassword([]byte(encpass.Password), 10)
	xHash := string(encpassword)
	// fmt.Println(xHash)
	myform.Password = xHash
	// fmt.Println(myform)

	err := models.CreateUser(controller.Db, &myform)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Gagal menyimpan user",
		})
	}
	// if succeed
	return c.JSON(fiber.Map{
		"message": "User berhasil dibuat",
	})
}

//Profile
func (controller *LoginController) Profile (c *fiber.Ctx)error {
	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	val := sess.Get("username")

	return c.JSON(fiber.Map{
		"username": val,
	})
}

//Logout 
func (controller *LoginController) Logout (c *fiber.Ctx)error {
	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	sess.Destroy()

	return c.JSON(fiber.Map{
		"message": "Berhasil Logout",
	})
}