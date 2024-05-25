package flag

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/softwarecheng/ord-bridge/common/util"
	"github.com/softwarecheng/ord-bridge/conf"
	serverCommon "github.com/softwarecheng/ord-bridge/service/common"
	"github.com/spf13/viper"
)

const CfgName = "config.yaml"

func loadConf(cfgPath string) (*conf.YamlConf, error) {
	cfgDir, cfgFileName, err := util.ParseFilePath(cfgPath)
	if err != nil {
		return nil, err
	}
	viper.AddConfigPath(cfgDir)
	viper.SetConfigType("yaml")
	viper.SetConfigName(cfgFileName)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	ret := &conf.YamlConf{}
	err = viper.Unmarshal(ret)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cfg: %s, error: %s", cfgPath, err)
	}

	return ret, nil
}

func defaultConf(chain string) (*conf.YamlConf, error) {
	bitcoinPort := 18332
	switch chain {
	case "mainnet":
		bitcoinPort = 8332
	case "testnet":
		bitcoinPort = 18332
	default:
		return nil, fmt.Errorf("unsupported chain: %s", chain)
	}
	ret := &conf.YamlConf{
		Chain: chain,
		DB: conf.DB{
			Path: "db" + chain,
		},
		ShareRPC: conf.ShareRPC{
			Bitcoin: conf.Bitcoin{
				Host:     "host",
				Port:     bitcoinPort,
				User:     "user",
				Password: "password",
			},
			Ord: conf.Ord{
				URL: "http://127.0.0.1:80",
			},
		},
		Log: conf.Log{
			Level: "error",
			Path:  "log",
		},
		Index: conf.Index{
			ExitOnReachHeight: false,
		},
		RPCService: serverCommon.RPCService{
			Addr:  "0.0.0.0:80",
			Proxy: chain,
			Swagger: serverCommon.Swagger{
				Host:    "127.0.0.0",
				Schemes: []string{"http"},
			},
			API: serverCommon.API{
				APIKeyList:     []serverCommon.APIKeyList{},
				NoLimitAPIList: []string{"/health"},
			},
		},
	}

	return ret, nil
}

func genDefaultCfg(chain string) error {
	conf, err := defaultConf(chain)
	if err != nil {
		return err
	}
	curDir, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath := curDir + string(filepath.Separator) + CfgName
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	// viper.Set("conf", conf)
	// util.ViperSetConfig(conf)
	// err = viper.WriteConfigAs(filePath)
	// if err != nil {
	// 	return err
	// }
	return nil
}
