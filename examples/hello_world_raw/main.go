package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "first-name",
				Usage:   "your first name",
				Aliases: []string{"fn"},
				Value:   "John",
			},
			&cli.StringFlag{
				Name:    "last-name",
				Usage:   "your last name",
				Aliases: []string{"ln"},
				Value:   "Doe",
			},
			&cli.IntFlag{
				Name:  "age",
				Usage: "your age",
				Value: -1,
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Hello %s %s (%d)\n", c.String("first-name"), c.String("last-name"), c.Int("age"))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
}
