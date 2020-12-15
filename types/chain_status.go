package types

type ChainStatus struct {
	SpecVersion        uint32 `json:"spec_version"`
	TransactionVersion uint32 `json:"tx_version"`
	BlockHash          string `json:"block_hash"`
	GenesisHash        string `json:"genesis_hash"`
}
