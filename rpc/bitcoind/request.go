package bitcoind

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/OLProtocol/go-bitcoind"
	"github.com/softwarecheng/ord-bridge/common/log"
	"github.com/softwarecheng/ord-bridge/common/util"
)

var BitcoinRpc *bitcoind.Bitcoind

func InitBitconRpc(host string, port int, user, passwd string, useSSL bool) error {
	var err error
	BitcoinRpc, err = bitcoind.New(
		host,
		port,
		user,
		passwd,
		useSSL,
		3600, // server timeout is 1 hour for debug
	)
	return err
}

func IsExistTxInMemPool(txid string) (bool, error) {
	_, err := BitcoinRpc.GetMemPoolEntry(txid)
	if err != nil {
		errNo := strings.Split(err.Error(), ":")[0]
		if errNo == "-5" {
			return false, nil
		}
		return false, nil
	}
	return true, nil
}

func IsExistUtxoInMemPool(utxo string) (bool, error) {
	txid, vout, err := util.ParseUtxo(utxo)
	if err != nil {
		return false, err
	}
	entry, err := BitcoinRpc.GetUnspendTxOutput(txid, vout, true)
	if err != nil {
		return false, err
	}
	return entry.Confirmations == 0, nil
}

func GetBatchUnspendTxOutput(utxoList []string, includeMempool bool, concurrentCount int, ctx context.Context) chan map[string]*bitcoind.UnspendTxOutput {
	if ctx == nil {
		ctx = context.Background()
	}
	unspendTxOutputList := make(map[string]*bitcoind.UnspendTxOutput)
	unspendTxOutputListChan := make(chan map[string]*bitcoind.UnspendTxOutput)

	index := 0
	totalTryCount := 0
	var unspendTxOutputMap sync.Map

	work := func(utxo string) {
		parts := strings.Split(utxo, ":")
		if len(parts) != 2 {
			log.Log.Errorf("GetBatchUnspendTxOutput-> invalid utxo: %s", utxoList[index])
			return
		}
		txid := parts[0]
		outIndex, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Log.Errorf("GetBatchUnspendTxOutput-> invalid utxo: %s", utxoList[index])
			return
		}
		unspendTxOutput, err := BitcoinRpc.GetUnspendTxOutput(txid, outIndex, includeMempool)
		tryCount := 0
		for err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				tryCount++
				totalTryCount++
				log.Log.Debugf("GetBatchUnspendTxOutput-> GetUnspendTxOutput failed: %v, tryCount: %d", err, tryCount)
				unspendTxOutput, err = BitcoinRpc.GetUnspendTxOutput(txid, outIndex, includeMempool)
			}
		}
		unspendTxOutputMap.Store(utxo, unspendTxOutput)
	}

	go func() {
		utxoCount := len(utxoList)
		var runningGoroutines int32
		var wg sync.WaitGroup
		for index < utxoCount {
			select {
			case <-ctx.Done():
				return
			default:
				if atomic.LoadInt32(&runningGoroutines) < int32(concurrentCount) {
					atomic.AddInt32(&runningGoroutines, 1)
					wg.Add(1)
					go func(utxo string) {
						defer wg.Done()
						defer atomic.AddInt32(&runningGoroutines, -1)
						work(utxo)
					}(utxoList[index])
					index++
				}
			}
		}
		wg.Wait()
		unspendTxOutputMap.Range(func(key, value interface{}) bool {
			unspendTxOutputList[key.(string)] = value.(*bitcoind.UnspendTxOutput)
			return true
		})
		unspendTxOutputListChan <- unspendTxOutputList
	}()
	if totalTryCount > 0 {
		log.Log.Debugf("Indexer.GetBatchUnspendTxOutput-> totalTryCount: %d", totalTryCount)
	}
	return unspendTxOutputListChan
}

// 提供一些接口，可以快速同步mempool中的数据，并将数据保存在本地kv数据库
// 1. 启动一个线程，或者一个被动的监听接口，监控内存池的新增tx的信息，
//    需要先获取mempool中所有tx（仅在初始化时调用），并且按照utxo为索引保存在数据库，
//    输入的utxo的spent为true，输出的utxo的spent为false
//    一个utxo很可能在生成后就马上被花费，所以生成时spent为false，被花费时设置为true
//    在上面的基础上，快速获取增量的tx（一般5s调用一次，期望10ms内完成操作）
// 2. 查询接口，查询一个utxo是否已经被花费，数据库查询，代替 IsExistUtxoInMemPool
// 3. 删除接口，删除一个UTXO（该utxo作为输入的tx所在block已经完成）
// 4. 以后可能会有很多基于内存池的操作，比如检查下内存池都是什么类型的tx，是否可以做RBF等等
