package conf

type YamlConf struct {
	Chain      string   `yaml:"chain" mapstructure:"chain"`
	DB         DB       `yaml:"db" mapstructure:"db"`
	ShareRPC   ShareRPC `yaml:"share_rpc" mapstructure:"share_rpc"`
	Log        Log      `yaml:"log" mapstructure:"log"`
	Index      Index    `yaml:"index" mapstructure:"index"`
	RPCService any      `yaml:"rpc_service" mapstructure:"rpc_service"`
}

type DB struct {
	Path string `yaml:"path" mapstructure:"path"`
}

type ShareRPC struct {
	Bitcoin Bitcoin `yaml:"bitcoin" mapstructure:"bitcoin"`
	Ord     Ord     `yaml:"ord" mapstructure:"ord"`
	Ordx    Ordx    `yaml:"ordx" mapstructure:"ordx"`
}

type Bitcoin struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
}

type Ord struct {
	URL string `yaml:"url" mapstructure:"url"`
}

type Ordx struct {
	URL string `yaml:"url" mapstructure:"url"`
}

type Log struct {
	Level string `yaml:"level" mapstructure:"level"`
	Path  string `yaml:"path" mapstructure:"path"`
}

type Index struct {
	ExitOnReachHeight bool `yaml:"exit_on_reach_height" mapstructure:"exit_on_reach_height"`
	EnableOrdxFix     bool `yaml:"enable_ordx_fix" mapstructure:"enable_ordx_fix"`
}
