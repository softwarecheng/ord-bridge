package ord

import (
	"github.com/softwarecheng/ord-bridge/rpc/ord"
	"github.com/softwarecheng/ord-bridge/service/common"
)

type InscriptionContentResp struct {
	common.BaseResp
	Data []byte `json:"data"`
}

type OrdStatusResp struct {
	common.BaseResp
	IsOrdx bool        `json:"isOrdx"`
	Data   *ord.Status `json:"data"`
}

type RBlockHashResp struct {
	common.BaseResp
	Data string `json:"data"`
}

type RBlockHeightResp struct {
	common.BaseResp
	Data uint64 `json:"data"`
}

type RBlockInfoResp struct {
	common.BaseResp
	Data *ord.RBlockInfo `json:"data"`
}

type RBlockTimestampResp struct {
	common.BaseResp
	Data uint64 `json:"data"`
}

type RInscriptionIdsPaginationResp struct {
	common.BaseResp
	Data *ord.RInscriptionIdsPagination `json:"data"`
}

type RInscriptionResp struct {
	common.BaseResp
	Data *ord.Inscription `json:"data"`
}

type RMetadataResp struct {
	common.BaseResp
	Data string `json:"data"`
}

type RInscriptionIdResp struct {
	common.BaseResp
	Data *ord.InscriptionId `json:"data"`
}
