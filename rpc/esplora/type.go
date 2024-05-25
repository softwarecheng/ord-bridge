package esplora

type ShareEsploraRpc struct {
	Url         string
	ReqTryTimes int
}

type UTXOStatus struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight uint64 `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   uint64 `json:"block_time"`
}

type UTXO struct {
	Txid   string     `json:"txid"`
	Vout   int        `json:"vout"`
	Status UTXOStatus `json:"status"`
	Value  uint64     `json:"value"`
}

type TxOutput struct {
	ScriptPubkey        string `json:"scriptpubkey"`
	ScriptPubkeyASM     string `json:"scriptpubkey_asm"`
	ScriptPubkeyType    string `json:"scriptpubkey_type"`
	ScriptPubkeyAddress string `json:"scriptpubkey_address"`
	Value               uint64 `json:"value"`
}

type TxInput struct {
	Txid          string   `json:"txid"`
	Vout          int      `json:"vout"`
	Prevout       TxOutput `json:"prevout"`
	ScriptSig     string   `json:"scriptsig"`
	Scriptsig_asm string   `json:"scriptsig_asm"`
	Witness       []string `json:"witness"`
	IsCoinbase    bool     `json:"is_coinbase"`
	Sequence      uint64   `json:"sequence"`
}

type Tx struct {
	Txid     string     `json:"txid"`
	Version  int        `json:"version"`
	Locktime uint64     `json:"locktime"`
	VInput   []TxInput  `json:"vin"`
	Vout     []TxOutput `json:"vout"`
	Size     uint64     `json:"size"`
	Weight   uint64     `json:"weight"`
	Fee      uint64     `json:"fee"`
	Status   UTXOStatus `json:"status"`
}

type BlockHeader struct {
	ID                string  `json:"id"`
	Height            uint64  `json:"height"`
	Version           int     `json:"version"`
	Timestamp         uint64  `json:"timestamp"`
	TxCount           int     `json:"tx_count"`
	Size              uint64  `json:"size"`
	Weight            uint64  `json:"weight"`
	MerkleRoot        string  `json:"merkle_root"`
	PreviousBlockHash string  `json:"previousblockhash"`
	MedianTime        uint64  `json:"mediantime"`
	Nonce             uint32  `json:"nonce"`
	Bits              uint32  `json:"bits"`
	Difficulty        float64 `json:"difficulty"`
}
