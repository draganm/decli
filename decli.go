package decli

import (
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/pkg/errors"
	"github.com/stoewer/go-strcase"
	"gopkg.in/urfave/cli.v2"
)

type Command interface {
	Run(args []string) error
}

type Before interface {
	Before(args []string) error
}

func RunAndFinish(cmd interface{}, args []string) {
	err := Run(cmd, args)
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
}

func Run(cmd interface{}, args []string) error {
	v := reflect.ValueOf(cmd).Elem()
	t := reflect.TypeOf(cmd).Elem()

	app := &cli.App{}

	rc, isCommand := cmd.(Command)
	if isCommand {
		app.Action = func(c *cli.Context) error {
			return rc.Run(c.Args().Slice())
		}
	}

	bf, isBefore := cmd.(Before)
	if isBefore {
		app.Before = func(c *cli.Context) error {
			return bf.Before(c.Args().Slice())
		}
	}

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)

		name := ft.Tag.Get("name")

		if name == "" {
			name = strcase.KebabCase(ft.Name)
		}
		usage := ft.Tag.Get("usage")
		hidden, err := strconv.ParseBool(ft.Tag.Get("hidden"))
		if err != nil {
			hidden = false
		}

		var envVars []string
		envVarsString := ft.Tag.Get("envVars")
		if envVarsString == "" {
			envVarsString = strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
		}
		for _, p := range strings.Split(envVarsString, ",") {
			envVars = append(envVars, strings.TrimSpace(p))
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
				EnvVars:     envVars,
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
				EnvVars:     envVars,
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
				EnvVars:     envVars,
				Destination: (*uint)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Int64:
			if fv.Type().String() == "time.Duration" {
				app.Flags = append(app.Flags, &cli.DurationFlag{
					Name:        name,
					Usage:       usage,
					Hidden:      hidden,
					Aliases:     aliases,
					DefaultText: defaultText,
					Value:       time.Duration(fv.Int()),
					EnvVars:     envVars,
					Destination: (*time.Duration)(unsafe.Pointer(fv.Addr().Pointer())),
				})
				continue
			}
			app.Flags = append(app.Flags, &cli.Int64Flag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.Int(),
				EnvVars:     envVars,
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
				EnvVars:     envVars,
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
				EnvVars:     envVars,
				Destination: (*float64)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Bool:
			app.Flags = append(app.Flags, &cli.BoolFlag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.Bool(),
				EnvVars:     envVars,
				Destination: (*bool)(unsafe.Pointer(fv.Addr().Pointer())),
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
	if name == "" {
		name = strcase.KebabCase(sf.Name)
	}

	cmd := &cli.Command{
		Name: name,
	}
	cm, isCommand := v.Interface().(Command)

	if isCommand {
		cmd.Action = func(c *cli.Context) error {
			return cm.Run(c.Args().Slice())
		}
	}

	bf, isBefore := v.Interface().(Before)
	if isBefore {
		cmd.Before = func(c *cli.Context) error {
			return bf.Before(c.Args().Slice())
		}
	}

	v = v.Elem()

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)

		name := ft.Tag.Get("name")
		if name == "" {
			name = strcase.KebabCase(ft.Name)
		}

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

		var envVars []string
		envVarsString := ft.Tag.Get("envVars")
		if envVarsString == "" {
			envVarsString = strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
		}

		for _, p := range strings.Split(envVarsString, ",") {
			envVars = append(envVars, strings.TrimSpace(p))
		}

		switch fv.Type().Kind() {
		case reflect.String:
			cmd.Flags = append(cmd.Flags, &cli.StringFlag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.String(),
				EnvVars:     envVars,
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
				EnvVars:     envVars,
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
				EnvVars:     envVars,
				Destination: (*uint)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Int64:
			if fv.Type().String() == "time.Duration" {
				cmd.Flags = append(cmd.Flags, &cli.DurationFlag{
					Name:        name,
					Usage:       usage,
					Hidden:      hidden,
					Aliases:     aliases,
					DefaultText: defaultText,
					Value:       time.Duration(fv.Int()),
					EnvVars:     envVars,
					Destination: (*time.Duration)(unsafe.Pointer(fv.Addr().Pointer())),
				})
				continue
			}
			cmd.Flags = append(cmd.Flags, &cli.Int64Flag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.Int(),
				EnvVars:     envVars,
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
				EnvVars:     envVars,
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
				EnvVars:     envVars,
				Destination: (*float64)(unsafe.Pointer(fv.Addr().Pointer())),
			})
		case reflect.Bool:
			cmd.Flags = append(cmd.Flags, &cli.BoolFlag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.Bool(),
				EnvVars:     envVars,
				Destination: (*bool)(unsafe.Pointer(fv.Addr().Pointer())),
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
