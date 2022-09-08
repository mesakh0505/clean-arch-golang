package config

import "github.com/joho/godotenv"

type Server struct {
	Address string `default:":9090"`
}

// Constructs the server configuration by loading variables
// from the environment, plus the filenames provided.
func LoadForServer(filenames ...string) Server {
	// we do not care if there is no .env file.
	_ = godotenv.Overload(filenames...)

	r := Server{}

	mustLoad("SERVER", &r)

	return r
}
