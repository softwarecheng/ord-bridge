// ConvertTimestampToISO8601 将时间戳转换为 ISO 8601 格式的字符串
package util

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseUtxo(utxo string) (txid string, vout int, err error) {
	parts := strings.Split(utxo, ":")
	if len(parts) != 2 {
		return txid, vout, fmt.Errorf("invalid utxo")
	}

	txid = parts[0]
	vout, err = strconv.Atoi(parts[1])
	if err != nil {
		return txid, vout, err
	}
	if vout < 0 {
		return txid, vout, fmt.Errorf("invalid vout")
	}
	return txid, vout, err
}

func ParseOrdSatPoint(satPoint string) (txid string, outputIndex int, offset int64, err error) {
	parts := strings.Split(satPoint, ":")
	if len(parts) != 3 {
		return "", 0, 0, fmt.Errorf("invalid satPoint")
	}
	txid = parts[0]
	outputIndex, err = strconv.Atoi(parts[1])
	if err != nil {
		return txid, outputIndex, 0, err
	}
	if outputIndex < 0 {
		return txid, outputIndex, 0, fmt.Errorf("invalid index")
	}

	offset, err = strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return txid, outputIndex, offset, err
	}
	return txid, outputIndex, offset, nil
}
