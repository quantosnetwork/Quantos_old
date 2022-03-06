package sdk

import (
	"github.com/quantosnetwork/Quantossdk/config"
	"github.com/quantosnetwork/Quantosuint512"
	"github.com/holiman/uint256"
	"math/big"
)

type AddressSDK interface {
	InitSDK(netID string)
	GenerateMasterWalletAddress() (*uint512.Address, string)
	VerifyAddress(in string, out bool)
	IsZeroAddress(in string, out bool)
	GenerateTXAddress(in InputData, out OutputData)
	GenerateBlockAddress(in InputData, out OutputData)
	GetZeroAddress() string
	DeriveFromMaster(master *uint512.Address, derivationPath string) string
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

func (a addrFunctions) GenerateMasterWalletAddress() (*uint512.Address, string) {

	addr := &uint512.Address{}
	am := addr.Raw.Create()
	m := am.Master()
	masterBig := new(big.Int).SetBytes(m)
	out1, _ := uint256.FromBig(masterBig)
	out := out1.String()

	return am, out
}

func (a addrFunctions) DeriveFromMaster(master *uint512.Address, derivationPath string) string {
	return master.Derive()
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
