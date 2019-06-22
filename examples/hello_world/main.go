package main

import (
	"fmt"
	"os"

	"github.com/draganm/decli"
)

type App struct {
	FirstName string `name:"first-name" usage:"your first name" aliases:"fn"`
	LastName  string `name:"last-name" usage:"your last name" aliases:"ln"`
	Age       int    `name:"age" usage:"your age"`
}

func (a App) Run([]string) error {
	fmt.Printf("Hello %s %s (%d)\n", a.FirstName, a.LastName, a.Age)
	return nil
}

func main() {
	decli.RunAndFinish(&App{FirstName: "John", LastName: "Doe", Age: -1}, os.Args)
}
