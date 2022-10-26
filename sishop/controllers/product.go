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

func (controller *ProductController) IndexProduct(c *fiber.Ctx) error {
	var products []models.Product
	err := models.ReadProducts(controller.Db, &products)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	return c.JSON(fiber.Map{
		"Message":  "Berhasil mendapatkan seluruh list products",
		"Products": products,
	})
}

//GET Product/create
// func (controller *ProductController) AddProduct(c *fiber.Ctx) error {
// 	return c.Render("addproduct", fiber.Map{
// 		"Title": "SiShop",
// 		"Content": "Tambah Produk",
// 	})
// }

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
		return c.JSON(fiber.Map{
			"status":  400,
			"message": "Bad Request, Product Form is not complete",
		})
	}

	myform.Image = filename
	// save product
	errr := models.CreateProduct(controller.Db, &myform)
	if errr != nil {
		return c.SendStatus(500)
	}
	// if succeed
	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Berhasil Menambahkan Product",
	})
}


func (controller *ProductController) DetailProduct(c *fiber.Ctx) error {
	id := c.Query("id")
	idn, _ := strconv.Atoi(id)

	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err!= nil{
		return c.JSON(fiber.Map{
			"Status":  500,
			"message": "Tidak ditemukan product dengan Id" + id,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Detail Product Dengan Id " + id,
		"product": product,
	})
}

//GET products/deleteproduct/xx
func (controller *ProductController) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn,_ := strconv.Atoi(id)

	var product models.Product
	err:=models.DeleteProductById(controller.Db, &product, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Berhasil Menghapus Product dengan Id " + id,
	})
}

// func (controller *ProductController) EditProduct(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	idn, _ := strconv.Atoi(id)
// 	// fmt.Print("a")

// 	var product models.Product
// 	err := models.ReadProductById(controller.Db, &product, idn)
// 	// fmt.Print("b")

// 	if err != nil {
// 		// fmt.Print("c")
// 		return c.SendStatus(500) // http 500 internal server error
// 	}

// 	fmt.Print("d")
// 	return c.Render("editproduct", fiber.Map{
// 		"Title": "SiShop",
// 		"Content": "Edit Produk",
// 		"Product": product,
// 	})
// }

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
		return c.JSON(fiber.Map{
			"status":  400,
			"message": "Bad Request, Product Form is not complete",
		})
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
	er := models.UpdateProduct(controller.Db, &product)
	if er != nil {
		return c.SendStatus(500)
	}

	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Berhasil Mengubah Product dengan Id " + id,
	})
}