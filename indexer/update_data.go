package indexer

import (
	"fmt"
	"runtime"
	"time"

	"github.com/softwarecheng/ord-bridge/common/db/proto3"
	"github.com/softwarecheng/ord-bridge/common/db/raw"
	"github.com/softwarecheng/ord-bridge/common/log"
	"github.com/softwarecheng/ord-bridge/indexer/pb"
)

func (s *Indexer) saveInscriptionList() error {
	if s.inscriptionMapCount == 0 {
		return nil
	}
	startTime := time.Now()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Log.Debugf("indexer.saveInscriptionList-> start stat: program alloc memory: %dMB, sys alloc memory: %dMB, inscription cache count: %d",
		int(float64(m.Alloc)/1024/1024), int(float64(m.Sys)/1024/1024), s.inscriptionMapCount)

	wb := s.db.NewWriteBatch()
	defer wb.Cancel()

	var err error
	inscriptionCount := 0
	cursedCount := uint64(0)
	blessedCount := uint64(0)
	s.inscriptionMap.Range(func(key, value any) bool {
		inscriptionId := key.(string)
		inscription := value.(*pb.Inscription)
		if inscription.Height <= s.status.SyncInscriptionHeight {
			// add inscription by inscriptionId
			key := s.getInscriptionIdDbKey(inscriptionId)
			err = proto3.BatchSet([]byte(key), inscription, wb)
			if err != nil {
				err = fmt.Errorf("proto3.BatchSet error: %s, key: %s, inscriptionId: %s", err, key, inscriptionId)
				log.Log.Error(err)
				return false
			}

			// add inscriptionId by number
			key = s.getInscriptionNumberDbKey(inscription.Number)
			err = raw.BatchSet([]byte(key), []byte(inscriptionId), wb)
			if err != nil {
				err = fmt.Errorf("proto3.BatchSet error: %s, key: %s, inscriptionId: %s", err, key, inscriptionId)
				log.Log.Error(err)
				return false
			}

			// add inscriptionId by sat
			key = s.getInscriptionSatDbKey(inscription.Sat) + "-" + inscriptionId
			err = raw.BatchSet([]byte(key), []byte(inscriptionId), wb)
			if err != nil {
				err = fmt.Errorf("proto3.BatchSet error: %s, key: %s, inscriptionId: %s", err, key, inscriptionId)
				log.Log.Error(err)
				return false
			}

			// add inscriptionId by geneses address and inscriptionId
			key = ""
			if inscription.GenesesAddress != "" {
				key = s.getInscriptionGenesesAddressDbKey(inscription.GenesesAddress) + "-" + inscriptionId
			} else {
				key = s.getInscriptionGenesesAddressDbKey(NOADDRESS) + "-" + inscriptionId
			}
			err = raw.BatchSet([]byte(key), []byte(inscriptionId), wb)
			if err != nil {
				err = fmt.Errorf("proto3.BatchSet error: %s, key: %s, inscriptionId: %s", err, key, inscriptionId)
				log.Log.Error(err)
				return false
			}

			// add inscriptionId by geneses address and count
			if inscription.GenesesAddress != "" {
				key = s.getStatInscriptionGenesesAddressDbKey(inscription.GenesesAddress)
				err = s.saveStatAddress(key, wb)
				if err != nil {
					err = fmt.Errorf("saveStatAddress error: %s, key: %s", err, key)
					log.Log.Error(err)
					return false
				}
			}

			// add address by address and inscriptionId
			if inscription.Address != "" {
				key = s.getInscriptionAddressDbKey(inscription.Address) + "-" + inscriptionId
			} else {
				key = s.getInscriptionAddressDbKey(NOADDRESS) + "-" + inscriptionId
			}
			err = raw.BatchSet([]byte(key), []byte(inscriptionId), wb)
			if err != nil {
				err = fmt.Errorf("proto3.BatchSet error: %s, key: %s, inscriptionId: %s", err, key, inscriptionId)
				log.Log.Error(err)
				return false
			}

			// add address by address and count
			if inscription.Address != "" {
				key = s.getStatInscriptionAddressDbKey(inscription.Address)
				err = s.saveStatAddress(key, wb)
				if err != nil {
					err = fmt.Errorf("saveStatAddress error: %v, key: %s", err, key)
					log.Log.Error(err)
					return false
				}
			}

			if inscription.Number < 0 {
				cursedCount++
			} else {
				blessedCount++
			}
			s.inscriptionMap.Delete(inscriptionId)
			s.inscriptionMapCount--
			inscriptionCount++
		}
		return true
	})

	if err != nil {
		return err
	}

	s.status.CursedInscriptions += cursedCount
	s.status.BlessedInscriptions += blessedCount
	err = proto3.BatchSet([]byte(DBKEY_Stats), s.status, wb)
	if err != nil {
		return fmt.Errorf("setDB error: %v", err)
	}

	log.Log.Infof("indexer.saveInscriptionList-> update inscription count: %d, syncStats:{ SyncHeight: %d, Blessed: %d, Cursed: %d }",
		inscriptionCount, s.status.SyncInscriptionHeight, s.status.BlessedInscriptions, s.status.CursedInscriptions)
	err = wb.Flush()
	if err != nil {
		return fmt.Errorf("wb.Flash error: %v", err)
	}

	runtime.ReadMemStats(&m)
	log.Log.Debugf("indexer.saveInscriptionList-> end stat: program alloc memory: %dMB, sys alloc memory: %dMB, inscription cache count: %d, elapsedTime: %v",
		int(float64(m.Alloc)/1024/1024), int(float64(m.Sys)/1024/1024), s.inscriptionMapCount, time.Since(startTime))
	return nil
}

func (s *Indexer) fetchAndSaveTransferInscriptionList(endHeight uint64) error {
	curHeight := s.status.SyncTransferInscriptionHeight + 1
	fetch2update := func(height uint64) error {
		blockTxOutputInscriptions := s.getRpcOrdxBlockTxOutputInscriptions(height)
		if blockTxOutputInscriptions == nil {
			log.Log.Infof("indexer.fetchAndSaveTransferInscriptionList-> getRpcOrdxBlockTxOutputInscriptions is nil, height: %d", height)
			return nil
		}

		log.Log.Debugf("indexer.fetchAndSaveTransferInscriptionList-> height: %d, need update tx count: %d", height, len(blockTxOutputInscriptions.Inscriptions))
		wb := s.db.NewWriteBatch()
		defer wb.Cancel()
		for _, ordInscription := range blockTxOutputInscriptions.Inscriptions {
			// check inscription in db
			inscription, err := s.GetInscription(ordInscription.Id)
			if err != nil {
				return fmt.Errorf("indexer.fetchAndSaveTransferInscriptionList-> GetInscription error: %s, inscriptionId: %s", err, ordInscription.Id)
			}
			if inscription == nil {
				return fmt.Errorf("indexer.fetchAndSaveTransferInscriptionList-> GetInscription is nil, inscriptionId: %s", ordInscription.Id)
			}

			if inscription.Address != ordInscription.Address {
				key := s.getInscriptionAddressDbKey(ordInscription.Address) + "-" + ordInscription.Id
				err = raw.BatchSet([]byte(key), []byte(inscription.Id), wb)
				if err != nil {
					return fmt.Errorf("raw.BatchSet error: %s, newAddressInscriptionIdKey: %s, inscriptionId: %s",
						err, key, ordInscription.Id)
				}
				key = s.getInscriptionAddressDbKey(inscription.Address) + "-" + ordInscription.Id
				err = wb.Delete([]byte(key))
				if err != nil {
					return fmt.Errorf("wb.Delete error: %s, oldAddressInscriptionIdKey: %s, inscriptionId: %s", err,
						key, ordInscription.Id)
				}

				inscription.Address = ordInscription.Address
				inscription.Children = ordInscription.Children
				if inscription.Fee != ordInscription.Fee {
					log.Log.Debugf("indexer.fetchAndSaveTransferInscriptionList: id:%s, Fee: new: %d, old: %d", ordInscription.Id, ordInscription.Fee, inscription.Fee)
				}
				inscription.Fee = ordInscription.Fee

				if inscription.Next != ordInscription.Next {
					log.Log.Debugf("indexer.fetchAndSaveTransferInscriptionList: id:%s, Next: new: %s, old: %s", ordInscription.Id, ordInscription.Next, inscription.Next)
				}
				inscription.Next = ordInscription.Next

				if inscription.Parent != ordInscription.Parent {
					log.Log.Debugf("indexer.fetchAndSaveTransferInscriptionList: id:%s, Parent: new: %s, old: %s", ordInscription.Id, ordInscription.Parent, inscription.Parent)
				}
				inscription.Parent = ordInscription.Parent

				if inscription.Previous != ordInscription.Previous {
					log.Log.Debugf("indexer.fetchAndSaveTransferInscriptionList: id:%s, Previous: new: %s, old: %s", ordInscription.Id, ordInscription.Previous, inscription.Previous)
				}
				inscription.Previous = ordInscription.Previous

				if inscription.Sat != ordInscription.Sat {
					log.Log.Debugf("indexer.fetchAndSaveTransferInscriptionList: id:%s, Sat: new: %d, old: %d", ordInscription.Id, ordInscription.Sat, inscription.Sat)
				}
				inscription.Sat = ordInscription.Sat

				if inscription.SatPoint != ordInscription.SatPoint {
					log.Log.Debugf("indexer.fetchAndSaveTransferInscriptionList: id:%s, SatPoint: new: %s, old: %s", ordInscription.Id, ordInscription.SatPoint, inscription.SatPoint)
				}
				inscription.SatPoint = ordInscription.SatPoint

				if inscription.Timestamp != ordInscription.Timestamp {
					log.Log.Debugf("indexer.fetchAndSaveTransferInscriptionList: id:%s, Timestamp: new: %d, old: %d", ordInscription.Id, ordInscription.Timestamp, inscription.Timestamp)
				}
				inscription.Timestamp = ordInscription.Timestamp

				if inscription.Value != ordInscription.Value {
					log.Log.Debugf("indexer.fetchAndSaveTransferInscriptionList: id:%s, Value: new: %d, old: %d", ordInscription.Id, ordInscription.Value, inscription.Value)
				}
				inscription.Value = ordInscription.Value

				key = s.getStatInscriptionAddressDbKey(inscription.Address)
				err = s.saveStatAddress(key, wb)
				if err != nil {
					return fmt.Errorf("saveStatAddress error: %s", err)
				}

				key = s.getInscriptionIdDbKey(inscription.Id)
				err = proto3.BatchSet([]byte(key), inscription, wb)
				if err != nil {
					return fmt.Errorf("proto3.BatchSet error: %s, inscriptionId: %s", err, ordInscription.Id)
				}
			}
		}

		s.status.SyncTransferInscriptionHeight = curHeight
		err := proto3.BatchSet([]byte(DBKEY_Stats), s.status, wb)
		if err != nil {
			err = fmt.Errorf("proto3.BatchSet error: %s", err)
			return err
		}

		err = wb.Flush()
		if err != nil {
			err = fmt.Errorf("wb.Flash error: %s", err)
			log.Log.Error(err)
			return err
		}
		log.Log.Infof("indexer.fetchAndSaveTransferInscriptionList-> flush height: %d, inscription count: %d", curHeight, len(blockTxOutputInscriptions.Inscriptions))
		return nil
	}

	for curHeight <= endHeight {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		default:
			err := fetch2update(curHeight)
			if err != nil {
				return err
			}
			curHeight++
		}
	}
	return nil
}
