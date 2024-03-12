package args_test

import (
	"testing"

	"crdx.org/mission/args"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestToEnvironmentVariables(t *testing.T) {
	args := args.New(map[string]string{"foo": "bar", "baz": "foo"}, "anon", "pass")
	env := args.ToEnvironmentVariables()

	tests := []string{
		"FOO_DIR=bar",
		"BAZ_DIR=foo",
		"TARGET_USER=anon",
		"PASS_BIN=pass",
	}

	for _, test := range tests {
		if !slices.Contains(env, test) {
			t.Errorf("env did not contain %s", test)
		}
	}

	assert.Equal(t, 4, len(env))
}
