package g

import (
	"os"

	"github.com/softwarecheng/ord-bridge/conf"
	"github.com/softwarecheng/ord-bridge/service"

	ordIndex "github.com/softwarecheng/ord-bridge/indexer"

	"github.com/dgraph-io/badger/v4"
)

type DbList struct {
	Ord *badger.DB
}

type OrdIndexerManager struct {
	Indexer *ordIndex.Indexer
}

var (
	Cfg *conf.YamlConf
	DB  = &DbList{}
)

var (
	SigInt         chan os.Signal
	sigIntFuncList = []func(){}
)

var (
	Rpc           *service.ServiceManager
	OrdIndexerMgr OrdIndexerManager
)
