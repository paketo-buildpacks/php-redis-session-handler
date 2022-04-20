package phpredishandler

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/paketo-buildpacks/packit/v2/servicebindings"
)

//go:generate faux --interface BuildBindingResolver --output fakes/build_binding_resolver.go
//go:generate faux --interface ConfigParser --output fakes/config_parser.go
//go:generate faux --interface ConfigWriter --output fakes/config_writer.go

type BuildBindingResolver interface {
	ResolveOne(typ, provider, platformDir string) (servicebindings.Binding, error)
}

type ConfigParser interface {
	Parse(dir string) (RedisConfig, error)
}

type ConfigWriter interface {
	Write(redisConfig RedisConfig, layerPath, cnbPath string) (string, error)
}

// Build will return a packit.BuildFunc that will be invoked during the build
// phase of the buildpack lifecycle.
//
func Build(redisBindingConfigParser ConfigParser, bindingResolver BuildBindingResolver, redisConfigWriter ConfigWriter, logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		logger.Debug.Process("Getting the layer associated with the redis configuration")
		phpRedisLayer, err := context.Layers.Get(PhpRedisLayer)
		if err != nil {
			return packit.BuildResult{}, err
		}
		logger.Debug.Subprocess(phpRedisLayer.Path)
		logger.Debug.Break()

		phpRedisLayer, err = phpRedisLayer.Reset()
		if err != nil {
			return packit.BuildResult{}, err
		}

		logger.Debug.Process("Resolving the %s service binding", RedisBindingType)
		binding, err := bindingResolver.ResolveOne(RedisBindingType, "", context.Platform.Path)
		if err != nil {
			return packit.BuildResult{}, err
		}
		logger.Debug.Break()

		logger.Debug.Process("Parsing the %s service binding", RedisBindingType)
		redisConfig, err := redisBindingConfigParser.Parse(binding.Path)
		if err != nil {
			return packit.BuildResult{}, err
		}
		logger.Debug.Break()

		// Use go templating to write the config file
		logger.Process("Writing the redis configuration")
		redisConfigPath, err := redisConfigWriter.Write(redisConfig, phpRedisLayer.Path, context.CNBPath)
		if err != nil {
			return packit.BuildResult{}, err
		}
		logger.Subprocess("Redis configuration written to: %s", redisConfigPath)
		logger.Break()

		phpRedisLayer.LaunchEnv.Append("PHP_INI_SCAN_DIR",
			phpRedisLayer.Path,
			string(os.PathListSeparator),
		)
		logger.EnvironmentVariables(phpRedisLayer)

		phpRedisLayer.Launch = true

		return packit.BuildResult{
			Layers: []packit.Layer{phpRedisLayer},
		}, nil
	}
}
