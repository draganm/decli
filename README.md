# De(clarative) CLI

Purpose of DeCLI is to simplify parsing complex CLI arguments and defining rich CLI tools with minimal amount of required code.

DeCLI builds on top of [cli](https://gopkg.in/urfave/cli.v2). It relies heavily on defining CLI parameters as field tags in Golang structs instead of programmatically defining the arguments using Go code.

## Why should I use it?

DeCLI enables you to write more concise code and will add type safety to your arguments.

As an illustration, here's a hello world example of a DeCLI:

```go
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

func (a *App) Run([]string) error {
    fmt.Printf("Hello %s %s (%d)\n", a.FirstName, a.LastName, a.Age)
    return nil
}

func main() {
    decli.RunAndFinish(&App{FirstName: "John", LastName: "Doe", Age: -1}, os.Args)
}
```

as opposed to using raw [cli](https://gopkg.in/urfave/cli.v2):

```go
package main

import (
    "fmt"
    "log"
    "os"

    "gopkg.in/urfave/cli.v2"
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
```

Both sources do the same, DeCLI version has almost the half of lines of code.
