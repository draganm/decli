package decli_test

import (
	"os"
	"testing"
	"time"

	"github.com/draganm/decli"
	"github.com/stretchr/testify/require"
)

type MyApp struct {
	SomeString   string
	SomeInt      int     `name:"some-int" usage:"an int"`
	SomeUint     uint    `name:"some-uint" usage:"a uint"`
	SomeInt64    int64   `name:"some-int64" usage:"an int64"`
	SomeUInt64   uint64  `name:"some-uint64" usage:"a uint64"`
	SomeFloat64  float64 `name:"some-float64" usage:"a float64"`
	SomeDuration time.Duration
	SomeBool     bool
	Sub          SubCommand
}

type SubCommand struct {
	Foo string `name:"foo"`
}

func (mc SubCommand) Run(args []string) error {
	return nil
}

func (m MyApp) Run(args []string) error {
	return nil
}

func TestDecli(t *testing.T) {

	require := require.New(t)

	x := &MyApp{
		Sub: SubCommand{},
	}

	os.Setenv("SOME_FLOAT64", "12.3")
	err := decli.Run(x, []string{
		"whatevs",
		"--some-string", "abc",
		"--some-int", "123",
		"--some-uint", "456",
		"--some-int64", "234",
		"--some-uint64", "789",
		"--some-duration", "5ms",
		"--some-bool", "true",
	})

	require.Nil(err)
	require.Equal("abc", x.SomeString)
	require.Equal(123, x.SomeInt)
	require.Equal(uint(456), x.SomeUint)
	require.Equal(int64(234), x.SomeInt64)
	require.Equal(uint64(789), x.SomeUInt64)
	require.Equal(12.3, x.SomeFloat64)
	require.Equal(5*time.Millisecond, x.SomeDuration)
	require.Equal(true, x.SomeBool)

}

func TestSubCommand(t *testing.T) {

	require := require.New(t)

	x := &MyApp{
		Sub: SubCommand{
			// Foo: "abc",
		},
	}
	err := decli.Run(x, []string{
		"whatevs",
		"sub",
		"--foo", "abc",
	})

	require.Nil(err)
	require.Equal("abc", x.Sub.Foo)

}
