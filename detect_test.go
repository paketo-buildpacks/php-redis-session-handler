package phpredishandler_test

import (
	"errors"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/servicebindings"
	phpredishandler "github.com/paketo-buildpacks/php-redis-session-handler"
	"github.com/paketo-buildpacks/php-redis-session-handler/fakes"
	"github.com/sclevine/spec"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		detectBindingResolver *fakes.DetectBindingResolver
		detect                packit.DetectFunc
	)

	it.Before(func() {
		detectBindingResolver = &fakes.DetectBindingResolver{}
		detectBindingResolver.ResolveCall.Returns.BindingSlice = []servicebindings.Binding{
			{
				Type: "php-redis-session",
			},
		}
		detect = phpredishandler.Detect(detectBindingResolver)
	})

	it("requires php during launch and provides nothing", func() {
		result, err := detect(packit.DetectContext{
			Platform: packit.Platform{
				Path: "some-platform-path",
			},
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Plan).To(Equal(packit.BuildPlan{
			Requires: []packit.BuildPlanRequirement{
				{
					Name: "php",
					Metadata: phpredishandler.BuildPlanMetadata{
						Launch: true,
					},
				},
			},
			Provides: []packit.BuildPlanProvision{},
		}))

		Expect(detectBindingResolver.ResolveCall.Receives.Typ).To(Equal("php-redis-session"))
		Expect(detectBindingResolver.ResolveCall.Receives.Provider).To(Equal(""))
		Expect(detectBindingResolver.ResolveCall.Receives.PlatformDir).To(Equal("some-platform-path"))
	})

	context("there are no php-redis-session bindings provided", func() {
		it.Before(func() {
			detectBindingResolver.ResolveCall.Returns.BindingSlice = []servicebindings.Binding{}
		})

		it("detection fails", func() {
			_, err := detect(packit.DetectContext{
				Platform: packit.Platform{
					Path: "some-platform-path",
				},
			})
			Expect(err).To(MatchError(packit.Fail.WithMessage("no service bindings of type `php-redis-session` provided")))
		})
	})

	context("failure cases", func() {
		context("the binding resolver fails to resolve bindings", func() {
			it.Before(func() {
				detectBindingResolver.ResolveCall.Returns.Error = errors.New("failed to resolve bindings")
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{
					Platform: packit.Platform{
						Path: "some-platform-path",
					},
				})
				Expect(err).To(MatchError("failed to resolve bindings"))
			})
		})
	})
}
