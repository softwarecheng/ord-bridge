package indexer

import (
	"fmt"
)

func (b *Indexer) getInscriptionIdDbKey(inscriptionId string) string {
	return fmt.Sprintf("insc-id-%s", inscriptionId)
}

func (b *Indexer) getInscriptionNumberDbKey(inscriptionNumber int64) string {
	return fmt.Sprintf("insc-num-%d", inscriptionNumber)
}

func (b *Indexer) getInscriptionSatDbKey(sat int64) string {
	return fmt.Sprintf("insc-sat-%d", sat)
}

func (b *Indexer) getInscriptionGenesesAddressDbKey(address string) string {
	return fmt.Sprintf("insc-genesesaddr-%s", address)
}

func (b *Indexer) getInscriptionAddressDbKey(address string) string {
	return fmt.Sprintf("insc-addr-%s", address)
}

func (b *Indexer) getStatInscriptionGenesesAddressDbKey(address string) string {
	return fmt.Sprintf("insc-stat-genesesaddr-%s", address)
}

func (b *Indexer) getStatInscriptionAddressDbKey(address string) string {
	return fmt.Sprintf("insc-stat-addr-%s", address)
}
