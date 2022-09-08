package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Root struct {
	App      App
	Postgres Postgres
}

// Constructs the root configuration by loading variables
// from the environment, plus the filenames provided.
func Load(filenames ...string) Root {
	// we do not care if there is no .env file.
	_ = godotenv.Overload(filenames...)

	r := Root{
		App:      App{},
		Postgres: Postgres{},
	}

	mustLoad("APP", &r.App)
	mustLoad("POSTGRES", &r.Postgres)

	return r
}

func mustLoad(prefix string, spec interface{}) {
	err := envconfig.Process(prefix, spec)
	if err != nil {
		panic(err)
	}
}
