package config

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func TestEnvExampleMustBeCompleteServer(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()

	_, b, _, _ := runtime.Caller(0)
	projectEnv := path.Join(filepath.Dir(b), "..", "..", ".env.example")
	_ = LoadForServer(projectEnv)
}
