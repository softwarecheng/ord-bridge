package flag

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/softwarecheng/ord-bridge/common/log"
	"github.com/softwarecheng/ord-bridge/common/util"
	"github.com/softwarecheng/ord-bridge/conf"
	"github.com/softwarecheng/ord-bridge/main/g"
	"github.com/spf13/viper"
)

var defaultCfgMap = map[string]any{
	"log.level":                   "info",
	"log.path":                    "log",
	"db.path":                     "db",
	"rpc_service.addr":            "log",
	"rpc_service.proxy":           "/",
	"rpc_service.log_path":        "log",
	"rpc_service.swagger.host":    "127.0.0.1",
	"rpc_service.swagger.schemes": []string{"http"},
}

func initConf(cfgPath string) (err error) {
	for k, v := range defaultCfgMap {
		viper.SetDefault(k, v)
	}
	g.Cfg, err = loadConf(cfgPath)
	if err != nil {
		return err
	}
	g.Cfg.Log.Path, err = util.ParseDirPath(g.Cfg.Log.Path)
	if err != nil {
		return err
	}
	g.Cfg.DB.Path, err = util.ParseDirPath(g.Cfg.DB.Path)
	if err != nil {
		return err
	}
	return nil
}

func initLog(cfg *conf.Log) error {
	var writers []io.Writer

	logPath := cfg.Path
	level := cfg.Level
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %s", err)
	}

	if logPath != "" {
		exePath, _ := os.Executable()
		executableName := filepath.Base(exePath)
		if strings.Contains(executableName, "debug") {
			executableName = "debug"
		}
		fileHook, err := rotatelogs.New(
			logPath+"/"+executableName+".%Y%m%d.log",
			rotatelogs.WithLinkName(logPath+"/"+executableName+".log"),
			rotatelogs.WithMaxAge(24*time.Hour),
			rotatelogs.WithRotationTime(1*time.Hour),
		)
		if err != nil {
			return fmt.Errorf("failed to create RotateFile hook, error: %s", err)
		}
		writers = append(writers, fileHook)
	}
	writers = append(writers, os.Stdout)
	log.Log.SetOutput(io.MultiWriter(writers...))
	log.Log.SetLevel(lvl)
	return nil
}
