package sdk

import (
	"Quantos/address"
	"Quantos/sdk/config"
)

type AddressSDK interface {
	InitSDK(netID string)
	GenerateQBITWalletAddress(out string) string
	VerifyAddress(in string, out bool)
	IsZeroAddress(in string, out bool)
	GenerateTXAddress(in InputData, out OutputData)
	GenerateBlockAddress(in InputData, out OutputData)
	GetZeroAddress() string
}

type InputData struct {
	data interface{}
}

type OutputData struct {
	data string
}

type addrFunctions struct{}

func (a addrFunctions) IsZeroAddress(in string, out bool) {
	//TODO implement me
	panic("implement me")
}

func (a addrFunctions) GenerateTXAddress(in InputData, out OutputData) {
	//TODO implement me
	panic("implement me")
}

func (a addrFunctions) GenerateBlockAddress(in InputData, out OutputData) {
	//TODO implement me
	panic("implement me")
}

var CurrentNetworkID config.NetworkID

func (a addrFunctions) GenerateQBITWalletAddress(out string) string {

	addr := address.GenerateNewQbitAddress(CurrentNetworkID, config.Version, config.QBIT_ADDRESS_PREFIX, uint32(0))
	out = addr.String()
	return out
}

func (a addrFunctions) InitSDK(netID string) {
	switch netID {
	case "live":

		CurrentNetworkID = config.LIVENET
		return
	case "test":
		CurrentNetworkID = config.TESTNET
		return
	case "local":
		CurrentNetworkID = config.LOCALNET
		return
	}
	return
}

func initZeroAddress() {

}

func (a addrFunctions) VerifyAddress(in string, out bool) {

}

func (a addrFunctions) GetZeroAddress() string {
	return config.ZEROADDRESS
}

func GetAddressSDK() AddressSDK {
	var a addrFunctions
	return a
}
