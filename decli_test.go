package decli_test

import (
	"testing"

	"github.com/draganm/decli"
	"github.com/stretchr/testify/require"
)

type MyApp struct {
	SomeString string `decli:"some-string,required,usage:a string or something"`
}

func (m *MyApp) Run(args []string) error {
	return nil
}

func TestDecli(t *testing.T) {

	require := require.New(t)

	x := &MyApp{}
	// err := decli.Run(x, []string{"whatevs", "help"})
	err := decli.Run(x, []string{"whatevs", "--some-string", "abc"})
	require.Nil(err)
	require.Equal("abc", x.SomeString)

}
