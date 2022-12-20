package args_test

import (
	"testing"

	"github.com/crdx/assert"
	"github.com/crdx/mission/args"
	"github.com/crdx/mission/util"
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
		if !util.Contains(env, test) {
			t.Errorf("env did not contain %s", test)
		}
	}

	assert.Equal(t, len(env), 4)
}
