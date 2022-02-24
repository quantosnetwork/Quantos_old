package address

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"log"
	"sync"
)

type Purpose = uint32

type Config chaincfg.Params

type Address string

func init() {
	config := chaincfg.Params{}
	config.Name = "Quantos"
	config.Checkpoints = make([]chaincfg.Checkpoint, 0)
	qbitBytes := []byte("QBIT")[:]
	config.HDPrivateKeyID = [4]byte{qbitBytes[0],qbitBytes[1],qbitBytes[2],qbitBytes[3]}
	config.BIP0034Height = int32(0)
	config.BIP0065Height = int32(0)
	config.BIP0066Height = int32(0)
	config.CoinbaseMaturity = 0
	config.DefaultPort = "55335"
	config.DNSSeeds = make([]chaincfg.DNSSeed,0)
	config.Bech32HRPSegwit = "QBit-"
	config.HDCoinType = CoinTypeQTO

	chaincfg.MainNetParams = config

	/*
	chaincfg.MainNetParams.HDCoinType = CoinTypeQTO
	chaincfg.Params{Name: ""}
	chaincfg.MainNetParams.HDPrivateKeyID = [4]byte{0x59, 0x58, 0x55, 0x56}
	chaincfg.MainNetParams.Bech32HRPSegwit = "QT"
	spew.Dump(chaincfg.MainNetParams)*/
}

const (
	PurposeBIP44 Purpose = 0x8000002C // 44' BIP44
	PurposeBIP49 Purpose = 0x80000031 // 49' BIP49
	PurposeBIP84 Purpose = 0x80000054 // 84' BIP84
)

type CoinType = uint32

const (
	// QTO = QuantusOS Coin
	CoinTypeQTO = 0x8000038a //906
)


const (
	Apostrophe uint32 = 0x80000000 // 0'
)

type Key struct {
	path     string
	bip32Key *bip32.Key
}

func (k *Key) Encode(compress bool) (wif, address, segwitBech32, segwitNested string, err error) {
	prvKey, _ := btcec.PrivKeyFromBytes(btcec.S256(),k.bip32Key.Key)
	return GenerateFromBytes(prvKey, compress)
}

func (k *Key) GetPath() string {
	return k.path
}

type KeyManager struct {
	mnemonic   string
	passphrase string
	keys       map[string]*bip32.Key
	mux        sync.Mutex
}

// NewKeyManager return new key manager
// bitSize has to be a multiple 32 and be within the inclusive range of {128, 256}
// 128: 12 phrases
// 256: 24 phrases
func NewKeyManager(bitSize int, passphrase, mnemonic string) (*KeyManager, error) {

	if mnemonic == "" {
		entropy, err := bip39.NewEntropy(bitSize)
		if err != nil {
			return nil, err
		}
		mnemonic, err = bip39.NewMnemonic(entropy)
		if err != nil {
			return nil, err
		}
	}

	km := &KeyManager{
		mnemonic:   mnemonic,
		passphrase: passphrase,
		keys:       make(map[string]*bip32.Key, 0),
	}
	return km, nil
}

func (km *KeyManager) GetMnemonic() string {
	return km.mnemonic
}

func (km *KeyManager) GetPassphrase() string {
	return km.passphrase
}

func (km *KeyManager) GetSeed() []byte {
	return bip39.NewSeed(km.GetMnemonic(), km.GetPassphrase())
}

func (km *KeyManager) getKey(path string) (*bip32.Key, bool) {
	km.mux.Lock()
	defer km.mux.Unlock()

	key, ok := km.keys[path]
	return key, ok
}

func (km *KeyManager) setKey(path string, key *bip32.Key) {
	km.mux.Lock()
	defer km.mux.Unlock()

	km.keys[path] = key
}

func (km *KeyManager) GetMasterKey() (*bip32.Key, error) {
	path := "m"

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	key, err := bip32.NewMasterKey(km.GetSeed())
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) GetPurposeKey(purpose uint32) (*bip32.Key, error) {
	path := fmt.Sprintf(`m/%d'`, purpose-Apostrophe)

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	parent, err := km.GetMasterKey()
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(purpose)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) GetCoinTypeKey(purpose, coinType uint32) (*bip32.Key, error) {
	path := fmt.Sprintf(`m/%d'/%d'`, purpose-Apostrophe, coinType-Apostrophe)

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	parent, err := km.GetPurposeKey(purpose)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(coinType)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) GetAccountKey(purpose, coinType, account uint32) (*bip32.Key, error) {
	path := fmt.Sprintf(`m/%d'/%d'/%d'`, purpose-Apostrophe, coinType-Apostrophe, account)

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	parent, err := km.GetCoinTypeKey(purpose, coinType)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(account + Apostrophe)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

// GetChangeKey ...
// https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki#change
// change constant 0 is used for external chain
// change constant 1 is used for internal chain (also known as change addresses)
func (km *KeyManager) GetChangeKey(purpose, coinType, account, change uint32) (*bip32.Key, error) {
	path := fmt.Sprintf(`m/%d'/%d'/%d'/%d`, purpose-Apostrophe, coinType-Apostrophe, account, change)

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	parent, err := km.GetAccountKey(purpose, coinType, account)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(change)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) GetKey(purpose, coinType, account, change, index uint32) (*Key, error) {
	path := fmt.Sprintf(`m/%d'/%d'/%d'/%d/%d`, purpose-Apostrophe, coinType-Apostrophe, account, change, index)

	key, ok := km.getKey(path)
	if ok {
		return &Key{path: path, bip32Key: key}, nil
	}

	parent, err := km.GetChangeKey(purpose, coinType, account, change)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(index)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return &Key{path: path, bip32Key: key}, nil
}

func Generate(compress bool) (wif, address, segwitBech32, segwitNested string, err error) {
	prvKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return "", "", "", "", err
	}
	return GenerateFromBytes(prvKey, compress)
}

func GenerateFromBytes(prvKey *btcec.PrivateKey, compress bool) (wif, address, segwitBech32, segwitNested string, err error) {
	// generate the wif(wallet import format) string

	btcwif, err := btcutil.NewWIF(prvKey, &chaincfg.MainNetParams, compress)
	if err != nil {
		return "", "", "", "", err
	}
	wif = btcwif.String()

	// generate a normal p2pkh address
	serializedPubKey := btcwif.SerializePubKey()
	addressPubKey, err := btcutil.NewAddressPubKey(serializedPubKey, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", "", "", err
	}
	address = addressPubKey.EncodeAddress()

	// generate a normal p2wkh address from the pubkey hash
	witnessProg := btcutil.Hash160(serializedPubKey)
	addressWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", "", "", err
	}
	segwitBech32 = addressWitnessPubKeyHash.EncodeAddress()

	// generate an address which is
	// backwards compatible to Bitcoin nodes running 0.6.0 onwards, but
	// allows us to take advantage of segwit's scripting improvments,
	// and malleability fixes.
	serializedScript, err := txscript.PayToAddrScript(addressWitnessPubKeyHash)
	if err != nil {
		return "", "", "", "", err
	}
	addressScriptHash, err := btcutil.NewAddressScriptHash(serializedScript, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", "", "", err
	}
	segwitNested = addressScriptHash.EncodeAddress()

	return wif, address, segwitBech32, segwitNested, nil
}

func NewAddress(compress bool, pass string, mnemonic string) {

	number := 1



		km, err := NewKeyManager(256, pass, mnemonic)
		if err != nil {
			log.Fatal(err)
		}
		masterKey, err := km.GetMasterKey()
		if err != nil {
			log.Fatal(err)
		}
		passphrase := km.GetPassphrase()
		if passphrase == "" {
			passphrase = "<none>"
		}
		fmt.Printf("\n%-18s %s\n", "Mnemonic (backup in a safe place!):", km.GetMnemonic())
		fmt.Printf("Passphrase: (if you setted one you know it) \n")
		//fmt.Printf("%-18s %x\n", "BIP39 Seed:", km.GetSeed())
		fmt.Printf("Your Master Key (do not share!): %s\n", masterKey.B58Serialize())


		for i := 0; i < number; i++ {
			key, err := km.GetKey(PurposeBIP84, CoinTypeQTO, 0, 0, uint32(i))
			if err != nil {
				log.Fatal(err)
			}
			_, _, segwitBech32, _, err := key.Encode(compress)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Your Quantos Wallet Address: %s\n", segwitBech32)
		}
		fmt.Println()

}