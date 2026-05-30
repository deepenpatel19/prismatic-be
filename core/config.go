package core

import (
	"fmt"
	"os"
	"strings"
)

func Get(key string) string {
	if path := os.Getenv(key + "_FILE"); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			panic(fmt.Sprintf("failed to read %s from %s: %v", key, path, err))
		}
		return strings.TrimSpace(string(data))
	}
	return os.Getenv(key)
}

func MustGet(key string) string {
	v := Get(key)
	if v == "" {
		panic(fmt.Sprintf("required config %s is not set", key))
	}
	return v
}
