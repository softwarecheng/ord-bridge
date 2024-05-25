package g

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/softwarecheng/ord-bridge/common/log"
	"github.com/softwarecheng/ord-bridge/rpc/bitcoind"
	"github.com/softwarecheng/ord-bridge/rpc/ord"
	"github.com/softwarecheng/ord-bridge/rpc/ordx"
)

func InitRpc() error {
	logLvl, err := logrus.ParseLevel(Cfg.Log.Level)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %s", err)
	}

	log.Log.WithFields(logrus.Fields{
		"BitcoinChain":   Cfg.Chain,
		"BitcoinRPCHost": Cfg.ShareRPC.Bitcoin.Host,
		"BitcoinRPCPort": Cfg.ShareRPC.Bitcoin.Port,
		"BitcoinRPCUser": Cfg.ShareRPC.Bitcoin.User,
		"BitcoinRPCPass": Cfg.ShareRPC.Bitcoin.Password,
		"DbDataDir":      Cfg.DB.Path,
		"LogLevel":       logLvl,
		"LogPath":        Cfg.Log.Path,
	}).Info("bitcoin using configuration")
	err = bitcoind.InitBitconRpc(
		Cfg.ShareRPC.Bitcoin.Host,
		Cfg.ShareRPC.Bitcoin.Port,
		Cfg.ShareRPC.Bitcoin.User,
		Cfg.ShareRPC.Bitcoin.Password,
		false,
	)
	if err != nil {
		return err
	}
	ord.InitOrdRpc(Cfg.ShareRPC.Ord.URL, 3)
	ordx.InitOrdxRpc(Cfg.ShareRPC.Ordx.URL, 3)
	return nil
}

func InitSigInt() {
	count := 0
	SigInt = make(chan os.Signal, 100)
	signal.Notify(SigInt, os.Interrupt)
	go func() {
		for {
			<-SigInt
			count++
			log.Log.Infof("Received SIGINT (CTRL+C), count %d, 3 times will close db and force exit", count)
			if count >= 3 {
				ReleaseRes()
				os.Exit(1)
			} else if count == 1 {
				for index := range sigIntFuncList {
					go sigIntFuncList[index]()
				}
			}
		}
	}()
}

func registSigIntFunc(callback func()) {
	sigIntFuncList = append(sigIntFuncList, callback)
}

func ReleaseRes() {
	ReleaseAllDB()
}
