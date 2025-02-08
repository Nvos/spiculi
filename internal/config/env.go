package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadEnv() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Executable: %v", err)
	}

	envPath := filepath.Join(wd, ".env")
	_, err = os.Stat(envPath)
	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return err
	}

	open, err := os.Open(".env")
	if err != nil {
		return fmt.Errorf("open .env: %w", err)
	}
	defer open.Close()

	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid .env line: %s", line)
		}

		if err := os.Setenv(parts[0], parts[1]); err != nil {
			return fmt.Errorf("set env from .env line: %s", line)
		}
	}

	return nil
}
