package ord

type ShareOrdRpc struct {
	Url         string
	ReqTryTimes int
}

type Inscription struct {
	Address       string   `json:"address"`
	Children      []string `json:"children"`
	ContentLength uint     `json:"content_length"`
	ContentType   string   `json:"content_type"`
	Fee           uint64   `json:"fee"`
	Height        uint64   `json:"height"`
	Id            string   `json:"id"`
	Next          string   `json:"next"`
	Number        int64    `json:"number"`
	Parent        string   `json:"parent"`
	Previous      string   `json:"previous"`
	Sat           int64    `json:"sat"`
	SatPoint      string   `json:"satpoint"`
	Timestamp     int64    `json:"timestamp"`
	Value         uint64   `json:"value"`
}

type Tx struct {
	Chain            string      `json:"chain"`
	Etching          interface{} `json:"etching"`
	InscriptionCount uint32      `json:"inscription_count"`
	Transaction      struct {
		Version  uint32 `json:"version"`
		LockTime uint64 `json:"lock_time"`
		Input    []struct {
			PreviousOutput string   `json:"previous_output"`
			ScriptSig      string   `json:"script_sig"`
			Sequence       uint64   `json:"sequence"`
			Witness        []string `json:"witness"`
		} `json:"input"`
		Output []struct {
			Value        uint64 `json:"value"`
			ScriptPubkey string `json:"script_pubkey"`
		} `json:"output"`
	} `json:"transaction"`
	Txid string `json:"txid"`
}

type Output struct {
	Address      string      `json:"address"`
	Indexed      bool        `json:"indexed"`
	Inscriptions []string    `json:"inscriptions"`
	SatRanges    [][2]uint64 `json:"sat_ranges"`
	ScriptPubkey string      `json:"script_pubkey"`
	Spent        bool        `json:"spent"`
	Transaction  string      `json:"transaction"`
	Value        uint64      `json:"value"`
}

type OrdxOutput struct {
	Address     string `json:"address"`
	Transaction string `json:"transaction"`
	Value       uint64 `json:"value"`
}

type Sat struct {
	Number       int64    `json:"number"`
	Decimal      string   `json:"decimal"`
	Degree       string   `json:"degree"`
	Name         string   `json:"name"`
	Block        uint32   `json:"block"`
	Cycle        uint32   `json:"cycle"`
	Epoch        uint32   `json:"epoch"`
	Period       uint32   `json:"period"`
	Offset       uint64   `json:"offset"`
	Rarity       string   `json:"rarity"`
	Percentile   string   `json:"percentile"`
	Satpoint     string   `json:"satpoint"`
	Timestamp    int64    `json:"timestamp"`
	Inscriptions []string `json:"inscriptions"`
}

type Block struct {
	Hash         string   `json:"hash"`
	Target       string   `json:"target"`
	BestHeight   uint32   `json:"best_height"`
	Height       uint64   `json:"height"`
	Inscriptions []string `json:"inscriptions"`
}

type Status struct {
	BlessedInscriptions uint64 `json:"blessed_inscriptions"`
	Chain               string `json:"chain"`
	ContentTypeCounts   []any  `json:"content_type_counts"`
	CursedInscriptions  uint64 `json:"cursed_inscriptions"`
	Height              uint64 `json:"height"`
	InitialSyncTime     struct {
		Secs  uint64 `json:"secs"`
		Nanos uint64 `json:"nanos"`
	} `json:"initial_sync_time"`
	Inscriptions            uint64 `json:"inscriptions"`
	LostSats                uint64 `json:"lost_sats"`
	MinimumRuneForNextBlock string `json:"minimum_rune_for_next_block"`
	RuneIndex               bool   `json:"rune_index"`
	Runes                   uint64 `json:"runes"`
	SatIndex                bool   `json:"sat_index"`
	Started                 string `json:"started"`
	TransactionIndex        bool   `json:"transaction_index"`
	UnrecoverablyReorged    bool   `json:"unrecoverably_reorged"`
	Uptime                  struct {
		Secs  uint64 `json:"secs"`
		Nanos uint64 `json:"nanos"`
	} `json:"uptime"`
}

type RBlockInfo struct {
	AverageFee       uint64  `json:"average_fee"`
	AverageFeeRate   uint64  `json:"average_fee_rate"`
	Bits             uint32  `json:"bits"`
	Chainwork        string  `json:"chainwork"`
	Confirmations    int32   `json:"confirmations"`
	Difficulty       float64 `json:"difficulty"`
	Hash             string  `json:"hash"`
	Height           uint32  `json:"height"`
	MaxFee           uint64  `json:"max_fee"`
	MaxFeeRate       uint64  `json:"max_fee_rate"`
	MaxTxSize        uint64  `json:"max_tx_size"`
	MedianFee        uint64  `json:"median_fee"`
	MedianTime       uint64  `json:"median_time"`
	MerkleRoot       string  `json:"merkle_root"`
	MinFee           uint64  `json:"min_fee"`
	MinFeeRate       uint64  `json:"min_fee_rate"`
	NextBlock        string  `json:"next_block"`
	Nonce            uint32  `json:"nonce"`
	PreviousBlock    string  `json:"previous_block"`
	Subsidy          uint64  `json:"subsidy"`
	Target           string  `json:"target"`
	Timestamp        uint64  `json:"timestamp"`
	TotalFee         uint64  `json:"total_fee"`
	TotalSize        uint64  `json:"total_size"`
	TotalWeight      uint64  `json:"total_weight"`
	TransactionCount uint64  `json:"transaction_count"`
	Version          uint32  `json:"version"`
}

type RInscriptionIdsPagination struct {
	Ids  []string `json:"ids" example:"79b0e9dbfaf11e664abafbd8fec7d734bfa2d59013f25c50aaac1264f700832di0"`
	More bool     `json:"more" example:"false"`
	Page int      `json:"page"`
}

type InscriptionId struct {
	Id string `json:"id" example:"79b0e9dbfaf11e664abafbd8fec7d734bfa2d59013f25c50aaac1264f700832di0"`
}

// Ordx customized for /ordx/block/inscriptions/:height

type OrdxBlockInscription struct {
	GenesesAddress string      `json:"genesesaddress"`
	Inscription    Inscription `json:"inscription"`
}

type OrdxBlockInscriptions struct {
	Height       uint64                 `json:"height"`
	Inscriptions []OrdxBlockInscription `json:"inscriptions"`
}

type OrdxBlockTxOutputInscriptions struct {
	Height       uint64        `json:"height"`
	Inscriptions []Inscription `json:"inscriptions"`
}
