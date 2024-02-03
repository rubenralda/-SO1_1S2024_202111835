package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Resp struct {
	Carnet string
	Nombre string
	Err    bool
	//Reporte_error string
	//Message       string
}

func estudiantes(c *fiber.Ctx) error {
	response := Resp{
		Carnet: "202111835",
		Nombre: "Rubén Alejandro Ralda Mejia",
		Err:    false,
		//Message: "Ejecucion realizada",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	// Middleware para imprimir las peticiones
	app.Use(func(c *fiber.Ctx) error {
		currentTime := time.Now()
		fmt.Print(currentTime.Format("2006-01-02")+" ", c.OriginalURL(), " ", c.Method(), " ")
		return c.Next()
	})
	// Registro de la respuesta
	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		fmt.Println(c.Response().StatusCode())
		return err
	})
	app.Get("/data", estudiantes)
	app.Listen(":5000")
}
