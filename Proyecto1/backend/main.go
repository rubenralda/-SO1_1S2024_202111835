package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type responsePorcentaje struct {
	Porcentaje float64
	Mensaje    string
}

type mierror struct {
	mensaje string
}

func (m mierror) Error() string {
	return m.mensaje
}

func copiarPorcentajeRam() (float64, error) {
	cmd := exec.Command("sh", "-c", "cat /proc/ram_porcentaje_uso")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, mierror{mensaje: "Proceso no corriendo"}
	}
	porcentaje, err := strconv.ParseFloat(string(out[:]), 32)
	if err != nil {
		return 0, mierror{mensaje: "Porcentaje incorrecto"}
	}
	return porcentaje, nil
}

func copiarPorcentajeCPU() (float64, error) {
	cmd := exec.Command("sh", "-c", "cat /proc/cpu_uso")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, mierror{mensaje: "Proceso no corriendo"}
	}
	porcentaje, err := strconv.ParseFloat(string(out[:]), 32)
	if err != nil {
		return 0, mierror{mensaje: "Porcentaje incorrecto"}
	}
	return porcentaje, nil
}

func PorcentajeActualRam(c *fiber.Ctx) error {
	porcentaje, err := copiarPorcentajeRam()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responsePorcentaje{Mensaje: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(responsePorcentaje{Porcentaje: porcentaje})
}

func PorcentajeActualCPU(c *fiber.Ctx) error {
	porcentaje, err := copiarPorcentajeCPU()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responsePorcentaje{Mensaje: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(responsePorcentaje{Porcentaje: porcentaje})
}

func InsertarDatos() {
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/sopes_proyecto1")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	//insertar info ram
	porcentaje, err := copiarPorcentajeRam()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	currentTime := time.Now().Local().Format("2006-01-02 15:04:05")
	insert, err := db.Query("INSERT INTO ram(porcentaje, fecha_tiempo) values(?,?)", porcentaje, currentTime)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer insert.Close()

	//insertar inf de cpu
	porcentaje_cpu, err := copiarPorcentajeCPU()
	if err != nil {
		fmt.Println(err)
		return
	}
	insert2, err := db.Query("INSERT INTO cpu(porcentaje, fecha_tiempo) values(?,?)", porcentaje_cpu, currentTime)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer insert2.Close()
}

type responseRegistros struct {
	Mensaje   string
	Registros []interface{}
}

type ram struct {
	Id           int     `json:"id"`
	Porcentaje   float64 `json:"porcentaje"`
	Fecha_tiempo string  `json:"time"`
}

type cpu struct {
	Id           int     `json:"id"`
	Porcentaje   float64 `json:"porcentaje"`
	Fecha_tiempo string  `json:"time"`
}

func TodoRam() ([]interface{}, error) {
	var todoRam []interface{}
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/sopes_proyecto1")
	if err != nil {
		return todoRam, err
	}
	defer db.Close()

	registros, err := db.Query("SELECT * FROM ram")
	if err != nil {
		return todoRam, err
	}
	for registros.Next() {
		var ram ram
		err = registros.Scan(&ram.Id, &ram.Porcentaje, &ram.Fecha_tiempo)
		if err != nil {
			return todoRam, err
		}
		todoRam = append(todoRam, ram)
		//fmt.Println(ram.Fecha_tiempo)
	}
	return todoRam, nil
}

func TodoCPU() ([]interface{}, error) {
	var todoRam []interface{}
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/sopes_proyecto1")
	if err != nil {
		return todoRam, err
	}
	defer db.Close()

	registros, err := db.Query("SELECT * FROM cpu")
	if err != nil {
		return todoRam, err
	}
	for registros.Next() {
		var ram cpu
		err = registros.Scan(&ram.Id, &ram.Porcentaje, &ram.Fecha_tiempo)
		if err != nil {
			return todoRam, err
		}
		todoRam = append(todoRam, ram)
		//fmt.Println(ram.Fecha_tiempo)
	}
	return todoRam, nil
}

func GetTodoRam(c *fiber.Ctx) error {
	rams, err := TodoRam()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responseRegistros{Mensaje: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(responseRegistros{Registros: rams})
}

func GetTodoCPU(c *fiber.Ctx) error {
	rams, err := TodoCPU()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responseRegistros{Mensaje: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(responseRegistros{Registros: rams})
}

func main() {
	ticker := time.NewTicker(3 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				InsertarDatos()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
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
	app.Get("/cpu/actual", PorcentajeActualCPU)
	app.Get("/ram", GetTodoRam)
	app.Get("/cpu", GetTodoCPU)
	app.Listen(":5000")
}
