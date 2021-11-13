package bootstrap

import (
	"github.com/aoaostar/v8cdn_panel/pkg"
	"github.com/juju/ratelimit"
	"time"
)

func InitRateLimit() {
	fillInterval := time.Second * time.Duration(pkg.Conf.RateLimit.FillInterval)
	pkg.RateLimitBucket = ratelimit.NewBucket(fillInterval, pkg.Conf.RateLimit.Capacity)
}
