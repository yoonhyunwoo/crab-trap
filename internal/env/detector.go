package env

import (
	"os"
	"runtime"
)

type OSType string

const (
	OSUnix    OSType = "unix"
	OSWindows OSType = "windows"
)

type Environment struct {
	OS      OSType   `json:"os"`
	Vars    []string `json:"vars"`
	Pattern string   `json:"pattern"`
}

func Detect() Environment {
	osType := OSUnix
	if runtime.GOOS == "windows" {
		osType = OSWindows
	}

	return Environment{
		OS:   osType,
		Vars: GetEnvVars(osType),
		Pattern: GetPattern(osType),
	}
}

func GetEnvVars(osType OSType) []string {
	vars := []string{}

	switch osType {
	case OSUnix:
		vars = []string{"HOSTNAME", "USER", "HOME", "PWD", "PATH", "SHELL"}
	case OSWindows:
		vars = []string{"COMPUTERNAME", "USERNAME", "USERPROFILE", "CD", "PATH", "COMSPEC"}
	}

	return vars
}

func GetPattern(osType OSType) string {
	switch osType {
	case OSUnix:
		return "$"
	case OSWindows:
		return "%"
	}
	return "$"
}

func GetValue(varName string) string {
	return os.Getenv(varName)
}

func GetAllValues(osType OSType) map[string]string {
	result := make(map[string]string)

	for _, varName := range GetEnvVars(osType) {
		result[varName] = GetValue(varName)
	}

	return result
}
