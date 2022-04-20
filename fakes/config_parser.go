package fakes

import (
	"sync"

	phpredishandler "github.com/paketo-buildpacks/php-redis-session-handler"
)

type ConfigParser struct {
	ParseCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Dir string
		}
		Returns struct {
			RedisConfig phpredishandler.RedisConfig
			Error       error
		}
		Stub func(string) (phpredishandler.RedisConfig, error)
	}
}

func (f *ConfigParser) Parse(param1 string) (phpredishandler.RedisConfig, error) {
	f.ParseCall.mutex.Lock()
	defer f.ParseCall.mutex.Unlock()
	f.ParseCall.CallCount++
	f.ParseCall.Receives.Dir = param1
	if f.ParseCall.Stub != nil {
		return f.ParseCall.Stub(param1)
	}
	return f.ParseCall.Returns.RedisConfig, f.ParseCall.Returns.Error
}
