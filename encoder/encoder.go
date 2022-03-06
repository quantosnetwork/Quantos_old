package encoder

import (
	"github.com/quantosnetwork/Quantos/crypto"

	"fmt"

	"sort"
	"sync"
	"unsafe"
)

//go:linkname memmov runtime.memmove
func memmov(to unsafe.Pointer, from unsafe.Pointer, n uintptr)

type Encoder struct {
	buffer []byte
	length int
	offset int
}

type sliceHeader struct {
	data unsafe.Pointer
	len  int
	cap  int
}

//go:nosplit
func (e *Encoder) grow(neededLength int) {
	availableLength := e.length - e.offset
	if availableLength >= neededLength {
		return
	}
	if e.length == 0 {
		if neededLength < 16 {
			neededLength = 16
		}
		e.length = neededLength
		availableLength = neededLength
	} else {
		for availableLength < neededLength {
			e.length += e.length
			availableLength = e.length - e.offset
		}
	}
	buffer := make([]byte, e.length)
	memmov(
		unsafe.Pointer(uintptr((*sliceHeader)(unsafe.Pointer(&buffer)).data)),
		(*sliceHeader)(unsafe.Pointer(&e.buffer)).data,
		uintptr(e.offset),
	)
	e.buffer = buffer

}

//go:nosplit
func (e *Encoder) write(data []byte) {
	length := len(data)
	memmov(
		unsafe.Pointer(uintptr((*sliceHeader)(unsafe.Pointer(&e.buffer)).data)+uintptr(e.offset)),
		(*sliceHeader)(unsafe.Pointer(&data)).data,
		uintptr(length),
	)
	e.offset += length
}

//go:nosplit
func (e *Encoder) writeByte(data byte) {
	*(*byte)(unsafe.Pointer(uintptr((*sliceHeader)(unsafe.Pointer(&e.buffer)).data) + uintptr(e.offset))) = data
	e.offset++
}

func (e *Encoder) EncodeTo(dst []byte, data interface{}) ([]byte, error) {
	if cap(dst) > len(dst) {
		dst = dst[:cap(dst)]
	} else if len(dst) == 0 {
		dst = make([]byte, 512)
	}
	e.buffer = dst
	e.length = cap(dst)
	err := e.encode(data)
	if err != nil {
		return nil, err
	}
	return e.buffer[:e.offset], nil
}

func (e *Encoder) encode(data interface{}) error {
	switch value := data.(type) {
	case int64:
		e.encodeInt(value)
	case int32:
		e.encodeInt(int64(value))
	case int16:
		e.encodeInt(int64(value))
	case int8:
		e.encodeInt(int64(value))
	case int:
		e.encodeInt(int64(value))
	case uint64:
		e.encodeInt(int64(value))
	case uint32:
		e.encodeInt(int64(value))
	case uint16:
		e.encodeInt(int64(value))
	case uint8:
		e.encodeInt(int64(value))
	case uint:
		e.encodeInt(int64(value))
	case []byte:
		e.encodeBytes(value)
	case string:
		e.encodeBytes(crypto.StringToBytes(value))
	case []interface{}:
		return e.encodeList(value)
	case map[string]interface{}:
		return e.encodeDictionary(value)
	case map[int]interface{}:
		return e.encodeHashTable(value)
	default:
		return fmt.Errorf("quantos encoding: unsupported type: %T", value)
	}
	return nil
}

//go:nosplit
func (e *Encoder) encodeBytes(data []byte) {
	dataLength := len(data)
	e.grow(dataLength + 23)
	e.writeInt(int64(len(data)))
	e.writeByte(':')
	e.write(data)
}

func (e *Encoder) encodeList(data []interface{}) error {
	e.grow(1)
	e.writeByte('l')
	for _, data := range data {
		err := e.encode(data)
		if err != nil {
			return err
		}
	}
	e.grow(1)
	e.writeByte('e')
	return nil
}

const stringsArrayLen = 20

var stringsArrayPool = sync.Pool{
	New: func() interface{} {
		return &[stringsArrayLen]string{}
	},
}

func sortStrings(ss []string) {
	if len(ss) <= stringsArrayLen {
		for i := 1; i < len(ss); i++ {
			for j := i; j > 0; j-- {
				if ss[j] >= ss[j-1] {
					break
				}
				ss[j], ss[j-1] = ss[j-1], ss[j]
			}
		}
	} else {
		sort.Strings(ss)
	}
}

func (e *Encoder) encodeDictionary(data map[string]interface{}) error {
	e.grow(1)
	e.writeByte('d')
	var keys []string
	if len(data) <= stringsArrayLen {
		stringsArray := stringsArrayPool.Get().(*[stringsArrayLen]string)
		defer stringsArrayPool.Put(stringsArray)
		keys = stringsArray[:0:len(data)]
	} else {
		keys = make([]string, 0, len(data))
	}
	for key, _ := range data {
		keys = append(keys, key)
	}
	sortStrings(keys)
	for _, key := range keys {
		e.encodeBytes(crypto.StringToBytes(key))
		err := e.encode(data[key])
		if err != nil {
			return err
		}
	}
	e.grow(1)
	e.writeByte('e')
	return nil
}

func (e *Encoder) encodeHashTable(data map[int]interface{}) error {
	e.grow(1)
	e.writeByte('h')
	var keys []int
	keys = make([]int, 0, len(data))
	for key, _ := range data {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	for _, kk := range keys {
		e.encodeInt(int64(kk))
	}

	e.grow(1)
	e.writeByte('e')
	return nil
}
