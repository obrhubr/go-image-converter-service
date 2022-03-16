package main

import (
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sunshineplan/imgconv"
)

type SuccessMessage struct {
	Message string `json:"message" xml:"message"`
}

type ErrorMessage struct {
	Message string `json:"message" xml:"message"`
}

func uploadImage(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		em := &ErrorMessage{
			Message: "Error while reading your request.",
		}
		return c.JSON(http.StatusInternalServerError, em)
	}

	src, err := file.Open()
	if err != nil {
		em := &ErrorMessage{
			Message: "Error while opening your image.",
		}
		return c.JSON(http.StatusInternalServerError, em)
	}
	defer src.Close()

	// Destination
	filename := uuid.New().String()
	dst, err := os.Create("./images/" + filename + "_orig.jpg")
	if err != nil {
		em := &ErrorMessage{
			Message: "Error while creating your image.",
		}
		return c.JSON(http.StatusInternalServerError, em)
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		em := &ErrorMessage{
			Message: "Error while copying your image.",
		}
		return c.JSON(http.StatusInternalServerError, em)
	}

	// Resize images
	imageToConvert, err := imgconv.Open("./images/" + filename + "_orig.jpg")
	if err != nil {
		em := &ErrorMessage{
			Message: "Error while creating your image.",
		}
		return c.JSON(http.StatusInternalServerError, em)
	}

	// Resize the image to the 3 different resolutions
	res1 := imgconv.Resize(imageToConvert, imgconv.ResizeOption{Width: 1920})
	res2 := imgconv.Resize(imageToConvert, imgconv.ResizeOption{Width: 1280})
	res3 := imgconv.Resize(imageToConvert, imgconv.ResizeOption{Width: 720})

	// Write the resulting image as jpg.
	res1Writer, err := os.Create("./images/" + filename + "_1.jpg")
	if err != nil {
		em := &ErrorMessage{
			Message: "Error while writing your image.",
		}
		return c.JSON(http.StatusInternalServerError, em)
	}
	err = imgconv.Write(res1Writer, res1, imgconv.FormatOption{Format: imgconv.JPEG})
	if err != nil {
		em := &ErrorMessage{
			Message: "Error while writing your image.",
		}
		return c.JSON(http.StatusInternalServerError, em)
	}
	res2Writer, err := os.Create("./images/" + filename + "_2.jpg")
	if err != nil {
		em := &ErrorMessage{
			Message: "Error while writing your image.",
		}
		return c.JSON(http.StatusInternalServerError, em)
	}
	err = imgconv.Write(res2Writer, res2, imgconv.FormatOption{Format: imgconv.JPEG})
	if err != nil {
		em := &ErrorMessage{
			Message: "Error while writing your image.",
		}
		return c.JSON(http.StatusInternalServerError, em)
	}
	res3Writer, err := os.Create("./images/" + filename + "_3.jpg")
	if err != nil {
		em := &ErrorMessage{
			Message: "Error while writing your image.",
		}
		return c.JSON(http.StatusInternalServerError, em)
	}
	err = imgconv.Write(res3Writer, res3, imgconv.FormatOption{Format: imgconv.JPEG})
	if err != nil {
		em := &ErrorMessage{
			Message: "Error while writing your image.",
		}
		return c.JSON(http.StatusInternalServerError, em)
	}

	sm := &SuccessMessage{
		Message: "Successfully uploaded and converted image.",
	}

	return c.JSON(http.StatusOK, sm)
}

func getImageWithRes(c echo.Context) error {
	image := c.Param("image")
	res := c.Param("res")

	switch res {
	case "1920px":
		return c.File("./images/" + image + "_1.jpg")
	case "1280px":
		return c.File("./images/" + image + "_2.jpg")
	case "720px":
		return c.File("./images/" + image + "_3.jpg")
	default:
		return c.File("./images/" + image + "_orig.jpg")
	}
}

func getImage(c echo.Context) error {
	image := c.Param("image")
	return c.File("./images/" + image + "_orig.jpg")
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/upload", uploadImage)
	e.GET("/:image", getImage)
	e.GET("/:image/:res", getImageWithRes)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
