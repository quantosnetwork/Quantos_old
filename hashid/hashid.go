package hashid

import (
	"Quantos/crypto"
	"bytes"
	"hash"
)

/*

	HashID =

	ChainCode + Genesis Hash + Current block hash + Chain hash up to that block

	@description
	## Purpose
	Faster chain sync and more flexibility than when dealing with full structure.

 */

type hashID interface {
	_buildHashId(
		chainCode []byte,
		genesisHash []byte,
		currBlockHash []byte,
		chainHash []byte) (hashID, error)
	_verifyHashID()
	_signHashID()
	_getHashIDFromHashTable(hID string) []byte
	_encodeHashId() []byte
	_decodeHashId() hashID
	_insertIntoHashTable() (int, error)
	Set()
	Get()
}

type HashID struct {
	raw []byte
	workBuffer *bytes.Buffer
	hashFunc hash.Hash
	ID hashID `json:"hash_id"`
	HtID int
}

type HashIDSignature struct {
	ID []byte
	crypto.Signature
}

func (hid *HashID) _buildHashId(
	chainCode []byte,
	genesisHash []byte,
	currBlockHash []byte,
	chainHash []byte) (hashID, error) {
	return nil, nil
}

func (hid *HashID) _verifyHashID() {}

func (hid *HashID) _signHashID() {}

func (hid *HashID) _getHashIDFromHashTable(hID string) []byte {
	return nil
}

func (hid *HashID) _encodeHashId() []byte{
	return nil
}

func (hid *HashID) _decodeHashId() hashID{
	return nil
}

func (hid *HashID) _insertIntoHashTable() (int, error) {
	return 0, nil
}

func (hid *HashID) Set() {}

func (hid *HashID) Get() {}
