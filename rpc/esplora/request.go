package esplora

import (
	"encoding/json"
	"strconv"

	"github.com/softwarecheng/ord-bridge/common/request"
)

var EsploraRpc *ShareEsploraRpc

func InitEsploraRpc(rpcUrl string, reqTryTimes int) *ShareEsploraRpc {
	EsploraRpc = &ShareEsploraRpc{
		Url:         rpcUrl,
		ReqTryTimes: reqTryTimes,
	}
	return EsploraRpc
}

// equal to getUtxoByAddr
func (e *ShareEsploraRpc) GetUTXOs(address string) ([]*UTXO, error) {
	path := "/address/" + address + "/utxo"
	body, _, err := request.RpcRequest(e.Url, path, e.ReqTryTimes)
	if err != nil {
		return nil, err
	}

	utxos := make([]*UTXO, 0)
	err = json.Unmarshal(body, &utxos)
	if err != nil {
		return nil, err
	}

	return utxos, nil
}

func (e *ShareEsploraRpc) GetTx(txid string) (*Tx, error) {
	path := "/tx/" + txid
	body, _, err := request.RpcRequest(e.Url, path, e.ReqTryTimes)
	if err != nil {
		return nil, err
	}

	var tx Tx
	err = json.Unmarshal(body, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (e *ShareEsploraRpc) GetTxByAddr(address string) ([]*Tx, error) {
	path := "/address/" + address + "/txs"
	body, _, err := request.RpcRequest(e.Url, path, e.ReqTryTimes)
	if err != nil {
		return nil, err
	}

	var txs []*Tx = make([]*Tx, 0)
	err = json.Unmarshal(body, &txs)
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (e *ShareEsploraRpc) GetBlockHash(height uint64) (string, error) {
	path := "/block-height/" + strconv.FormatUint(height, 10)
	body, _, err := request.RpcRequest(e.Url, path, e.ReqTryTimes)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (e *ShareEsploraRpc) GetRawBlock(blockHash string) (string, error) {
	path := "/block/" + blockHash + "/raw"
	body, _, err := request.RpcRequest(e.Url, path, e.ReqTryTimes)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (e *ShareEsploraRpc) GetBlockHeader(blockHash string) (*BlockHeader, error) {
	path := "/block/" + blockHash
	body, _, err := request.RpcRequest(e.Url, path, e.ReqTryTimes)
	if err != nil {
		return nil, err
	}

	var blockHeader *BlockHeader = &BlockHeader{}
	err = json.Unmarshal(body, &blockHeader)
	if err != nil {
		return nil, err
	}
	return blockHeader, nil
}

func (e *ShareEsploraRpc) GetBlockTxList(blockHash string) ([]*Tx, error) {
	ret := make([]*Tx, 0)
	const pageSize = 25

	blockHeader, err := e.GetBlockHeader(blockHash)
	if err != nil {
		return nil, err
	}

	offset := 0
	for offset < blockHeader.TxCount {
		path := "/block/" + blockHash + "/txs/" + strconv.Itoa(offset)
		body, _, err := request.RpcRequest(e.Url, path, e.ReqTryTimes)
		if err != nil {
			return nil, err
		}
		txListResp := make([]*Tx, 0)
		err = json.Unmarshal(body, &txListResp)
		if err != nil {
			return nil, err
		}
		ret = append(ret, txListResp...)
		offset += pageSize
	}
	return ret, nil
}
