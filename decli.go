package decli

import (
	"log"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

type Command interface {
	Run(args []string) error
}

func RunAndFinish(cmd Command, args []string) {
	err := Run(cmd, args)
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
}

func Run(cmd Command, args []string) error {
	v := reflect.ValueOf(cmd).Elem()
	t := reflect.TypeOf(cmd).Elem()

	app := &cli.App{
		Action: func(c *cli.Context) error {

			return cmd.Run(c.Args().Slice())
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
				Value:       fv.String(),
				Destination: (*string)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Int:
			app.Flags = append(app.Flags, &cli.IntFlag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       int(fv.Int()),
				Destination: (*int)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Uint:
			app.Flags = append(app.Flags, &cli.UintFlag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       uint(fv.Uint()),
				Destination: (*uint)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Int64:
			app.Flags = append(app.Flags, &cli.Int64Flag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.Int(),
				Destination: (*int64)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Uint64:
			app.Flags = append(app.Flags, &cli.Uint64Flag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.Uint(),
				Destination: (*uint64)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Float64:
			app.Flags = append(app.Flags, &cli.Float64Flag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.Float(),
				Destination: (*float64)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		// case reflect.Struct:
		case reflect.Ptr:
			cmd, err := createCommand(fv, ft)
			if err != nil {
				return err
			}
			app.Commands = append(app.Commands, cmd)
		default:
			return errors.Errorf("not supported type %#v", fv.Type().String())
		}
	}

	return app.Run(args)
}

func createCommand(v reflect.Value, sf reflect.StructField) (*cli.Command, error) {
	name := sf.Tag.Get("name")
	cm, isCommand := v.Interface().(Command)
	if !isCommand {
		return nil, errors.Errorf("%#v is not a Command", v.Kind())
	}

	v = v.Elem()

	cmd := &cli.Command{
		Name: name,
		Action: func(c *cli.Context) error {
			return cm.Run(c.Args().Slice())
		},
	}

	t := v.Type()

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
			log.Println("name", name)
			cmd.Flags = append(cmd.Flags, &cli.StringFlag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.String(),
				Destination: (*string)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Int:
			cmd.Flags = append(cmd.Flags, &cli.IntFlag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       int(fv.Int()),
				Destination: (*int)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Uint:
			cmd.Flags = append(cmd.Flags, &cli.UintFlag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       uint(fv.Uint()),
				Destination: (*uint)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Int64:
			cmd.Flags = append(cmd.Flags, &cli.Int64Flag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.Int(),
				Destination: (*int64)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Uint64:
			cmd.Flags = append(cmd.Flags, &cli.Uint64Flag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.Uint(),
				Destination: (*uint64)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Float64:
			cmd.Flags = append(cmd.Flags, &cli.Float64Flag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.Float(),
				Destination: (*float64)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		// case reflect.Struct:
		case reflect.Ptr:
			c, err := createCommand(fv, ft)
			if err != nil {
				return nil, err
			}
			cmd.Subcommands = append(cmd.Subcommands, c)
		default:
			return nil, errors.Errorf("not supported type %#v", fv.Type().String())
		}
	}

	return cmd, nil
}
