package env

import (
	"os"
	"strings"
)

func GetVar(varName string) string {
	return os.Getenv(varName)
}

func GetVarOtherwiseDefault(varName string, defaultValue string) string {
	v, p := os.LookupEnv(varName)
	if p {
		return v
	}
	return defaultValue
}

func IsSet(varName string) bool {
	_, p := os.LookupEnv(varName)
	if p {
		return true
	}
	return false
}

func SetEnvVars(envVars map[string]string) {
	for key, val := range envVars {
		os.Setenv(strings.TrimSpace(key), strings.TrimSpace(val))
	}
}
