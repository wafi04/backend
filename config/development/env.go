package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadConfig(env string) error {
	rootDir, err := getRootDir()
	if err != nil {
		return fmt.Errorf("failed to get root directory: %v", err)
	}

	envFile := fmt.Sprintf(".env.%s", env)
	err = godotenv.Load(filepath.Join(rootDir, envFile))
	if err != nil {
		return fmt.Errorf("error loading %s file: %v", envFile, err)
	}

	log.Printf("Loaded environment configuration from %s", envFile)
	return nil
}

func LoadEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Warning: Environment variable %s not found", key)
	}
	return value
}

func getRootDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find project root (no go.mod file found)")
		}
		dir = parent
	}
}
