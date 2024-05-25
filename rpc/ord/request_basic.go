package ord

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/softwarecheng/ord-bridge/common/log"
	"github.com/softwarecheng/ord-bridge/common/request"
)

var OrdRpc *ShareOrdRpc

func InitOrdRpc(rpcUrl string, reqTryTimes int) *ShareOrdRpc {
	OrdRpc = &ShareOrdRpc{
		Url:         rpcUrl,
		ReqTryTimes: reqTryTimes,
	}
	return OrdRpc
}

func (s *ShareOrdRpc) GetInscription(inscriptionIdOrNum string) (*Inscription, error) {
	path := "/inscription/" + inscriptionIdOrNum
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return nil, err
	}
	var resp Inscription
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Error("Error unmarshalling JSON")
		return nil, err
	}
	return &resp, nil
}

func (s *ShareOrdRpc) GetOutput(utxo string) (*Output, error) {
	path := "/output/" + utxo
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return nil, err
	}
	var resp Output
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Error("Error unmarshalling JSON")
		return nil, err
	}
	return &resp, nil
}

func (s *ShareOrdRpc) GetTx(txid string) (*Tx, error) {
	path := "/tx/" + txid
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return nil, err
	}
	var resp Tx
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Error("Error unmarshalling JSON")
		return nil, err
	}
	return &resp, nil
}

func (s *ShareOrdRpc) GetSat(satNumber string) (*Sat, error) {
	path := "/sat/" + satNumber
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return nil, err
	}
	var resp Sat
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Error("Error unmarshalling JSON")
		return nil, err
	}
	return &resp, nil
}

func (s *ShareOrdRpc) GetBlock(heightOrHash string) (*Block, error) {
	path := "/block/" + heightOrHash
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return nil, err
	}
	var resp Block
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Error("Error unmarshalling JSON")
		return nil, err
	}

	return &resp, nil
}

func (s *ShareOrdRpc) GetStatus() (*Status, error) {
	path := "/status"
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return nil, err
	}
	var resp Status
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Error("Error unmarshalling JSON")
		return nil, err
	}
	return &resp, nil
}

// the content of the inscription with inscriptionId
func (s *ShareOrdRpc) GetInscriptionContent(inscriptionId string) ([]byte, http.Header, error) {
	path := "/content/" + inscriptionId
	body, header, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return nil, nil, err
	}
	return body, header, nil
}

// the content of the inscription with inscriptionId
func (s *ShareOrdRpc) GetInscriptionPreview(inscriptionId string) ([]byte, http.Header, error) {
	path := "/preview/" + inscriptionId
	body, header, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return nil, nil, err
	}
	return body, header, nil
}

// Recursion: get block hash at given block height
func (s *ShareOrdRpc) GetRBlockHash(height uint64) (string, error) {
	path := "/r/blockhash/" + strconv.FormatUint(height, 10)
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return "", err
	}

	var resp string
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Warn("Error unmarshalling JSON")
		return "", err
	}
	return resp, nil
}

// Recursion: get latest block hash
func (s *ShareOrdRpc) GetRLastestBlockHash() (string, error) {
	path := "/r/blockhash"
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return "", err
	}
	var resp string
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Warn("Error unmarshalling JSON")
		return "", err
	}

	return resp, nil
}

// Recursion: get latest block height
func (s *ShareOrdRpc) GetRLastestBlockHeight() (uint64, error) {
	path := "/r/blockheight"
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return 0, err
	}
	height, err := strconv.ParseUint(string(body), 10, 64)
	if err != nil {
		return 0, err
	}
	return height, nil
}

// Recursion: latest block info at given block height or hash
func (s *ShareOrdRpc) GetRBlockInfo(heightOrHash string) (*RBlockInfo, error) {
	path := "/r/blockinfo/" + heightOrHash
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return nil, err
	}
	var resp RBlockInfo
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Warn("Error unmarshalling JSON")
		return nil, err
	}
	return &resp, nil
}

// Recursion: latest block timestamp
func (s *ShareOrdRpc) GetRLatestBlockTimestamp() (uint64, error) {
	path := "/r/blocktime"
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return 0, err
	}
	height, err := strconv.ParseUint(string(body), 10, 64)
	if err != nil {
		return 0, err
	}
	return height, nil
}

// Recursion: get children of inscription id list
func (s *ShareOrdRpc) GetRChildrenInscriptionIdList(inscritipnId string, page string) (*RInscriptionIdsPagination, error) {
	path := "/r/children/" + inscritipnId
	if page != "" {
		path += "/" + page
	}
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return nil, err
	}
	var resp RInscriptionIdsPagination
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Warn("Error unmarshalling JSON")
		return nil, err
	}
	return &resp, nil
}

// Recursion: get inscription info
func (s *ShareOrdRpc) GetRInscriptionInfo(inscritipnId string) (*Inscription, error) {
	path := "/r/inscription/" + inscritipnId
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return nil, err
	}
	var resp Inscription
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Warn("Error unmarshalling JSON")
		return nil, err
	}
	return &resp, nil
}

// Recursion: get inscription info
func (s *ShareOrdRpc) GetRMetaData(inscritipnId string) (string, error) {
	path := "/r/metadata/" + inscritipnId
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return "", err
	}
	var resp string
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Warn("Error unmarshalling JSON")
		return "", err
	}
	return resp, nil
}

// Recursion: get inscriptions ids for given sat number
func (s *ShareOrdRpc) GetRSatInscriptionIdList(satNumber string, page string) (*RInscriptionIdsPagination, error) {
	path := "/r/sat/" + satNumber
	if page != "" {
		path += "/" + page
	}
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return nil, err
	}
	var resp RInscriptionIdsPagination
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Warn("Error unmarshalling JSON")
		return nil, err
	}
	return &resp, nil
}

// Recursion: get inscriptions id for given sat number
func (s *ShareOrdRpc) GetRSatInscriptionId(satNumber string, index string) (*InscriptionId, error) {
	path := "/r/sat/" + satNumber
	if index != "" {
		path += "/at/" + index
	}
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return nil, err
	}
	var resp InscriptionId
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Warn("Error unmarshalling JSON")
		return nil, err
	}
	return &resp, nil
}

// Ordx customized: get inscription and output data collection with block height
func (s *ShareOrdRpc) GetOrdxBlockInscriptionList(height uint64) (*OrdxBlockInscriptions, error) {
	path := "/ordx/block/inscriptions/" + strconv.FormatUint(height, 10)
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return nil, err
	}
	var resp OrdxBlockInscriptions
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Warn("Error unmarshalling JSON")
		return nil, err
	}
	return &resp, nil
}

func (s *ShareOrdRpc) GetOrdxBlockTxOutputInscriptionList(height uint64) (*OrdxBlockTxOutputInscriptions, error) {
	path := "/ordx/block/tx/outputs/inscriptions/" + strconv.FormatUint(height, 10)
	body, _, err := request.RpcRequest(s.Url, path, s.ReqTryTimes)
	if err != nil {
		return nil, err
	}
	var resp OrdxBlockTxOutputInscriptions
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Log.WithField("body", string(body)).Warn("Error unmarshalling JSON")
		return nil, err
	}
	return &resp, nil
}
