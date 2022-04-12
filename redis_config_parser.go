package phpredishandler

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/paketo-buildpacks/packit/v2/fs"
)

type RedisConfig struct {
	Hostname string
	Port     int
	Password string
}

type RedisConfigParser struct {
}

func NewRedisConfigParser() RedisConfigParser {
	return RedisConfigParser{}
}

func (p RedisConfigParser) Parse(dir string) (RedisConfig, error) {
	hostFilepath := filepath.Join(dir, "host")
	hostnameFilepath := filepath.Join(dir, "hostname")

	hostFileExists, err := fs.Exists(hostFilepath)
	if err != nil {
		return RedisConfig{}, err
	}

	hostname := "127.0.0.1"
	if hostFileExists {
		hostnameBytes, err := os.ReadFile(hostFilepath)
		if err != nil {
			return RedisConfig{}, err
		}

		hostname = strings.TrimSpace(string(hostnameBytes))
	} else {
		hostnameFileExists, err := fs.Exists(hostnameFilepath)
		if err != nil {
			// untested
			return RedisConfig{}, err
		}

		if hostnameFileExists {
			hostnameBytes, err := os.ReadFile(hostnameFilepath)
			if err != nil {
				return RedisConfig{}, err
			}

			hostname = strings.TrimSpace(string(hostnameBytes))
		}
	}

	portFilepath := filepath.Join(dir, "port")

	portFileExists, err := fs.Exists(portFilepath)
	if err != nil {
		// untested
		return RedisConfig{}, err
	}

	port := 6379
	if portFileExists {
		portBytes, err := os.ReadFile(portFilepath)
		if err != nil {
			return RedisConfig{}, err
		}

		port, err = strconv.Atoi(strings.TrimSpace(string(portBytes)))
		if err != nil {
			return RedisConfig{}, err
		}
	}

	passwordFilepath := filepath.Join(dir, "password")

	passwordFileExists, err := fs.Exists(passwordFilepath)
	if err != nil {
		// untested
		return RedisConfig{}, err
	}

	password := ""
	if passwordFileExists {
		passwordBytes, err := os.ReadFile(passwordFilepath)
		if err != nil {
			return RedisConfig{}, err
		}

		password = strings.TrimSpace(string(passwordBytes))
	}

	return RedisConfig{
		Hostname: hostname,
		Port:     port,
		Password: password,
	}, nil
}
