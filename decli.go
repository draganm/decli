package decli

import (
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"gopkg.in/urfave/cli.v2"
)

type Command interface {
	Run(args []string) error
}

func Run(cmd Command, args []string) error {
	v := reflect.ValueOf(cmd).Elem()
	t := reflect.TypeOf(cmd).Elem()

	// fields := map[string]string{}

	app := &cli.App{
		Action: func(c *cli.Context) error {
			return cmd.Run(nil)
		},
	}

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)

		fv := v.Field(i)

		name := ft.Tag.Get("name")
		usage := ft.Tag.Get("usage")
		hidden, err := strconv.ParseBool(ft.Tag.Get("hidden"))
		if err != nil {
			hidden = false
		}

		defaultText := ft.Tag.Get("defaultText")

		var aliases []string

		if ft.Tag.Get("aliases") != "" {
			aliases = strings.Split(ft.Tag.Get("aliases"), " ")
		}

		switch fv.Type().Kind() {
		case reflect.String:
			app.Flags = append(app.Flags, &cli.StringFlag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Destination: (*string)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Int:
			app.Flags = append(app.Flags, &cli.IntFlag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Destination: (*int)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Uint:
			app.Flags = append(app.Flags, &cli.UintFlag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Destination: (*uint)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Int64:
			app.Flags = append(app.Flags, &cli.Int64Flag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Destination: (*int64)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Uint64:
			app.Flags = append(app.Flags, &cli.Uint64Flag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Destination: (*uint64)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Float64:
			app.Flags = append(app.Flags, &cli.Float64Flag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Destination: (*float64)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		}
	}

	return app.Run(args)
}
