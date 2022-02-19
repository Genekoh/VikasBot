package env

import (
	"os"
	"strings"
)

const defaultFile = ".env"

func readFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func GetEnviron(envPath string) (map[string]string, error) {
	if envPath == "" {
		envPath = defaultFile
	}

	data, err := readFile(envPath)
	if err != nil {
		return nil, err
	}

	environ := make(map[string]string)

	lines := strings.Split(data, "\n")
	for _, l := range lines {
		kv := strings.SplitN(l, "=", 2)

		if len(kv) == 0 {
			continue
		}

		key := kv[0]
		value := removeSuffix(kv[1])
		environ[key] = value
	}
	return environ, nil
}

func removeSuffix(s string) string {
	x := strings.TrimSuffix(s, "\n")
	return strings.TrimSuffix(x, "\r")
}
