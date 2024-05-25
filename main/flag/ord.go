package flag

import (
	"context"

	"github.com/softwarecheng/ord-bridge/common/log"
	"github.com/softwarecheng/ord-bridge/indexer"
	"github.com/softwarecheng/ord-bridge/main/g"
)

func importOrdData(chain, dbDir, importFilePath string) error {
	g.InitRpc()

	const DBSize = 3000 << 20
	g.InitOrdDB(dbDir, DBSize)
	defer g.ReleaseAllDB()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-g.SigInt
		log.Log.Info("Received SIGINT (CTRL+C).")
		cancel()
	}()

	ordIndex, err := indexer.NewIndexer(g.DB.Ord, chain)
	if err != nil {
		return err
	}
	err = ordIndex.ImportData(importFilePath, ctx)
	if err != nil {
		return err
	}
	return err
}
