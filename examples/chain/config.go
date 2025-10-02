package main

import (
	"bufio"
	"os"
	"strings"
)

// loadEnvFile loads environment variables from a .env file
func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

// loadConfig loads configuration from environment or .env file
func loadConfig() {
	// Try to load from .env file first
	if err := loadEnvFile(".env"); err != nil {
		// If .env file doesn't exist, that's okay - use system environment variables
	}

	// Also try to load from chain directory
	if err := loadEnvFile("chain/.env"); err != nil {
		// If chain/.env file doesn't exist, that's okay - use system environment variables
	}
}
