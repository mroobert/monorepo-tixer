// This package provides support to manage environment variables.
package env

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// LoadEnvFile loads environment variables from a file.
func LoadEnvFile() error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	file, err := os.Open(fmt.Sprintf("%s/.env", currDir))
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

// LoadEnv loads an environment variable or returns an error if it is not set.
func LoadEnv(env string) (string, error) {
	v := os.Getenv(env)
	if v == "" {
		return "", fmt.Errorf("env var %s not set", env)
	}
	return v, nil
}

// LoadEnvOrDefault loads an environment variable or returns a default value if it is not set.
func LoadEnvOrDefault(env string, defaultValue string) string {
	v := os.Getenv(env)
	if v == "" {
		return defaultValue
	}
	return v
}

// LoadIntEnvOrDefault loads an environment variable as an integer or returns a default value if it is not set.
func LoadIntEnvOrDefault(env string, defaultValue int) (int, error) {
	v := os.Getenv(env)
	if v == "" {
		return defaultValue, nil
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue, err
	}
	return i, nil
}

// LoadDurationEnvOrDefault loads an environment variable as a duration or returns a default value if it is not set.
func LoadDurationEnvOrDefault(env string, defaultValue time.Duration) (time.Duration, error) {
	v := os.Getenv(env)
	if v == "" {
		return defaultValue, nil
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return defaultValue, err
	}

	return d, nil
}
