package indexer

import (
	"sort"

	"github.com/dgraph-io/badger/v4"
	"github.com/softwarecheng/ord-bridge/common/db/proto3"
	"github.com/softwarecheng/ord-bridge/common/db/raw"
	"github.com/softwarecheng/ord-bridge/common/log"
	"github.com/softwarecheng/ord-bridge/indexer/pb"
)

var IndexVersion = "0.1.0"
var DbVersion = "0.1.0"

func (s *Indexer) loadStats() error {
	err := s.db.View(func(txn *badger.Txn) error {
		status := &pb.Status{
			IndexVersion:                  IndexVersion,
			DbVersion:                     DbVersion,
			SyncInscriptionHeight:         s.ordFirstHeight - 1,
			SyncTransferInscriptionHeight: s.ordFirstHeight - 1,
			BlessedInscriptions:           0,
			CursedInscriptions:            0,
		}
		err := proto3.TxnGet([]byte(DBKEY_Status), status, txn)
		if err == badger.ErrKeyNotFound {
			log.Log.Info("indexer.loadStatus-> No status found in db")
		} else if err != nil {
			return err
		}
		s.status = status
		return nil
	})
	return err
}

func (s *Indexer) calcSyncInscriptionHeight() (minHeight uint64) {
	minHeight = 0
	syncHeightList := make([]uint64, 0, len(s.syncFetchDataList))
	for k := range s.syncFetchDataList {
		syncHeightList = append(syncHeightList, k)
	}
	if len(syncHeightList) == 1 {
		return syncHeightList[0]
	} else if len(syncHeightList) == 0 {
		return minHeight
	}
	sort.Slice(syncHeightList, func(i, j int) bool {
		return syncHeightList[i] < syncHeightList[j]
	})
	for i := 0; i < len(syncHeightList)-1; i++ {
		if syncHeightList[i+1] != syncHeightList[i]+1 {
			minHeight = syncHeightList[i]
			return minHeight
		}
	}
	minHeight = syncHeightList[len(syncHeightList)-1]
	return minHeight
}

func (s *Indexer) clearSyncFetchDataList(height uint64) {
	for k := range s.syncFetchDataList {
		if k < height {
			delete(s.syncFetchDataList, k)
		}
	}
}

func (s *Indexer) saveStatAddress(key string, wb *badger.WriteBatch) error {
	data, err := raw.Get([]byte(key), s.db)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}
	if data == nil {
		data = make([]byte, 1)
		data[0] = uint8(1)
		err = raw.BatchSet([]byte(key), data, wb)
		if err != nil {
			return err
		}
	}
	return nil
}
