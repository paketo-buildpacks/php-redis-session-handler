package phpredishandler

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"text/template"

	"github.com/paketo-buildpacks/packit/v2/scribe"
)

type RedisConfigWriter struct {
	logger scribe.Emitter
}

func NewRedisConfigWriter(logger scribe.Emitter) RedisConfigWriter {
	return RedisConfigWriter{
		logger: logger,
	}
}

func (c RedisConfigWriter) Write(redisConfig RedisConfig, layerPath, cnbPath string) (string, error) {
	tmpl, err := template.New("php-redis.ini").ParseFiles(filepath.Join(cnbPath, "config", "php-redis.ini"))
	if err != nil {
		return "", fmt.Errorf("failed to parse PHP redis config template: %w", err)
	}

	sessionSavePath := fmt.Sprintf("tcp://%s:%d", redisConfig.Hostname, redisConfig.Port)
	c.logger.Debug.Subprocess("Including session save path: %s", sessionSavePath)

	if redisConfig.Password != "" {
		sessionSavePath = fmt.Sprintf("%s?auth=%s", sessionSavePath, url.QueryEscape(redisConfig.Password))
		c.logger.Debug.Subprocess("Including a password on the session save path")
	}

	// Configuration set by this buildpack
	var b bytes.Buffer
	err = tmpl.Execute(&b, sessionSavePath)
	if err != nil {
		// not tested
		return "", err
	}

	f, err := os.OpenFile(filepath.Join(layerPath, "php-redis.ini"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, &b)
	if err != nil {
		// not tested
		return "", err
	}

	return f.Name(), nil
}
