package index

import (
	"github.com/softwarecheng/ord-bridge/service/common"
)

type HealthResp struct {
	Status string `json:"status"`
}

type IndexStatusResp struct {
	IndexVersion                  string `json:"indexVersion"`
	DbVersion                     string `json:"dbVersion"`
	SyncInscriptionHeight         uint64 `json:"syncInscriptionHeight"`
	SyncTransferInscriptionHeight uint64 `json:"syncTransferInscriptionHeight"`
	BlessedInscriptions           uint64 `json:"blessedInscriptions"`
	CursedInscriptions            uint64 `json:"cursedInscriptions"`
	AddressCount                  uint64 `json:"addressCount"`
	GenesesAddressCount           uint64 `json:"genesesAddressCount"`
}

type InscriptionContentResp struct {
	common.BaseResp
	Data []byte `json:"data"`
}

type OrdInscriptionListData struct {
	common.ListResp
	Detail []*OrdInscriptionAsset `json:"detail"`
}

type OrdInscriptionListResp struct {
	common.BaseResp
	Data *OrdInscriptionListData `json:"data"`
}

type OrdInscriptionResp struct {
	common.BaseResp
	Data *OrdInscriptionAsset `json:"data"`
}
