package main

import (
	"github.com/softwarecheng/ord-bridge/common/log"
	"github.com/softwarecheng/ord-bridge/main/flag"
	"github.com/softwarecheng/ord-bridge/main/g"
)

func init() {
	err := flag.ParseCmdParams()
	if err != nil {
		log.Log.Fatal(err)
	}
	err = g.InitRpc()
	if err != nil {
		log.Log.Fatal(err)
	}
	err = g.InitDB()
	if err != nil {
		log.Log.Fatal(err)
	}
	g.InitSigInt()
}

func main() {
	log.Log.Info("Starting...")
	defer func() {
		g.ReleaseRes()
		log.Log.Info("shut down")
	}()
	var err error
	exitProgramChan := make(chan bool)
	err = g.InitOrdIndexer(exitProgramChan)
	if err != nil {
		log.Log.Error(err)
		return
	}

	err = g.InitRpcService()
	if err != nil {
		log.Log.Error(err)
		return
	}

	for range exitProgramChan {
		if g.OrdIndexerMgr.Indexer == nil {
			break
		}
	}
	log.Log.Info("prepare to release resource...")
}
