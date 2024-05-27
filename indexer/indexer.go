package indexer

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/softwarecheng/ord-bridge/common/db/proto3"
	"github.com/softwarecheng/ord-bridge/common/define"
	"github.com/softwarecheng/ord-bridge/common/log"
	"github.com/softwarecheng/ord-bridge/indexer/pb"
)

func NewIndexer(db *badger.DB, chain string) (*Indexer, error) {
	ordFirstHeight := uint64(define.HEIGHT_TESTNET_ORD)
	switch chain {
	case "testnet":
		ordFirstHeight = define.HEIGHT_TESTNET_ORD
	case "mainnet":
		ordFirstHeight = define.HEIGHT_MAINNET_ORD
	default:
		return nil, fmt.Errorf("unsupported chain: %s", chain)
	}

	concurrentCount := runtime.NumCPU()
	ret := &Indexer{
		db:                          db,
		ordFirstHeight:              ordFirstHeight,
		concurrentFetchCount:        concurrentCount,
		fetchOrdInscriptionListChan: make(chan *FetchOrdInscriptinList, concurrentCount),
		syncFetchDataList:           make(map[uint64]uint64),
		ctx:                         context.Background(),
		indexerStoppedChan:          make(chan bool),
		stopIndexerChan:             make(chan bool),
	}

	return ret, nil
}

func (s *Indexer) StartDaemon(stopChan chan bool, notifyChan chan bool) {
	isRunning := false
	const n = 10
	ticker := time.NewTicker(time.Duration(n) * time.Second)
	tick := func() {
		isRunning = true
		err := s.run()
		isRunning = false
		if err != nil {
			log.Log.Infof("indexer.StartDaemon-> Run error: %v", err)
			return
		}
	}

	go tick()
	isSendSignal := false
	for {
		select {
		case <-ticker.C:
			if !isRunning {
				if s.ExitOnReachHeight {
					if !isSendSignal {
						isSendSignal = true
						pid := os.Getpid()
						process, err := os.FindProcess(pid)
						if err != nil {
							log.Log.Errorf("indexer.StartDaemon-> FindProcess error: %v", err)
							continue
						}
						err = process.Signal(syscall.SIGINT)
						if err != nil {
							log.Log.Errorf("indexer.StartDaemon-> signal error: %v", err)
							continue
						}
						log.Log.Info("indexer.StartDaemon-> exit on reach height")
						continue
					} else {
						return
					}
				}
				go tick()
			}
		case <-stopChan:
			if isRunning {
				isRunning = false
				log.Log.Info("indexer.StartDaemon-> receive external stop notify signal, send internal stop signal")
				s.stopIndexerChan <- true
				log.Log.Info("indexer.StartDaemon-> receive internal stop signal, wait internal stoped signal")
				<-s.indexerStoppedChan
				log.Log.Info("indexer.StartDaemon-> receive internal stoped signal, send external stoped notify signal")
			}
			notifyChan <- true
			return
		}
	}
}

func (s *Indexer) run() error {
	statusResp := s.getRpcOrdStatus()
	if statusResp == nil {
		return fmt.Errorf("indexer.Run-> getOrdStatus error")
	}
	err := s.loadStats()
	if err != nil {
		return err
	}

	err = s.syncData(statusResp.Height)
	if err != nil {
		return err
	}
	return nil
}

func (s *Indexer) syncData(maxHeight uint64) error {
	if s.status.SyncTransferInscriptionHeight < s.status.SyncInscriptionHeight {
		err := s.fetchAndSaveTransferInscriptionList(s.status.SyncInscriptionHeight)
		if err != nil {
			return err
		}
	}
	if s.status.SyncInscriptionHeight >= maxHeight {
		log.Log.Infof("indexer.syncData-> already synced to block %d", maxHeight)
		return nil
	}

	startHeight := s.status.SyncInscriptionHeight + 1
	log.Log.Infof("indexer.syncData-> start handle all blocks: from %d to %d, startTime: %v",
		startHeight, maxHeight, time.Now().Format("2006-01-02 15:04:05"))

	go s.fetchData(startHeight, maxHeight)

	for s.status.SyncInscriptionHeight < maxHeight {
		select {
		case <-s.stopIndexerChan:
			err := s.saveInscriptionList()
			if err != nil {
				return err
			}
			log.Log.Info("indexer.syncData-> Graceful close service.")
			s.indexerStoppedChan <- true
			return nil
		default:
			fetchData := <-s.fetchOrdInscriptionListChan
			if fetchData == nil {
				log.Log.Warnf("indexer.syncData-> fetchDataChan is nil")
				return nil
			}
			s.syncFetchDataList[fetchData.BlockHeight] = fetchData.BlockHeight
			inscriptionDataList := fetchData.InscriptionList
			if len(inscriptionDataList) == 0 {
				s.status.SyncInscriptionHeight = s.calcSyncInscriptionHeight()
				err := proto3.Set([]byte(DBKEY_Status), s.status, s.db)
				if err != nil {
					return err
				}
				continue
			}

			for _, inscriptionData := range inscriptionDataList {
				inscription := &pb.Inscription{
					GenesesAddress: inscriptionData.GenesesAddress,
					Address:        inscriptionData.Address,
					Children:       inscriptionData.Children,
					ContentLength:  uint32(inscriptionData.ContentLength),
					ContentType:    inscriptionData.ContentType,
					Fee:            inscriptionData.Fee,
					Height:         inscriptionData.Height,
					Id:             inscriptionData.Id,
					Next:           inscriptionData.Next,
					Number:         inscriptionData.Number,
					Parent:         inscriptionData.Parent,
					Previous:       inscriptionData.Previous,
					Sat:            inscriptionData.Sat,
					SatPoint:       inscriptionData.SatPoint,
					Timestamp:      inscriptionData.Timestamp,
					Value:          inscriptionData.Value,
				}
				s.inscriptionMap.Store(inscriptionData.Id, inscription)
				s.inscriptionMapCount++
			}
			s.status.SyncInscriptionHeight = s.calcSyncInscriptionHeight()
			log.Log.Infof("indexer.syncData-> fetch block data for height %d, calculated sync height: %d",
				fetchData.BlockHeight, s.status.SyncInscriptionHeight)

			if (maxHeight - s.status.SyncInscriptionHeight) == 0 {
				s.syncFetchDataList = make(map[uint64]uint64)
				err := s.saveInscriptionList()
				if err != nil {
					return err
				}
				return nil
			}

			if len(s.syncFetchDataList) > periodFlushToDB {
				err := s.saveInscriptionList()
				if err != nil {
					return err
				}
				s.clearSyncFetchDataList(s.status.SyncInscriptionHeight)
			}
		}
	}
	return nil
}

func (s *Indexer) fetchData(startHeight, endHeight uint64) {
	curHeight := startHeight
	var runningGoroutines int32
	var wg sync.WaitGroup

	work := func(height uint64) {
		blockInscriptions := s.getRpcOrdxBlockInscriptions(height)
		if len(blockInscriptions.Inscriptions) == 0 {
			s.fetchOrdInscriptionListChan <- &FetchOrdInscriptinList{InscriptionList: nil, BlockHeight: height}
			return
		}
		inscriptionDataList := make([]*FetchOrdInscription, 0)
		for _, ordxBlockInscription := range blockInscriptions.Inscriptions {
			inscription := ordxBlockInscription.Inscription
			ordInscription := &FetchOrdInscription{
				GenesesAddress: ordxBlockInscription.GenesesAddress,
				Inscription:    inscription,
			}
			inscriptionDataList = append(inscriptionDataList, ordInscription)
		}
		s.fetchOrdInscriptionListChan <- &FetchOrdInscriptinList{InscriptionList: inscriptionDataList, BlockHeight: height}
	}

	for curHeight <= endHeight {
		select {
		case <-s.ctx.Done():
			return
		default:
			if atomic.LoadInt32(&runningGoroutines) < int32(s.concurrentFetchCount) {
				atomic.AddInt32(&runningGoroutines, 1)
				wg.Add(1)
				go func(height uint64) {
					defer atomic.AddInt32(&runningGoroutines, -1)
					defer wg.Done()
					work(height)
				}(curHeight)
				curHeight++
			}
		}
	}
	wg.Wait()
}
