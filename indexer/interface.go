package indexer

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/softwarecheng/ord-bridge/common/db/proto3"
	"github.com/softwarecheng/ord-bridge/common/db/raw"
	"github.com/softwarecheng/ord-bridge/common/util"
	"github.com/softwarecheng/ord-bridge/indexer/pb"
	"github.com/softwarecheng/ord-bridge/rpc/ordx"
)

func getInscriptionId(txn *badger.Txn, key string) (string, error) {
	data, err := raw.TxnGet([]byte(key), txn)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", fmt.Errorf("inscriptionId is empty")
	}
	return string(data), nil
}

func getInscriptionIdList(key string, txn *badger.Txn) (map[string]string, error) {
	itemList, err := raw.TxnGetList([]byte(key), txn)
	ret := make(map[string]string)
	if err != nil {
		return nil, err
	}
	for k, v := range itemList {
		ret[k] = string(v)
	}
	return ret, nil
}

func (s *Indexer) fixInscriptionInfo(inscription *pb.Inscription) error {
	if !s.EnableOrdxFix {
		return nil
	}
	if inscription != nil {
		inscription.OrdSat = inscription.Sat
		txid, voutIndex, offset, err := util.ParseOrdSatPoint(inscription.SatPoint)
		if err != nil {
			return fmt.Errorf("ParseOrdSatPoint error: %s, satPoint: %v", err, inscription.SatPoint)
		}
		utxo := fmt.Sprintf("%s:%d", txid, voutIndex)
		resp, err := ordx.Rpc.GetExoticWithUtxo(utxo)
		if err != nil {
			return fmt.Errorf("getExoticWithUtxo error: %v, utox: %s", err, utxo)
		}
		exoticSatRangeUtxo := resp.Data

		if exoticSatRangeUtxo == nil {
			return fmt.Errorf("getExoticWithUtxo is nil, utxo: %s, inscriptionId: %s", utxo, inscription.Id)
		}

		findIndex := -1

		fixOffset := offset
		for index, satDetailInfo := range exoticSatRangeUtxo.Sats {
			if fixOffset < satDetailInfo.Size {
				findIndex = index
				break
			} else {
				fixOffset -= satDetailInfo.Size
			}
		}
		if findIndex == -1 {
			return fmt.Errorf("not find, sat: %d", inscription.Sat)
		} else {
			start := exoticSatRangeUtxo.Sats[findIndex].Start
			inscription.Sat = start + fixOffset
		}

		inscription.Timestamp = inscription.Timestamp * 1000
	}
	return nil
}

func (s *Indexer) getInscriptionList(key string) (inscriptionList map[string]*pb.Inscription, err error) {
	err = s.db.View(func(txn *badger.Txn) error {
		inscriptionIdList, err := getInscriptionIdList(key, txn)
		if err != nil {
			return err
		}
		if len(inscriptionIdList) == 0 {
			return nil
		}
		inscriptionList = make(map[string]*pb.Inscription)
		for _, inscriptionId := range inscriptionIdList {
			key = s.getInscriptionIdDbKey(inscriptionId)
			inscription := new(pb.Inscription)
			err = proto3.TxnGetWithProto3([]byte(key), inscription, txn)
			if err != nil {
				return err
			}
			err = s.fixInscriptionInfo(inscription)
			if err != nil {
				return err
			}
			inscriptionList[inscriptionId] = inscription
		}
		return err
	})
	return inscriptionList, err
}

func (s *Indexer) GetInscription(inscriptionId string) (inscription *pb.Inscription, err error) {
	err = s.db.View(func(txn *badger.Txn) error {
		key := s.getInscriptionIdDbKey(inscriptionId)
		inscription = new(pb.Inscription)
		return proto3.TxnGetWithProto3([]byte(key), inscription, txn)
	})
	if err != nil {
		return nil, err
	}
	err = s.fixInscriptionInfo(inscription)
	if err != nil {
		return nil, err
	}
	return inscription, err
}

func (s *Indexer) GetInscriptionWithNumber(inscriptionNum int64) (inscription *pb.Inscription, err error) {
	err = s.db.View(func(txn *badger.Txn) error {
		key := s.getInscriptionNumberDbKey(inscriptionNum)
		inscriptionId, err := getInscriptionId(txn, key)
		if err != nil {
			return err
		}
		key = s.getInscriptionIdDbKey(inscriptionId)
		inscription = new(pb.Inscription)
		return proto3.TxnGetWithProto3([]byte(key), inscription, txn)
	})
	if err != nil {
		return nil, err
	}
	err = s.fixInscriptionInfo(inscription)
	if err != nil {
		return nil, err
	}
	return inscription, err
}

func (s *Indexer) GetInscriptionWithAddress(address string) (inscriptionList map[string]*pb.Inscription, err error) {
	key := s.getInscriptionAddressDbKey(address)
	return s.getInscriptionList(key)
}

func (s *Indexer) GetInscriptionWithSat(sat int64) (inscriptionList map[string]*pb.Inscription, err error) {
	key := s.getInscriptionSatDbKey(sat)
	return s.getInscriptionList(key)
}

func (s *Indexer) GetInscriptionWithGenesesAddress(genesesAddress string) (inscriptionList map[string]*pb.Inscription, err error) {
	if genesesAddress == "" {
		genesesAddress = NOADDRESS
	}
	key := s.getInscriptionGenesesAddressDbKey(genesesAddress)
	return s.getInscriptionList(key)
}

func (s *Indexer) GetStatus() *pb.Status {
	return s.status
}

func (s *Indexer) GetAddressCount() (uint64, error) {
	var ret uint64 = 0
	key := s.getStatInscriptionAddressDbKey("")
	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek([]byte(key)); it.ValidForPrefix([]byte(key)); it.Next() {
			ret += 1
		}
		return nil
	})
	return ret, err
}

func (s *Indexer) GetGenesesAddressCount() (uint64, error) {
	var ret uint64 = 0
	key := s.getStatInscriptionGenesesAddressDbKey("")
	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek([]byte(key)); it.ValidForPrefix([]byte(key)); it.Next() {
			ret += 1
		}
		return nil
	})
	return ret, err
}
