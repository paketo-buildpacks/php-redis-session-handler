package phpredishandler_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/packit/v2/scribe"
	phpredishandler "github.com/paketo-buildpacks/php-redis-session-handler"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testRedisConfigWriter(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		layerDir          string
		cnbDir            string
		redisConfig       phpredishandler.RedisConfig
		redisConfigWriter phpredishandler.RedisConfigWriter
	)

	it.Before(func() {
		var err error
		layerDir, err = os.MkdirTemp("", "php-redis-layer")
		Expect(err).NotTo(HaveOccurred())
		Expect(os.Chmod(layerDir, os.ModePerm)).To(Succeed())

		cnbDir, err = os.MkdirTemp("", "cnb")
		Expect(err).NotTo(HaveOccurred())

		Expect(os.MkdirAll(filepath.Join(cnbDir, "config"), os.ModePerm)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(cnbDir, "config", "php-redis.ini"), []byte(`session.save_path = "{{.}}"`), os.ModePerm)).To(Succeed())

		redisConfig = phpredishandler.RedisConfig{
			Hostname: "some-hostname",
			Port:     1234,
			Password: "some-password",
		}
		logEmitter := scribe.NewEmitter(bytes.NewBuffer(nil))
		redisConfigWriter = phpredishandler.NewRedisConfigWriter(logEmitter)
	})

	it.After(func() {
		Expect(os.RemoveAll(layerDir)).To(Succeed())
		Expect(os.RemoveAll(cnbDir)).To(Succeed())
	})

	it("writes a redis config ini file into the redis config layer", func() {
		redisConfigFilePath, err := redisConfigWriter.Write(redisConfig, layerDir, cnbDir)
		Expect(err).NotTo(HaveOccurred())

		Expect(redisConfigFilePath).To(Equal(filepath.Join(layerDir, "php-redis.ini")))

		contents, err := os.ReadFile(redisConfigFilePath)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(contents)).To(ContainSubstring(`session.save_path = "tcp://some-hostname:1234?auth=some-password"`))
	})

	context("when there is no password", func() {
		it.Before(func() {
			redisConfig.Password = ""
		})

		it("writes a redis config ini file into the redis config layer without a password", func() {
			redisConfigFilePath, err := redisConfigWriter.Write(redisConfig, layerDir, cnbDir)
			Expect(err).NotTo(HaveOccurred())

			Expect(redisConfigFilePath).To(Equal(filepath.Join(layerDir, "php-redis.ini")))

			contents, err := os.ReadFile(redisConfigFilePath)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(contents)).To(ContainSubstring(`session.save_path = "tcp://some-hostname:1234"`))
		})
	})

	context("failure cases", func() {
		context("when template is not parseable", func() {
			it.Before(func() {
				Expect(os.WriteFile(filepath.Join(cnbDir, "config", "php-redis.ini"), []byte(`{{.`), os.ModePerm)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := redisConfigWriter.Write(redisConfig, layerDir, cnbDir)
				Expect(err).To(MatchError(ContainSubstring("failed to parse PHP redis config template")))
			})
		})

		context("when redis config file can't be opened for writing", func() {
			it.Before(func() {
				Expect(os.WriteFile(filepath.Join(layerDir, "php-redis.ini"), nil, 0400)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := redisConfigWriter.Write(redisConfig, layerDir, cnbDir)
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})
		})
	})
}
