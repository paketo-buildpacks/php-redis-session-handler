package phpredishandler_test

import (
	"os"
	"path/filepath"
	"testing"

	phpredishandler "github.com/paketo-buildpacks/php-redis-session-handler"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testRedisConfigParser(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir string

		parser phpredishandler.RedisConfigParser
	)

	it.Before(func() {
		var err error

		workingDir, err = os.MkdirTemp("", "workingDir")
		Expect(err).NotTo(HaveOccurred())

		parser = phpredishandler.NewRedisConfigParser()
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	it("parses with default values", func() {
		config, err := parser.Parse(workingDir)
		Expect(err).NotTo(HaveOccurred())

		Expect(config).To(Equal(phpredishandler.RedisConfig{
			Hostname: "127.0.0.1",
			Port:     6379,
			Password: "",
		}))
	})

	context("when the host file exists", func() {
		it.Before(func() {
			Expect(os.WriteFile(filepath.Join(workingDir, "host"), []byte("some-host"), os.ModePerm)).To(Succeed())
		})

		it("uses the value from the host file", func() {
			config, err := parser.Parse(workingDir)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Hostname).To(Equal("some-host"))
		})

		context("when there is whitespace", func() {
			it.Before(func() {
				Expect(os.WriteFile(filepath.Join(workingDir, "host"), []byte("  \tsome-host\n\n"), os.ModePerm)).To(Succeed())
			})

			it("strips whitespace", func() {
				config, err := parser.Parse(workingDir)
				Expect(err).NotTo(HaveOccurred())

				Expect(config.Hostname).To(Equal("some-host"))
			})
		})
	})

	context("when the hostname file exists", func() {
		it.Before(func() {
			Expect(os.WriteFile(filepath.Join(workingDir, "hostname"), []byte("some-other-host"), os.ModePerm)).To(Succeed())
		})

		it("uses the value from the hostname file", func() {
			config, err := parser.Parse(workingDir)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Hostname).To(Equal("some-other-host"))
		})

		context("when there is whitespace", func() {
			it.Before(func() {
				Expect(os.WriteFile(filepath.Join(workingDir, "hostname"), []byte("  \tsome-other-host\n\n"), os.ModePerm)).To(Succeed())
			})

			it("strips whitespace", func() {
				config, err := parser.Parse(workingDir)
				Expect(err).NotTo(HaveOccurred())

				Expect(config.Hostname).To(Equal("some-other-host"))
			})
		})
	})

	context("when the port file exists", func() {
		it.Before(func() {
			Expect(os.WriteFile(filepath.Join(workingDir, "port"), []byte("1234"), os.ModePerm)).To(Succeed())
		})

		it("uses the value from the port file", func() {
			config, err := parser.Parse(workingDir)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Port).To(Equal(1234))
		})

		context("when there is whitespace", func() {
			it.Before(func() {
				Expect(os.WriteFile(filepath.Join(workingDir, "port"), []byte("  \t1234\n\n"), os.ModePerm)).To(Succeed())
			})

			it("strips whitespace", func() {
				config, err := parser.Parse(workingDir)
				Expect(err).NotTo(HaveOccurred())

				Expect(config.Port).To(Equal(1234))
			})
		})
	})

	context("when the password file exists", func() {
		it.Before(func() {
			Expect(os.WriteFile(filepath.Join(workingDir, "password"), []byte("some-password"), os.ModePerm)).To(Succeed())
		})

		it("uses the value from the password file", func() {
			config, err := parser.Parse(workingDir)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Password).To(Equal("some-password"))
		})

		context("when there is whitespace", func() {
			it.Before(func() {
				Expect(os.WriteFile(filepath.Join(workingDir, "password"), []byte("  \tsome-password\n\n"), os.ModePerm)).To(Succeed())
			})

			it("strips whitespace", func() {
				config, err := parser.Parse(workingDir)
				Expect(err).NotTo(HaveOccurred())

				Expect(config.Password).To(Equal("some-password"))
			})
		})
	})

	context("failure cases", func() {
		context("when there is an error determining if the files exist", func() {
			it.Before(func() {
				Expect(os.Chmod(workingDir, 0000)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := parser.Parse(workingDir)
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})
		})

		context("when there is an error reading the host file", func() {
			it.Before(func() {
				Expect(os.WriteFile(filepath.Join(workingDir, "host"), []byte("some-host"), 0000)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := parser.Parse(workingDir)
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})
		})

		context("when there is an error reading the hostname file", func() {
			it.Before(func() {
				Expect(os.WriteFile(filepath.Join(workingDir, "hostname"), []byte("some-host"), 0000)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := parser.Parse(workingDir)
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})
		})

		context("when there is an error reading the port file", func() {
			it.Before(func() {
				Expect(os.WriteFile(filepath.Join(workingDir, "port"), []byte("some-port"), 0000)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := parser.Parse(workingDir)
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})
		})

		context("when port file contents cannot be parsed as an int", func() {
			it.Before(func() {
				Expect(os.WriteFile(filepath.Join(workingDir, "port"), []byte("not-an-int"), os.ModePerm)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := parser.Parse(workingDir)
				Expect(err).To(MatchError(ContainSubstring("invalid syntax")))
			})
		})

		context("when there is an error reading the password file", func() {
			it.Before(func() {
				Expect(os.WriteFile(filepath.Join(workingDir, "password"), []byte("some-password"), 0000)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := parser.Parse(workingDir)
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})
		})
	})
}
