# PHP Redis Session Handler Cloud Native Buildpack
A Cloud Native Buildpack for configuring a [Redis](https://redis.io/) [session
handler](https://www.php.net/manual/en/class.sessionhandler.php) in PHP apps.

The buildpack generates an `.ini` configuration snippet to allow for connecting
to an external Redis server as a session handler. The `host`, `port`, and `password` 
are configurable via service bindings.

## Integration

The PHP Redis Session Handler CNB provides nothing, and only requires
`php` at launch time. It detects on the presence of a service binding of
type `php-redis-session`.

## Service Binding Configuration

As mentioned above, the buildpack participates in the build if the user
provides a [service
binding](https://paketo.io/docs/howto/configuration/#bindings) of `type php-redis-session`.

The build command will look like:
```
pack build myapp --env SERVICE_BINDING_ROOT=/bindings --volume <absolute-path-to-binding>:/bindings/php-redis-session

```

Inside of the binding itself, the following configuration can be set:

- `host` or `hostname` (Default `127.0.0.1`): Redis instance IP address
- `port` (Default 6379): Redis instance port
- `password` (No default): Redis instance password, if there is one

The configurations from the service binding are parsed and used to create a
`php-redis.ini` file with session configurations. The `php-redis.ini` file is
available in the PHP Redis Session Handler buildpack layer on the image, and
its path is appended to the `PHP_INI_SCAN_DIR` for usage when the app starts up.

## Usage

To package this buildpack for consumption:

```
$ ./scripts/package.sh
```

This builds the buildpack's Go source using `GOOS=linux` by default. You can
supply another value as the first argument to `package.sh`.

## Run Tests

To run all unit tests, run:
```
./scripts/unit.sh
```

To run all integration tests, run:
```
./scripts/integration.sh
```

## Debug Logs
For extra debug logs from the image build process, set the `$BP_LOG_LEVEL`
environment variable to `DEBUG` at build-time (ex. `pack build my-app --env
BP_LOG_LEVEL=DEBUG` or through a  [`project.toml`
file](https://github.com/buildpacks/spec/blob/main/extensions/project-descriptor.md).
