package ordx

import (
	"encoding/json"

	"github.com/softwarecheng/ord-bridge/common/log"
	"github.com/softwarecheng/ord-bridge/common/request"
)

var Rpc *OrdxRpc

func InitOrdxRpc(rpcUrl string, reqTryTimes int) *OrdxRpc {
	Rpc = &OrdxRpc{
		Url:         rpcUrl,
		ReqTryTimes: reqTryTimes,
	}
	return Rpc
}

func (s *OrdxRpc) GetExoticWithUtxo(utxo string) (*SatRangeResp, error) {
	path := "/exotic/utxo/" + utxo
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return nil, err
	}
	resp := SatRangeResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Error("Error unmarshalling JSON")
		return nil, err
	}
	return &resp, nil
}
