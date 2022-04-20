package phpredishandler_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitPhpRedisHandler(t *testing.T) {
	suite := spec.New("php-redis-handler", spec.Report(report.Terminal{}), spec.Parallel())
	suite("Build", testBuild)
	suite("Detect", testDetect)
	suite("RedisConfigParser", testRedisConfigParser)
	suite("RedisConfigWriter", testRedisConfigWriter)
	suite.Run(t)
}
