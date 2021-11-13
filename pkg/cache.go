package pkg

import (
	"github.com/patrickmn/go-cache"
)

var Cache = map[string]*cache.Cache{}

//cache.New(24*time.Hour, 10*time.Second)
