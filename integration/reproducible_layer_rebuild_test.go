package integration_test

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testReproducibleLayerRebuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		pack   occam.Pack
		docker occam.Docker
		source string
		name   string
	)

	it.Before(func() {
		pack = occam.NewPack().WithVerbose()
		docker = occam.NewDocker()
	})

	context("when an app is rebuilt with pack build", func() {
		var (
			imageIDs  map[string]struct{}
			container occam.Container

			redisContainer occam.Container
			binding        string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())

			source, err = occam.Source(filepath.Join("testdata", "default_app"))
			Expect(err).NotTo(HaveOccurred())
			binding = filepath.Join(source, "binding")

			redisContainer, err = docker.Container.Run.
				WithPublish("6379").
				Execute("redis:latest")
			Expect(err).NotTo(HaveOccurred())

			ipAddress, err := redisContainer.IPAddressForNetwork("bridge")
			Expect(err).NotTo(HaveOccurred())

			Expect(os.WriteFile(filepath.Join(source, "binding", "host"), []byte(ipAddress), os.ModePerm)).To(Succeed())

			imageIDs = map[string]struct{}{}
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(redisContainer.ID)).To(Succeed())
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())

			for id := range imageIDs {
				Expect(docker.Image.Remove.Execute(id)).To(Succeed())
			}

			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("creates a layer with the same SHA as before", func() {
			var (
				err error

				firstImage  occam.Image
				secondImage occam.Image
			)

			build := pack.WithNoColor().Build.
				WithPullPolicy("never").
				WithBuildpacks(
					phpBuildpack,
					buildpack,
					phpBuiltinServerBuildpack,
				).
				WithEnv(map[string]string{
					"BP_PHP_WEB_DIR":       "htdocs",
					"BP_LOG_LEVEL":         "DEBUG",
					"SERVICE_BINDING_ROOT": "/bindings",
				}).
				WithVolumes(fmt.Sprintf("%s:/bindings/php-redis-session", binding))

			firstImage, _, err = build.Execute(name, source)
			Expect(err).NotTo(HaveOccurred())

			imageIDs[firstImage.ID] = struct{}{}
			Expect(firstImage.Buildpacks).To(HaveLen(3))

			Expect(firstImage.Buildpacks[1].Key).To(Equal(buildpackInfo.Buildpack.ID))
			Expect(firstImage.Buildpacks[1].Layers).To(HaveKey("php-redis-config"))

			secondImage, _, err = build.Execute(name, source)
			Expect(err).NotTo(HaveOccurred())

			imageIDs[secondImage.ID] = struct{}{}

			Expect(secondImage.Buildpacks).To(HaveLen(3))

			Expect(secondImage.Buildpacks[1].Key).To(Equal(buildpackInfo.Buildpack.ID))
			Expect(secondImage.Buildpacks[1].Layers).To(HaveKey("php-redis-config"))
			Expect(secondImage.Buildpacks[1].Layers["php-redis-config"].SHA).To(Equal(firstImage.Buildpacks[1].Layers["php-redis-config"].SHA))

			container, err = docker.Container.Run.
				WithEnv(map[string]string{"PORT": "8080"}).
				WithPublish("8080").
				Execute(secondImage.ID)
			Expect(err).NotTo(HaveOccurred())

			jar, err := cookiejar.New(nil)
			Expect(err).NotTo(HaveOccurred())

			client := &http.Client{
				Jar: jar,
			}

			Eventually(container).Should(Serve(ContainSubstring("1")).WithClient(client).OnPort(8080).WithEndpoint("/index.php"))
			Eventually(container).Should(Serve(ContainSubstring("2")).WithClient(client).OnPort(8080).WithEndpoint("/index.php"))
		})
	})
}
