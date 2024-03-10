package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type responseRam struct {
	Porcentaje int
	Mensaje    string
}

func PorcentajeActualRam(c *fiber.Ctx) error {
	cmd := exec.Command("sh", "-c", "cat /proc/ram_porcentaje_uso")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(responseRam{Mensaje: "Proceso no corriendo"})
	}
	porcentaje, err := strconv.Atoi(string(out[:]))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responseRam{Mensaje: "Porcentaje incorrecto"})
	}
	return c.Status(fiber.StatusOK).JSON(responseRam{Porcentaje: porcentaje})
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
	app.Get("/ram/actual", PorcentajeActualRam)
	app.Listen(":5000")
}
