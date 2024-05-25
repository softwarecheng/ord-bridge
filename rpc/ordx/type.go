package ordx

import "github.com/softwarecheng/ord-bridge/service/common"

type OrdxRpc struct {
	Url         string
	ReqTryTimes int
}

type Satribute string

type SatRange struct {
	Start  int64 `json:"start"`
	Size   int64 `json:"size"`
	Offset int64 `json:"offset"`
}

type SatributeRange struct {
	SatRange
	Satributes []Satribute `json:"satributes"`
}

type SatDetailInfo struct {
	SatributeRange
	Block int   `json:"block"`
	Time  int64 `json:"time"`
}

type ExoticSatRangeUtxo struct {
	Utxo  string          `json:"utxo"`
	Value int64           `json:"value"`
	Sats  []SatDetailInfo `json:"sats"`
}

type SatRangeResp struct {
	common.BaseResp
	Data *ExoticSatRangeUtxo `json:"data"`
}
