package decoder

import (
	"Quantos/crypto"
	"bytes"
	"errors"
	"strconv"
)

type Decoder struct {
	data []byte
	length int
	cursor int
}

func (d *Decoder) Decode(data []byte) (interface{}, error) {
	d.data = data
	d.length = len(data)
	return d.decode()
}

func (d *Decoder) decode() (interface{}, error) {
	switch d.data[d.cursor] {
	case 'i':
		return d.decodeInt()
	case 'l':
		d.cursor += 1
		list := []interface{}{}
		for {
			if d.cursor == d.length {
				return nil, errors.New("bencode: invalid list field")
			}
			if d.data[d.cursor] == 'e' {
				d.cursor += 1
				return list, nil
			}
			value, err := d.decode()
			if err != nil {
				return nil, err
			}
			list = append(list, value)
		}
	case 'd':
		d.cursor += 1
		dictionary := map[string]interface{}{}
		for {
			if d.cursor == d.length {
				return nil, errors.New("bencode: invalid dictionary field")
			}
			if d.data[d.cursor] == 'e' {
				d.cursor += 1
				return dictionary, nil
			}
			key, err := d.decodeBytes()
			if err != nil {
				return nil, errors.New("bencode: non-string dictionary key")
			}
			value, err := d.decode()
			if err != nil {
				return nil, err
			}
			dictionary[crypto.BytesToString(key)] = value
		}
	case 'h':
		d.cursor += 1
		hashtable := map[int]interface{}{}
		for {
			if d.cursor == d.length {
				return nil, errors.New("quantos decoding: invalid hashtable")
			}
			if d.data[d.cursor] == 'e' {
				d.cursor += 1
				return hashtable, nil
			}
			key, err := d.decodeBytes()
			if err != nil {
				return nil, errors.New("bencode: non-string dictionary key")
			}
			value, err := d.decode()
			if err != nil {
				return nil, err
			}
			hashtable[crypto.BytesToInt(key)] = value
		}

	default:
		return d.decodeBytes()
	}
}

func (d *Decoder) decodeBytes() ([]byte, error) {
	if d.data[d.cursor] < '0' || d.data[d.cursor] > '9' {
		return nil, errors.New("quantos bytes decoder: invalid string field")
	}
	index := bytes.IndexByte(d.data[d.cursor:], ':')
	if index == -1 {
		return nil, errors.New("quantos bytes decoder:  invalid string field")
	}
	index += d.cursor
	stringLength, err := d.parseInt(d.data[d.cursor:index])
	if err != nil {
		return nil, err
	}
	index += 1
	endIndex := index + int(stringLength)
	if endIndex > d.length {
		return nil, errors.New("quantos bytes decoder: not a valid string")
	}
	value := d.data[index:endIndex]
	d.cursor = endIndex
	return value, nil
}

var (
	pow10i64 = [...]int64{
		1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
		1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18,
	}
	pow10i64Len = len(pow10i64)
)

func (d *Decoder) parseInt(data []byte) (int64, error) {
	isNegative := false
	if data[0] == '-' {
		data = data[1:]
		isNegative = true
	}
	maxDigit := len(data)
	if maxDigit > pow10i64Len {
		return 0, errors.New("quantos int parser: invalid length of number")
	}
	sum := int64(0)
	for i, b := range data {
		if b < '0' || b > '9' {
			return 0, errors.New("quantos int parser: invalid integer byte: " + strconv.FormatUint(uint64(b), 10))
		}
		c := int64(b) - 48
		digitValue := pow10i64[maxDigit-i-1]
		sum += c * digitValue
	}
	if isNegative {
		return -1 * sum, nil
	}
	return sum, nil
}

func (d *Decoder) decodeInt() (interface{}, error) {
	d.cursor += 1
	index := bytes.IndexByte(d.data[d.cursor:], 'e')
	if index == -1 {
		return nil, errors.New("quantos int decoder: invalid integer field")
	}
	index += d.cursor
	integer, err := d.parseInt(d.data[d.cursor:index])
	if err != nil {
		return nil, err
	}
	d.cursor = index + 1
	return integer, nil
}