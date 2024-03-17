package main

import (
	"database/sql"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/utils"
)

type responsePorcentaje struct {
	Porcentaje float64
	Mensaje    string
	Procesos   string
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

func GetProcesos(c *fiber.Ctx) error {
	cmd := exec.Command("sh", "-c", "cat /proc/procesos_cpu")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responsePorcentaje{Mensaje: "Proceso no corriendo"})
	}
	return c.Status(fiber.StatusOK).JSON(responsePorcentaje{Procesos: string(out[:])})
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
	Mensaje string
	Labels  []string
	Data    []float64
}

func TodoRam() ([]float64, []string, error) {
	var data []float64
	var labels []string
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/sopes_proyecto1")
	if err != nil {
		return data, labels, err
	}
	defer db.Close()

	registros, err := db.Query("SELECT porcentaje, fecha_tiempo FROM ram")
	if err != nil {
		return data, labels, err
	}
	for registros.Next() {
		var porcentaje float64
		var fecha_tiempo string
		err = registros.Scan(&porcentaje, &fecha_tiempo)
		if err != nil {
			return data, labels, err
		}
		data = append(data, porcentaje)
		labels = append(labels, fecha_tiempo)
	}
	return data, labels, nil
}

func TodoCPU() ([]float64, []string, error) {
	var data []float64
	var labels []string
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/sopes_proyecto1")
	if err != nil {
		return data, labels, err
	}
	defer db.Close()

	registros, err := db.Query("SELECT porcentaje, fecha_tiempo FROM cpu")
	if err != nil {
		return data, labels, err
	}
	for registros.Next() {
		var porcentaje float64
		var fecha_tiempo string
		err = registros.Scan(&porcentaje, &fecha_tiempo)
		if err != nil {
			return data, labels, err
		}
		data = append(data, porcentaje)
		labels = append(labels, fecha_tiempo)
	}
	return data, labels, nil
}

func GetTodoRam(c *fiber.Ctx) error {
	data, labels, err := TodoRam()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responseRegistros{Mensaje: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(responseRegistros{Labels: labels, Data: data})
}

func GetTodoCPU(c *fiber.Ctx) error {
	data, labels, err := TodoCPU()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responseRegistros{Mensaje: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(responseRegistros{Labels: labels, Data: data})
}

type responseSimulacion struct {
	Ok       bool
	Mensaje  string
	Estado   string
	Procesos []interface{}
}

func InsertarProceso(pid int) error {
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/sopes_proyecto1")
	if err != nil {
		return err
	}
	defer db.Close()

	currentTime := time.Now().Local().Format("2006-01-02 15:04:05")
	insert, err := db.Query("INSERT INTO procesos(pid, fecha_tiempo, estado) values(?,?,?)", pid, currentTime, "new")
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	defer insert.Close()
	return nil
}

func UpdateProceso(pid int, estado string) error {
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/sopes_proyecto1")
	if err != nil {
		return err
	}
	defer db.Close()

	insert, err := db.Query("UPDATE procesos set estado=? where pid=?", estado, pid)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	defer insert.Close()
	return nil
}

func SelectProceso(pid int) (string, error) {
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/sopes_proyecto1")
	if err != nil {
		return "", err
	}
	defer db.Close()

	registros, err := db.Query("SELECT estado FROM procesos where pid=?", pid)
	if err != nil {
		return "", err
	}
	var estado string
	for registros.Next() {

		err = registros.Scan(&estado)
		if err != nil {
			return "", err
		}
	}
	return estado, nil
}

func SelectProcesos() ([]interface{}, error) {
	var procesos []interface{}
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/sopes_proyecto1")
	if err != nil {
		return procesos, err
	}
	defer db.Close()

	registros, err := db.Query("SELECT id, estado, pid FROM procesos")
	if err != nil {
		return procesos, err
	}
	type estado struct {
		Id     int
		Pid    int
		Estado string
	}
	for registros.Next() {
		var estado estado
		err = registros.Scan(&estado.Id, &estado.Estado, &estado.Pid)
		if err != nil {
			return procesos, err
		}
		procesos = append(procesos, estado)
	}
	fmt.Println(procesos...)
	return procesos, nil
}

func GetProceso(c *fiber.Ctx) error {
	pidStr := utils.CopyString(c.Params("pid"))
	if pidStr == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: "Se requiere el parámetro 'pid'"})
	}
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: "El parámetro 'pid' debe ser un número entero"})
	}
	estado, err := SelectProceso(pid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: "Error al buscar en la BD"})
	}
	return c.Status(fiber.StatusOK).JSON(responseSimulacion{Estado: estado})
}

func GetProcesosSimulados(c *fiber.Ctx) error {
	procesos, err := SelectProcesos()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: "Error al buscar en la BD"})
	}
	return c.Status(fiber.StatusOK).JSON(responseSimulacion{Procesos: procesos})
}

func StartProcess(c *fiber.Ctx) error {
	cmd := exec.Command("bash", "-c", "./prueba1")
	err := cmd.Start()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: "Error al iniciar el proceso"})
	}

	// Obtener el PID del proceso hijo
	childPID := cmd.Process.Pid
	if InsertarProceso(childPID) != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: "Error al insertar a la BD"})
	}
	return c.Status(fiber.StatusOK).JSON(responseSimulacion{Ok: true, Mensaje: fmt.Sprintf("Proceso iniciado con PID del hijo: %d", childPID)})
}

func StopProcess(c *fiber.Ctx) error {
	pidStr := utils.CopyString(c.Params("pid"))
	if pidStr == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: "Se requiere el parámetro 'pid'"})
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: "El parámetro 'pid' debe ser un número entero"})
	}

	// Enviar señal SIGSTOP al proceso con el PID proporcionado stop kill -SIGSTOP PID
	cmd := exec.Command("kill", "-SIGSTOP", strconv.Itoa(pid))
	err = cmd.Run()
	if err != nil {
		mensaje := fmt.Sprintf("Error al detener el proceso con PID %d", pid)
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: mensaje})
	}
	if UpdateProceso(pid, "stop") != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: "Error al actualizar estado"})
	}
	mensaje := fmt.Sprintf("Proceso con PID %d detenido", pid)
	return c.Status(fiber.StatusOK).JSON(responseSimulacion{Ok: true, Mensaje: mensaje})
}

func ResumeProcess(c *fiber.Ctx) error {
	pidStr := utils.CopyString(c.Params("pid"))
	if pidStr == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: "Se requiere el parámetro 'pid'"})
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: "El parámetro 'pid' debe ser un número entero"})
	}
	// Enviar señal SIGCONT al proceso con el PID proporcionado kill -SIGCONT PID
	cmd := exec.Command("kill", "-SIGCONT", strconv.Itoa(pid))
	err = cmd.Run()
	if err != nil {
		mensaje := fmt.Sprintf("Error al reanudar el proceso con PID %d", pid)
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: mensaje})
	}
	if UpdateProceso(pid, "resume") != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: "Error al actualizar estado"})
	}
	mensaje := fmt.Sprintf("Proceso con PID %d reanudado", pid)
	return c.Status(fiber.StatusOK).JSON(responseSimulacion{Ok: true, Mensaje: mensaje})
}

func KillProcess(c *fiber.Ctx) error {
	pidStr := utils.CopyString(c.Params("pid"))
	if pidStr == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: "Se requiere el parámetro 'pid'"})
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: "El parámetro 'pid' debe ser un número entero"})
	}

	// Enviar señal SIGCONT al proceso con el PID proporcionado KILL -9 PID o KILL SIGTERM PID
	cmd := exec.Command("kill", "-9", strconv.Itoa(pid))
	err = cmd.Run()
	if err != nil {
		mensaje := fmt.Sprintf("Error al intentar terminar el proceso con PID  %d", pid)
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: mensaje})
	}
	if UpdateProceso(pid, "kill") != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responseSimulacion{Mensaje: "Error al actualizar estado"})
	}
	mensaje := fmt.Sprintf("Proceso con PID %d terminado", pid)
	return c.Status(fiber.StatusOK).JSON(responseSimulacion{Ok: true, Mensaje: mensaje})
}

func main() {
	/* ticker := time.NewTicker(3 * time.Second)
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
	}() */
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
	app.Get("/procesos", GetProcesos)
	app.Get("/procesos/simulados", GetProcesosSimulados)
	app.Get("/procesos/:pid", GetProceso)
	app.Get("/start", StartProcess)
	app.Get("/stop/:pid", StopProcess)
	app.Get("/resume/:pid", ResumeProcess)
	app.Get("/kill/:pid", KillProcess)
	app.Listen(":5000")
}
