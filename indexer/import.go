package indexer

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/softwarecheng/ord-bridge/common/log"
	"github.com/softwarecheng/ord-bridge/indexer/pb"
	ord_rpc "github.com/softwarecheng/ord-bridge/rpc/ord"
)

func (s *Indexer) ImportData(filePath string, ctx context.Context) error {
	err := s.loadStats()
	if err != nil {
		return err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	totalCount := uint64(0)
	needFlushCount := uint64(0)
	startTime := time.Now()

	scanner := bufio.NewScanner(file)
	const maxCapacity = 1024 * 1024 * 1024 * 5 // 50MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	blockInscriptions := ord_rpc.OrdxBlockInscriptions{}
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			err = s.saveInscriptionList()
			if err != nil {
				return err
			}
			log.Log.Infof("indexer.ImportData-> cancelled by user, finish stats: totalCount: %d, latest height: %d, syncStats: %+v",
				totalCount, blockInscriptions.Height, s.status)
			return nil
		default:
			line := scanner.Bytes()
			if len(line) == 0 || line[0] == '\n' {
				continue
			}
			err = json.Unmarshal(line, &blockInscriptions)
			if err != nil {
				fileName := "errline.txt"
				err := os.WriteFile(fileName, line, 0644)
				if err != nil {
					fmt.Printf("Error writing to %s: %v", fileName, err)
					os.Exit(1)
				}
				return err
			}
			if blockInscriptions.Height <= s.status.SyncInscriptionHeight || len(blockInscriptions.Inscriptions) == 0 {
				log.Log.Infof("indexer.ImportData-> skip block %d", blockInscriptions.Height)
				continue
			}
			log.Log.Infof("indexer.ImportData-> processing block %d", blockInscriptions.Height)

			for _, ordxBlockInscription := range blockInscriptions.Inscriptions {
				inscription := &pb.Inscription{
					GenesesAddress: ordxBlockInscription.GenesesAddress,
					Address:        ordxBlockInscription.Inscription.Address,
					Children:       ordxBlockInscription.Inscription.Children,
					ContentLength:  uint32(ordxBlockInscription.Inscription.ContentLength),
					ContentType:    ordxBlockInscription.Inscription.ContentType,
					Fee:            ordxBlockInscription.Inscription.Fee,
					Height:         ordxBlockInscription.Inscription.Height,
					Id:             ordxBlockInscription.Inscription.Id,
					Next:           ordxBlockInscription.Inscription.Next,
					Number:         ordxBlockInscription.Inscription.Number,
					Parent:         ordxBlockInscription.Inscription.Parent,
					Previous:       ordxBlockInscription.Inscription.Previous,
					Sat:            ordxBlockInscription.Inscription.Sat,
					SatPoint:       ordxBlockInscription.Inscription.SatPoint,
					Timestamp:      ordxBlockInscription.Inscription.Timestamp,
					Value:          ordxBlockInscription.Inscription.Value,
				}
				s.inscriptionMap.Store(inscription.Id, inscription)
				s.inscriptionMapCount++
				needFlushCount++
			}

			s.status.SyncInscriptionHeight = blockInscriptions.Height
			s.status.SyncTransferInscriptionHeight = s.status.SyncInscriptionHeight
			if needFlushCount%1000 == 0 {
				err = s.saveInscriptionList()
				if err != nil {
					return err
				}
				totalCount += needFlushCount
				needFlushCount = 0
			}
		}
	}

	err = s.saveInscriptionList()
	if err != nil {
		return err
	}

	log.Log.Infof("indexer.ImportData-> totalCount: %d, latest height: %d, syncStats: %+v",
		totalCount, blockInscriptions.Height, s.status)

	if err = scanner.Err(); err != nil {
		return err
	}

	log.Log.Infof("ordinals.Indexer.ImportData-> success import ordx inscriptions data, syncStats: %+v, cost: %v",
		s.status, time.Since(startTime))
	return nil
}
