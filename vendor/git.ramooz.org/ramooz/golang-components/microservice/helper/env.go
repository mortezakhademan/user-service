package helper

import (
	"github.com/joho/godotenv"
	"os"
)

func SetEnvFile(filePath ...string) error {
	if err := godotenv.Load(filePath...); err != nil {
		return err
	}
	return nil
}

func SetEnvFromMap(env map[string]string) error {
	for key, value := range env {
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}
	return nil
}
