package g

import (
	"fmt"

	"github.com/softwarecheng/ord-bridge/common/log"
	"github.com/softwarecheng/ord-bridge/common/util"
	"github.com/softwarecheng/ord-bridge/service"
	serverCommon "github.com/softwarecheng/ord-bridge/service/common"
)

func InitRpcService() error {
	if OrdIndexerMgr.Indexer == nil {
		return fmt.Errorf("indexer not set")
	}
	Rpc = service.NewManager(OrdIndexerMgr.Indexer)
	addr := ""
	host := ""
	scheme := ""
	proxy := ""
	logPath := ""
	var apiCfgData any
	if Cfg != nil {
		rpcServiceCfg, err := util.ParseCfgField[serverCommon.RPCService](Cfg.RPCService)
		if err != nil {
			return err
		}
		addr = rpcServiceCfg.Addr
		host = rpcServiceCfg.Swagger.Host
		for _, v := range rpcServiceCfg.Swagger.Schemes {
			scheme += v + ","
		}
		proxy = rpcServiceCfg.Proxy
		if proxy[0] != '/' {
			proxy = "/" + proxy
		}
		logPath = rpcServiceCfg.LogPath
		apiCfgData = rpcServiceCfg.API

	} else {
		return fmt.Errorf("configuration not set")
	}

	err := Rpc.Start(addr, host, scheme,
		proxy, logPath, apiCfgData)
	if err != nil {
		return err
	}
	log.Log.Info("rpc started")

	return nil
}
