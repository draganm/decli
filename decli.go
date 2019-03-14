package decli

import (
	"log"
	"reflect"
	"strings"

	"gopkg.in/urfave/cli.v2"
)

type Command interface {
	Run(args []string) error
}

func Run(cmd Command, args []string) error {
	t := reflect.TypeOf(cmd).Elem()

	fields := map[string]string{}

	app := &cli.App{
		Action: func(c *cli.Context) error {

			v := reflect.ValueOf(cmd).Elem()

			for fn, an := range fields {
				fv := v.FieldByName(fn)
				fv.SetString(c.String(an))
				log.Println("set", fn, "to", c.String(an))
			}

			return cmd.Run(nil)
		},
	}

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		t, found := ft.Tag.Lookup("decli")
		if found {
			parts := strings.Split(t, ",")
			name := parts[0]
			usage := ""
			for _, flag := range parts[1:] {
				strings.HasPrefix(flag, "usage:")
				usage = flag[len("usage:"):]
			}

			fields[ft.Name] = name
			app.Flags = append(app.Flags, &cli.StringFlag{
				Name:  name,
				Usage: usage,
			})
		}

	}

	return app.Run(args)
}
