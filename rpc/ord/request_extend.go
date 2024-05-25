package ord

import (
	"fmt"
	"strconv"
	"strings"
)

// check utxo whether has ord asset. key: insciption_id
func (s *ShareOrdRpc) GetInscriptionByUtxo(utxo string) (map[string]int64, error) {
	resp, err := s.GetOutput(utxo)
	if err != nil {
		return nil, err
	}
	inscriptions := make(map[string]int64)
	for _, inscid := range resp.Inscriptions {
		inscriptions[inscid] = 0
	}
	return inscriptions, nil
}

// key: tx; value: id
func (s *ShareOrdRpc) GetInscriptionListInBlock(height uint64) (map[string][]int, error) {
	resp, err := s.GetBlock(strconv.FormatUint(height, 10))
	if err != nil {
		return nil, err
	}

	inscriptions := make(map[string][]int)
	for _, insc_id := range resp.Inscriptions {
		tx, id, err := parseInscriptionId(insc_id)
		if err != nil {
			continue
		}

		inscriptions[tx] = append(inscriptions[tx], id)
	}
	return inscriptions, nil
}

func parseInscriptionId(inscid string) (string, int, error) {
	parts := strings.Split(inscid, "i")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid inscription_id format %s", inscid)
	}

	id, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return "", 0, err
	}
	return parts[0], int(id), nil
}
