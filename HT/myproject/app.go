package main

import (
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) DatosRam() string {
	fmt.Println("DATOS OBTENIDOS DESDE EL MODULO:")
	fmt.Println("")

	/* cmd := exec.Command("sh", "-c", "cat /proc/modulo_ram")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	output := string(out[:])
	myram := ramstruct{}
	errr := json.Unmarshal([]byte(output), &myram)
	if errr != nil {
		fmt.Println("Error:", err)
	}
	return myram */
	return "hola"
}

type ramstruct struct {
	Total_Ram         int `json:"Total_Ram"`
	Ram_en_Uso        int `json:"Ram_en_Uso"`
	Ram_libre         int `json:"Ram_libre"`
	Porcentaje_en_uso int `json:"Porcentaje_en_uso"`
}
