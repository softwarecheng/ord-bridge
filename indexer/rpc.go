package indexer

import (
	"time"

	"github.com/softwarecheng/ord-bridge/common/log"
	"github.com/softwarecheng/ord-bridge/rpc/ord"
)

func (s *Indexer) getRpcOrdStatus() *ord.Status {
	resp, err := ord.OrdRpc.GetStatus()
	for err != nil {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			log.Log.Warnf("indexer.getRpcOrdStatus-> failed to get status, error:%v", err)
			time.Sleep(RpcTimeOut)
			resp, err = ord.OrdRpc.GetStatus()
		}
	}
	return resp
}

// func (s *Indexer) getRpcBlockHash(height uint64) string {
// 	resp, err := bitcoin_rpc.ShareBitconRpc.GetBlockHash(height)
// 	for err != nil {
// 		select {
// 		case <-s.ctx.Done():
// 			return ""
// 		default:
// 			log.Log.Warnf("indexer.getRpcBlockHash-> failed to get block hash for block height:%d, error:%v", height, err)
// 			time.Sleep(RpcTimeOut)
// 			resp, err = bitcoin_rpc.ShareBitconRpc.GetBlockHash(height)
// 		}
// 	}
// 	return resp
// }

// func (s *Indexer) getBlock(blockhash string) *bitcoind.Block {
// 	resp, err := bitcoin_rpc.ShareBitconRpc.GetBlock(blockhash)
// 	for err != nil {
// 		select {
// 		case <-s.ctx.Done():
// 			return nil
// 		default:
// 			log.Log.Warnf("indexer.getBlock-> failed to get block for block height:%d, error:%v", blockhash, err)
// 			time.Sleep(RpcTimeOut)
// 			resp, err = bitcoin_rpc.ShareBitconRpc.GetBlock(blockhash)
// 		}
// 	}
// 	return &resp
// }

// func (s *Indexer) getRpcBlockV2(blockhash string) *bitcoind.BlockV2 {
// 	resp, err := bitcoin_rpc.ShareBitconRpc.GetBlockV2(blockhash)
// 	for err != nil {
// 		select {
// 		case <-s.ctx.Done():
// 			return nil
// 		default:
// 			log.Log.Warnf("indexer.getRpcBlockV2-> failed to get block for block height: %s, error: %v", blockhash, err)
// 			time.Sleep(RpcTimeOut)
// 			resp, err = bitcoin_rpc.ShareBitconRpc.GetBlockV2(blockhash)
// 		}
// 	}
// 	return &resp
// }

// bitcoin_rpc.ShareBitconRpc.GetTransaction(blockhash)

// func (s *Indexer) getOrdBlock(height uint64) *ord_rpc.Block {
// 	heightStr := strconv.FormatUint(uint64(height), 10)
// 	resp, err := ord_rpc.ShareOrdinalsRpc.GetBlock(heightStr)
// 	for err != nil {
// 		select {
// 		case <-s.ctx.Done():
// 			return nil
// 		default:
// 			log.Log.Warnf("indexer.getOrdBlock-> failed to get ord block for block height:%d, error:%v", height, err)
// 			time.Sleep(RpcTimeOut)
// 			resp, err = ord_rpc.ShareOrdinalsRpc.GetBlock(heightStr)
// 		}
// 	}
// 	return resp
// }

// func (s *Indexer) getRpcOrdInscription(inscriptionIdOrNum string) *ord_rpc.Inscription {
// 	resp, err := ord_rpc.ShareOrdinalsRpc.GetInscription(inscriptionIdOrNum)
// 	for err != nil {
// 		select {
// 		case <-s.ctx.Done():
// 			return nil
// 		default:
// 			log.Log.Warnf("indexer.getRpcOrdInscription-> failed to get block inscription for inscriptionIdOrNum:%s, error:%v", inscriptionIdOrNum, err)
// 			time.Sleep(RpcTimeOut)
// 			resp, err = ord_rpc.ShareOrdinalsRpc.GetInscription(inscriptionIdOrNum)
// 		}
// 	}
// 	return resp
// }

// func (s *Indexer) getRpcOrdOutput(utxo string) *ord_rpc.Output {
// 	resp, err := ord_rpc.ShareOrdinalsRpc.GetOutput(utxo)
// 	for err != nil {
// 		select {
// 		case <-s.ctx.Done():
// 			return nil
// 		default:
// 			log.Log.Warnf("indexer.getRpcOrdOutput-> error error:%v, utxo:%s", err, utxo)
// 			time.Sleep(RpcTimeOut)
// 			resp, err = ord_rpc.ShareOrdinalsRpc.GetOutput(utxo)
// 		}
// 	}
// 	return resp
// }

func (s *Indexer) getRpcOrdxBlockInscriptions(height uint64) *ord.OrdxBlockInscriptions {
	resp, err := ord.OrdRpc.GetOrdxBlockInscriptionList(height)
	for err != nil {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			log.Log.Warnf("indexer.getRpcOrdxBlockInscriptions-> failed to get block inscription list for block height:%d, error:%v", height, err)
			time.Sleep(RpcTimeOut)
			resp, err = ord.OrdRpc.GetOrdxBlockInscriptionList(height)
		}
	}
	return resp
}

func (s *Indexer) getRpcOrdxBlockTxOutputInscriptions(height uint64) *ord.OrdxBlockTxOutputInscriptions {
	resp, err := ord.OrdRpc.GetOrdxBlockTxOutputInscriptionList(height)
	for err != nil {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			log.Log.Warnf("indexer.getRpcOrdxBlockTxOutputInscriptions-> failed to get block inscription list for block height:%d, error:%v", height, err)
			time.Sleep(RpcTimeOut)
			resp, err = ord.OrdRpc.GetOrdxBlockTxOutputInscriptionList(height)
		}
	}
	return resp
}
