package indexer

import (
	"context"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/softwarecheng/ord-bridge/indexer/pb"
	"github.com/softwarecheng/ord-bridge/rpc/ord"
)

var (
	RpcTimeOut = 3 * time.Second
)

const (
	DBKEY_Stats     = "stats"
	periodFlushToDB = 100
	NOADDRESS       = "UNKNOWN"
)

type Indexer struct {
	db                          *badger.DB
	ordFirstHeight              uint64
	syncFetchDataList           map[uint64]uint64
	fetchOrdInscriptionListChan chan *FetchOrdInscriptinList
	concurrentFetchCount        int
	inscriptionMap              sync.Map
	inscriptionMapCount         uint32
	ctx                         context.Context
	stopIndexerChan             chan bool
	indexerStoppedChan          chan bool
	status                      *pb.Status
	ExitOnReachHeight           bool
	EnableOrdxFix               bool
}

type FetchOrdInscription struct {
	ord.Inscription
	GenesesAddress string `json:"genesesaddress"`
}

type FetchOrdInscriptinList struct {
	InscriptionList []*FetchOrdInscription
	BlockHeight     uint64
}
