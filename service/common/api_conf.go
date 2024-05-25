package common

type RPCService struct {
	Addr    string  `yaml:"addr" mapstructure:"addr"`
	Proxy   string  `yaml:"proxy" mapstructure:"proxy"`
	LogPath string  `yaml:"log_path" mapstructure:"log_path"`
	Swagger Swagger `yaml:"swagger" mapstructure:"swagger"`
	API     API     `yaml:"api" mapstructure:"api"`
}

type Swagger struct {
	Host    string   `yaml:"host" mapstructure:"host"`
	Schemes []string `yaml:"schemes" mapstructure:"schemes"`
}

type API struct {
	APIKeyList     []APIKeyList `yaml:"apikey_list" mapstructure:"apikey_list"`
	NoLimitAPIList []string     `yaml:"nolimit_api_list" mapstructure:"nolimit_api_list"`
}

type APIKeyList struct {
	APIKey    string     `yaml:"api_key" mapstructure:"api_key"`
	UserName  string     `yaml:"user_name" mapstructure:"user_name"`
	RateLimit *RateLimit `yaml:"rate_limit" mapstructure:"rate_limit"`
}

type RateLimit struct {
	PerSecond int `yaml:"per_second" mapstructure:"per_second"`
	PerDay    int `yaml:"per_day" mapstructure:"per_day"`
}

type APIInfo struct {
	UserName  string     `yaml:"user_name" mapstructure:"user_name"`
	RateLimit *RateLimit `yaml:"rate_limit" mapstructure:"rate_limit"`
}
