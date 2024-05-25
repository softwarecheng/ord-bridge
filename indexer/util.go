package indexer

import (
	"fmt"
	"strconv"
	"strings"
)

func GetInscriptionId(mintutxo string, id int) string {
	parts := strings.Split(mintutxo, ":")
	idstr := strconv.Itoa(int(id))

	return parts[0] + "i" + idstr
}

func ParseInscriptionId(inscId string) (string, int, error) {
	parts := strings.Split(inscId, "i")
	if len(parts) != 2 {
		return inscId, 0, fmt.Errorf("wrong format %s", inscId)
	}

	i, err := strconv.Atoi(parts[1])
	if err != nil {
		return inscId, 0, err
	}

	return parts[0], (i), nil
}

func ParseSatPoint(satpoint string) (string, error) {
	parts := strings.Split(satpoint, ":")
	if len(parts) != 3 {
		return satpoint, fmt.Errorf("wrong format %s", satpoint)
	}

	return parts[0] + parts[1], nil
}
