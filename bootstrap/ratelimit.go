package bootstrap

import (
	"github.com/aoaostar/v8cdn_panel/config"
	"github.com/juju/ratelimit"
	"time"
)

func InitRateLimit() {
	fillInterval := time.Second * time.Duration(config.Conf.RateLimit.FillInterval)
	config.RateLimitBucket = ratelimit.NewBucket(fillInterval, config.Conf.RateLimit.Capacity)
}
