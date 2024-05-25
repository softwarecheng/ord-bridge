package index

import "github.com/softwarecheng/ord-bridge/indexer/pb"

type OrdInscriptionAsset struct {
	Utxo            string `json:"utxo"`
	*pb.Inscription `json:"inscription"`
}
