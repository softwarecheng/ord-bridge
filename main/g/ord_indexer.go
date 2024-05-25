package g

import (
	"fmt"

	"github.com/softwarecheng/ord-bridge/common/log"
	"github.com/softwarecheng/ord-bridge/indexer"
)

func InitOrdIndexer(receiveChan chan bool) error {
	var err error
	var chain string
	if Cfg != nil {
		chain = Cfg.Chain
	} else {
		return fmt.Errorf("cfg is not set")
	}
	OrdIndexerMgr.Indexer, err = indexer.NewIndexer(DB.Ord, chain)
	if err != nil {
		return err
	}

	OrdIndexerMgr.Indexer.EnableOrdxFix = Cfg.Index.EnableOrdxFix
	OrdIndexerMgr.Indexer.ExitOnReachHeight = Cfg.Index.ExitOnReachHeight

	stopChan := make(chan bool)
	notifyChan := make(chan bool)
	cb := func() {
		log.Log.Info("handle SIGINT for close ord indexer")
		stopChan <- true
		<-notifyChan
		OrdIndexerMgr.Indexer = nil
		receiveChan <- true
	}
	registSigIntFunc(cb)
	go OrdIndexerMgr.Indexer.StartDaemon(stopChan, notifyChan)
	log.Log.Info("ord indexer started")
	return nil
}
