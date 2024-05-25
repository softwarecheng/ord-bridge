package ord

import (
	"testing"
)

var (
	testOrdRpcURL       = "http://192.168.1.102"
	OfficeOrdMainnetUrl = "https://testnet.ordinals.com"
)

func init() {
	InitOrdRpc(testOrdRpcURL, 1)
}

func TestGetInscriptionListInBlock(t *testing.T) {
	blockHeight := uint64(2581379)
	resp, err := OrdRpc.GetInscriptionListInBlock(blockHeight)
	if err != nil {
		t.Fatalf("GetInscriptionListInBlock error: %v", err)
	}
	t.Logf("Response: %+v", resp)
}

func TestGetInscriptionByUtxo(t *testing.T) {
	utxo := "784d96f949dc97543cc23606ca05b5414a484c4eb93a1ea5caa48f9c2a44fedb:0"
	resp, err := OrdRpc.GetInscriptionByUtxo(utxo)
	if err != nil {
		t.Fatalf("GetInscriptionByUtxo error: %v", err)
	}
	t.Logf("Response: %+v", resp)
}

func TestGetStatus(t *testing.T) {
	resp, err := OrdRpc.GetStatus()
	if err != nil {
		t.Fatalf("GetStatus error: %v", err)
	}
	t.Logf("Response: %+v", resp)
}

func TestGetContent(t *testing.T) {
	inscriptionId := "dc032ca45585697ebfed7ba5fc21cbee9e43641f008e7d43656be3e02725e642i0"
	resp, _, err := OrdRpc.GetInscriptionContent(inscriptionId)
	if err != nil {
		t.Fatalf("GetContent error: %v", err)
	}
	t.Logf("Response: %+v", resp)
}
