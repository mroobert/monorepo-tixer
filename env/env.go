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

// LoadEnvFile loads environment variables from .env file.
func LoadEnvFile() error {
	currDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	file, err := os.Open(fmt.Sprintf("%s/.env", currDir))
	if err != nil {
		return fmt.Errorf("failed to open .env file: %w", err)
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

func LoadEnv(env string) (string, error) {
	v := os.Getenv(env)
	if v == "" {
		return "", fmt.Errorf("env var %s not set", env)
	}
	return v, nil
}

func LoadEnvOrDefault(env string, defaultValue string) string {
	v := os.Getenv(env)
	if v == "" {
		return defaultValue
	}
	return v
}

func LoadInt32EnvOrDefault(env string, defaultValue int32) (int32, error) {
	v := os.Getenv(env)
	if v == "" {
		return defaultValue, nil
	}

	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return defaultValue, err
	}
	return int32(i), nil
}

func LoadInt64EnvOrDefault(env string, defaultValue int64) (int64, error) {
	v := os.Getenv(env)
	if v == "" {
		return defaultValue, nil
	}

	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return defaultValue, err
	}
	return i, nil
}

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
