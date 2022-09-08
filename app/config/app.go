package config

import "time"

type App struct {
	Env            string        `default:"local"`
	ContextTimeout time.Duration `envconfig:"CONTEXT_TIMEOUT" default:"2s"`
}
