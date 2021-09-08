package config

import (
	"github.com/juju/ratelimit"
	"github.com/patrickmn/go-cache"
	"time"
)

var RateLimitBucket *ratelimit.Bucket
var RateLimitCache = cache.New(24*time.Hour, 10*time.Second)
