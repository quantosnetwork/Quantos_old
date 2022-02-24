package hashtable

import (
	"crypto/sha256"
	"github.com/davecgh/go-spew/spew"
	"log"
	"math/rand"
	"testing"
	"time"
)



func TestHashTable_Put(t *testing.T) {
	htable := &HashTable{}
	contentKey, contentVal := generateHashtableContent(10)
	for i, c := range contentKey {
		htable.Put(c, contentVal[i])
	}


	//htable.PrintHashTable()

	/*if htable.Size() < 10 {
		log.Fatalf("hashtable should have 10 records it only has: %v", htable.Size())
	}*/

}

func createMockHashtable(size int) *HashTable {
	htable := &HashTable{}
	contentKey, contentVal := generateHashtableContent(size)
	for i, c := range contentKey {
		htable.Put(c, contentVal[i])
	}
	return htable
}

func TestHashTable_ToBytes(t *testing.T) {

	ht := createMockHashtable(10124)
	b, err := ht.ToBytes()
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(b)

}

func randomBytes() []byte {
	buf := make([]byte, 32)
	rand.Seed(time.Now().UnixNano())
	rand.Read(buf)
	return buf
}

func generateHashtableContent(size int) ([]Key, []Value) {

	k := make([]Key, size)
	v := make([]Value, size)

	for i := 0; i < size; i++ {
		k[i] = Key(randomBytes())
		v[i] = Value(generateRandomHash(randomBytes()))
	}
	return k, v

}

func generateRandomHash(src []byte) []byte {
	hasher := sha256.New()
	hasher.Write(src)
	return hasher.Sum(nil)
}