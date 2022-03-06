package config

import (
	p2pConf "github.com/quantosnetwork/Quantos/p2p/config"
)

type ChainConfig struct {
	NetID     NetworkID
	Version   [2]byte
	P2PConfig p2pConf.Config
	Host      string
	Port      string
	BlockChan chan interface{}
}

type NetworkID [2]byte

var (
	LIVENET  NetworkID = [2]byte{0x38, 0x0a}
	TESTNET  NetworkID = [2]byte{0x0a, 0x00}
	LOCALNET NetworkID = [2]byte{0x00, 0xff}
)

var Version = [2]byte{0x00, 0x01}

type AddressPrefixes uint32

const (
	DEFAULT_ADDRESS_PREFIX uint32 = iota
	QBIT_ADDRESS_PREFIX
	TX_ADDRESS_PREFIX
	BLOCK_ADDRESS_PREFIX
	CONTRACT_ADDRESS_PREFIX
)

const ZEROADDRESS = "0xQ532Sbdoigjcofdhaylrzqxehahocf4G2Rm6V5Iam57Logocjyw4"
