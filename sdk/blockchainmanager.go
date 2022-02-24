package sdk

type BlockchainManager interface {
	CreateNewBlockchain()
	SetConfigurationFile()
	GetConfiguration(config string) (value interface{}, err error)
	GetPublicInfo()
	GetGenesisBlock()
	ValidateGenesisBlock()
	GetValidators()
	ValidateValidators()
	GetBlockById(blockID string)
	ValidateBlock(blockID string)
	CloseBlock(blockID string)
	CreateBlock(data ...interface{})
	GetLastBlock()
	GetBlockMerkleRoot(blockID string)
	GetTXMerkleRoot(txID string)
	GetReceiptMerkleRoot(rID string)
	CreateTx(txData ...interface{})
	SendTx(txData ...interface{})
	SignTx(signData ...interface{})
	GetLastTimeStamp()
	GetBlockByTxID(txId string)
	GetTxByID(txId string)
	Consensus() interface{}
	GetOrphanBlocks()
	GetPendingTxs()
	GetPendingBlocks()
	GetTxQueue()
	GetBlockQueue()
	Version()
	CoinbaseAddress()
	Coin()
	Tokens()
	Contracts()

}

type Coins interface{
	MintTo(addr string, amt uint)
	BurnTo(addr string, amt uint)
	TotalAvailable()
	Unspent()
	Pairings()
	PairingValue()
	Description()
	URI()
	Coinbase()
	Name()
	JSONInfo()
	Contract()
}
type Token interface{
	Address()
	ContractAddress()
	TokenType()
	Coins
}
type Contract interface{
	Code()
	Hex()
	CompiledBinary()
	WASM()
	VM() VM
}

