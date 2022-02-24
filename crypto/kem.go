package crypto


import  (
oqs "github.com/open-quantum-safe/liboqs-go/oqs"
	"log"
)

type KemKeys struct {
	Client oqs.KeyEncapsulation
	PubKey []byte
	SharedSecret chan []byte
}

type ServerKeys struct {
	Server oqs.KeyEncapsulation
	ciphertext []byte
	SharedSecret []byte
	ClientPubKey chan []byte
}



func NewKemClient() *KemKeys {
	k := &KemKeys{}
	client := oqs.KeyEncapsulation{}
	k.Client = client
	defer k.Client.Clean()
	if err := k.Client.Init("kyber512", nil); err != nil {
		panic(err)
	}

	pk, err := k.Client.GenerateKeyPair()
	if err != nil {
		panic(err)
	}
	k.PubKey = pk
	k.SharedSecret = make(chan []byte)

	return k

}

func NewKemServer() *ServerKeys {
	k := &ServerKeys{}
	k.Server = oqs.KeyEncapsulation{}
	defer k.Server.Clean()
	if err := k.Server.Init("kyber512", nil); err != nil {
		log.Fatal(err)
	}
	k.ClientPubKey = make(chan []byte)
	return k
}

