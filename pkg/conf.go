package pkg

type Config struct {
	AppName    string `mapstructure:"app_name"`
	Debug      bool   `mapstructure:"debug"`
	Listen     string `mapstructure:"listen"`
	Static     string `mapstructure:"static"`
	Cloudflare struct {
		Email   string `mapstructure:"email"`
		HostKey string `mapstructure:"host_key"`
		DefaultRecord string `mapstructure:"default_record"`
	}
	JwtSecret string `mapstructure:"jwt_secret"`
	RateLimit struct {
		Enabled      bool  `mapstructure:"enabled"`
		FillInterval int64 `mapstructure:"fill_interval"`
		Capacity     int64 `mapstructure:"capacity"`
	}
}

var Conf = new(Config)
