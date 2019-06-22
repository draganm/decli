package decli

import (
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

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

func extractFlagsAndCommands(v reflect.Value, path []int) ([]func(c *cli.Context) error, []cli.Flag, []*cli.Command, error) {

	t := v.FieldByIndex(path).Type()

	if t.Kind() != reflect.Struct {
		return nil, nil, nil, errors.New("command must be a struct")
	}

	runBefore := [](func(c *cli.Context) error){}
	flags := []cli.Flag{}
	commands := []*cli.Command{}

	for i := 0; i < t.NumField(); i++ {

		newPath := append([]int(nil), path...)
		newPath = append(newPath, i)

		ft := t.Field(i)
		fv := v.FieldByIndex(newPath)

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
			runBefore = append(runBefore, func(c *cli.Context) error {
				fv.Set(reflect.ValueOf(c.String(name)))
				return nil
			})
			flags = append(flags, &cli.StringFlag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.String(),
				EnvVars:     envVars,
			})
		case reflect.Int:
			runBefore = append(runBefore, func(c *cli.Context) error {
				fv.Set(reflect.ValueOf(c.Int(name)))
				return nil
			})
			flags = append(flags, &cli.IntFlag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       int(fv.Int()),
				EnvVars:     envVars,
			})
		case reflect.Uint:
			runBefore = append(runBefore, func(c *cli.Context) error {
				fv.Set(reflect.ValueOf(c.Uint(name)))
				return nil
			})
			flags = append(flags, &cli.UintFlag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       uint(fv.Uint()),
				EnvVars:     envVars,
			})
		case reflect.Int64:
			if fv.Type().String() == "time.Duration" {
				runBefore = append(runBefore, func(c *cli.Context) error {
					fv.Set(reflect.ValueOf(c.Duration(name)))
					return nil
				})
				flags = append(flags, &cli.DurationFlag{
					Name:        name,
					Usage:       usage,
					Hidden:      hidden,
					Aliases:     aliases,
					DefaultText: defaultText,
					Value:       time.Duration(fv.Int()),
					EnvVars:     envVars,
				})
				continue
			}
			runBefore = append(runBefore, func(c *cli.Context) error {
				fv.Set(reflect.ValueOf(c.Int64(name)))
				return nil
			})
			flags = append(flags, &cli.Int64Flag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.Int(),
				EnvVars:     envVars,
			})
		case reflect.Uint64:
			runBefore = append(runBefore, func(c *cli.Context) error {
				fv.Set(reflect.ValueOf(c.Uint64(name)))
				return nil
			})
			flags = append(flags, &cli.Uint64Flag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.Uint(),
				EnvVars:     envVars,
			})
		case reflect.Float64:
			runBefore = append(runBefore, func(c *cli.Context) error {
				fv.Set(reflect.ValueOf(c.Float64(name)))
				return nil
			})
			flags = append(flags, &cli.Float64Flag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.Float(),
				EnvVars:     envVars,
			})
		case reflect.Bool:
			runBefore = append(runBefore, func(c *cli.Context) error {
				fv.Set(reflect.ValueOf(c.Bool(name)))
				return nil
			})
			flags = append(flags, &cli.BoolFlag{
				Name:        name,
				Usage:       usage,
				Hidden:      hidden,
				Aliases:     aliases,
				DefaultText: defaultText,
				Value:       fv.Bool(),
				EnvVars:     envVars,
			})
		case reflect.Struct:
			cmd, err := createCommand(v, newPath)
			if err != nil {
				return nil, nil, nil, err
			}
			commands = append(commands, cmd)
		default:
			return nil, nil, nil, errors.Errorf("not supported type %#v", fv.Type().String())
		}
	}

	return runBefore, flags, commands, nil

}

func Run(cmd interface{}, args []string) error {

	app := &cli.App{}

	v := reflect.ValueOf(cmd).Elem()

	runBefore, flags, commands, err := extractFlagsAndCommands(v, []int{})
	if err != nil {
		return errors.Wrap(err, "while configuring decli")
	}

	app.Commands = commands
	app.Flags = flags

	rc, isCommand := cmd.(Command)
	if isCommand {
		app.Action = func(c *cli.Context) error {
			for _, rbf := range runBefore {
				err = rbf(c)
				if err != nil {
					return err
				}
			}
			return rc.Run(c.Args().Slice())
		}
	}

	bf, isBefore := cmd.(Before)
	if isBefore {
		app.Before = func(c *cli.Context) error {
			for _, rbf := range runBefore {
				err = rbf(c)
				if err != nil {
					return err
				}
			}
			return bf.Before(c.Args().Slice())
		}
	}

	return app.Run(args)
}

func createCommand(v reflect.Value, path []int) (*cli.Command, error) {

	t := v.Type()
	sf := t.FieldByIndex(path)

	name := sf.Tag.Get("name")
	if name == "" {
		name = strcase.KebabCase(sf.Name)
	}

	runBefore, flags, commands, err := extractFlagsAndCommands(v, path)
	if err != nil {
		return nil, errors.Wrapf(err, "while configuring command %s", name)
	}

	cliCommand := &cli.Command{
		Name:        name,
		Flags:       flags,
		Subcommands: commands,
	}
	cm, isCommand := v.Interface().(Command)

	if isCommand {
		cliCommand.Action = func(c *cli.Context) error {
			for _, rbf := range runBefore {
				err = rbf(c)
				if err != nil {
					return err
				}
			}
			return cm.Run(c.Args().Slice())
		}
	}

	bf, isBefore := v.Interface().(Before)
	if isBefore {
		cliCommand.Before = func(c *cli.Context) error {
			for _, rbf := range runBefore {
				err = rbf(c)
				if err != nil {
					return err
				}
			}
			return bf.Before(c.Args().Slice())
		}
	}

	return cliCommand, nil
}
