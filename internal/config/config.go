package config

import (
	"os"

	"github.com/joho/godotenv"
)

func Load(path string) error {
	// Пытаемся загрузить .env файл, но не падаем, если его нет
	// Переменные окружения могут быть заданы через систему
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
