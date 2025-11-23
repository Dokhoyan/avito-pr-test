package config

import (
	"os"

	"github.com/joho/godotenv"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

type HTTPConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}
