package protocol

import (
	"Quantos/crypto"
	"bufio"
	"errors"
	"fmt"
	"github.com/open-quantum-safe/liboqs-go/oqs"
	"io"
	"log"
	"net"
)

const KEMNAME = "kyber512"


type keyExchange interface {
	initKemKX(host, port string)
	handleKemKX(conn net.Conn)
}

type KeyExchange struct {
	conn net.Conn
	clientKeys *crypto.KemKeys
	serverKeys *crypto.ServerKeys
}

func (kex KeyExchange) initKemKX(host, port string) {
	kex.clientKeys = crypto.NewKemClient()
	kex.serverKeys = crypto.NewKemServer()

	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		kex.conn, err = listener.Accept()

		if err != nil {
			panic(err)
		}
		go kex.handleKemKX(kex.conn)
	}

}

func (kex KeyExchange) handleKemKX(conn net.Conn) {
	defer conn.Close()
	_, err := fmt.Println(conn, KEMNAME)
	if err != nil {
		conn.Close()
	}
	clientPublicKey := make([]byte, kex.serverKeys.Server.Details().LengthPublicKey)
	_, err = io.ReadFull(conn, clientPublicKey)
	if err != nil {
		conn.Write([]byte("could not read key\n"))
		conn.Close()
		return
	}
	cipherText, sharedSecret, _  := kex.serverKeys.Server.EncapSecret(clientPublicKey)
	_, err = conn.Write(cipherText)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("\nConnection #%d - server shared secret:\n% X ... % X\n\n", sharedSecret[0:8],
		sharedSecret[len(sharedSecret)-8:])

	conn.Write([]byte("AUTHENTICATED"))
}

func StartKeyExchange(host, port string) {

	var kex KeyExchange
	kex.initKemKX(host, port)

}

func StartKexClient() {
	conn, err := net.Dial("tcp", "127.0.0.1:55225")
	if err != nil {
		log.Fatal(errors.New("client cannot connect " +
			"to host on port 55225"))
	}
	defer conn.Close()
	client := oqs.KeyEncapsulation{}
	defer client.Clean()
	kemName, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		panic(err)
		log.Fatal(errors.New("client cannot receive the " +
			"KEM name from the server"))
	}
	kemName = kemName[:len(kemName)-1]
	if err := client.Init(kemName, nil); err != nil {
		log.Fatal(err)
	}
	clientPublicKey, err := client.GenerateKeyPair()
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write(clientPublicKey)
	if err != nil {
		log.Fatal(errors.New("client cannot send the public key to the " +
			"server"))
	}

	// listen for reply from the server, e.g. for the encapsulated secret
	ciphertext := make([]byte, client.Details().LengthCiphertext)
	n, err := io.ReadFull(conn, ciphertext)
	if err != nil {
		log.Fatal(err)
	} else if n != client.Details().LengthCiphertext {
		log.Fatal(errors.New("client expected to read " +
			string(client.Details().LengthCiphertext) + " bytes, but instead " +
			"read " + string(n)))
	}

	// decapsulate the secret and extract the shared secret
	sharedSecretClient, err := client.DecapSecret(ciphertext)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(client.Details())
	fmt.Printf("\nClient shared secret:\n% X ... % X\n",
		sharedSecretClient[0:8], sharedSecretClient[len(sharedSecretClient)-8:])

}