package index

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/softwarecheng/ord-bridge/common/util"
	"github.com/softwarecheng/ord-bridge/indexer"
	"github.com/softwarecheng/ord-bridge/indexer/pb"
)

func (s *Model) getInscriptionAssetList(inscriptionList map[string]*pb.Inscription, start, limit int) ([]*OrdInscriptionAsset, int, error) {
	utxosort := make([]*OrdInscriptionAsset, 0)
	for _, inscription := range inscriptionList {
		txid, outputIndex, _, err := util.ParseOrdSatPoint(inscription.SatPoint)
		if err != nil {
			return nil, 0, err
		}
		utxo := txid + ":" + fmt.Sprintf("%d", outputIndex)
		utxosort = append(utxosort, &OrdInscriptionAsset{utxo, inscription})
	}
	sort.Slice(utxosort, func(i, j int) bool {
		if utxosort[i].Value == utxosort[j].Value {
			return utxosort[i].Utxo > utxosort[j].Utxo
		} else {
			return utxosort[i].Value > utxosort[j].Value
		}
	})

	totalRecords := len(utxosort)
	if totalRecords < start {
		return nil, 0, nil
	}
	if totalRecords < start+limit {
		limit = totalRecords - start
	}
	end := start + limit
	utxoresult := utxosort[start:end]

	for i := 0; i < len(utxoresult); i++ {
		if utxoresult[i].GenesesAddress == "" {
			utxoresult[i].GenesesAddress = indexer.NOADDRESS
		}
		if utxoresult[i].Address == "" {
			utxoresult[i].Address = indexer.NOADDRESS
		}
	}
	return utxoresult, totalRecords, nil
}

func (s *Model) GetInscriptionAssetList(address string, isGenesesAddr bool, start, limit int) ([]*OrdInscriptionAsset, int, error) {
	var inscriptionList map[string]*pb.Inscription
	var err error

	if isGenesesAddr {
		inscriptionList, err = s.indexer.GetInscriptionWithGenesesAddress(address)
	} else {
		inscriptionList, err = s.indexer.GetInscriptionWithAddress(address)
	}
	if err != nil {
		return nil, 0, err
	}
	return s.getInscriptionAssetList(inscriptionList, start, limit)
}

func (s *Model) GetOrdInscriptionWithSat(query string, start, limit int) ([]*OrdInscriptionAsset, int, error) {
	sat, err := strconv.ParseInt(query, 10, 64)
	if err != nil {
		return nil, 0, err
	}
	inscriptionList, err := s.indexer.GetInscriptionWithSat(sat)
	if err != nil {
		return nil, 0, err
	}

	return s.getInscriptionAssetList(inscriptionList, start, limit)
}

func (s *Model) GetOrdInscriptionWithId(query string) (*OrdInscriptionAsset, error) {
	inscription, err := s.indexer.GetInscription(query)
	if err != nil {
		return nil, err
	}
	txid, outputIndex, _, err := util.ParseOrdSatPoint(inscription.SatPoint)
	if err != nil {
		return nil, err
	}
	utxo := txid + ":" + fmt.Sprintf("%d", outputIndex)
	ret := &OrdInscriptionAsset{utxo, inscription}
	return ret, nil
}

func (s *Model) GetOrdInscriptionWithNumer(query string) (*OrdInscriptionAsset, error) {
	inscriptionNumber, _ := strconv.ParseInt(query, 10, 64)
	inscription, err := s.indexer.GetInscriptionWithNumber(inscriptionNumber)
	if err != nil {
		return nil, err
	}
	txid, outputIndex, _, err := util.ParseOrdSatPoint(inscription.SatPoint)
	if err != nil {
		return nil, err
	}
	utxo := txid + ":" + fmt.Sprintf("%d", outputIndex)
	ret := &OrdInscriptionAsset{utxo, inscription}
	return ret, nil
}
