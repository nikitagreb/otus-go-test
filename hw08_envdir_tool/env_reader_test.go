package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	testData := []struct {
		name     string
		dir      string
		expected Environment
	}{
		{
			"basic test",
			"./testdata/env",
			Environment{
				"BAR": EnvValue{
					Value:      "bar",
					NeedRemove: false,
				},
				"EMPTY": EnvValue{
					Value:      "",
					NeedRemove: true,
				},
				"FOO": EnvValue{
					Value:      "   foo\nwith new line",
					NeedRemove: false,
				},
				"HELLO": EnvValue{
					Value:      "\"hello\"",
					NeedRemove: false,
				},
				"UNSET": EnvValue{
					Value:      "",
					NeedRemove: true,
				},
			},
		},
	}

	for _, tData := range testData {
		t.Run(tData.name, func(t *testing.T) {
			envs, err := ReadDir(tData.dir)

			require.NoError(t, err)
			require.Equal(t, envs, tData.expected)
		})
	}
}

func TestReadDirNotFound(t *testing.T) {
	environment, err := ReadDir("./testdata/foobar")
	require.Error(t, err)
	require.Nil(t, environment)
}
