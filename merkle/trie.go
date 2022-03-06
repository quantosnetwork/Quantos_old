package merkle

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/zeebo/blake3"
	"go.uber.org/atomic"
	"log"
	"strings"
	"sync"
	"time"
)

/*

	@dev Info about the merkle trie structure

	Blake3 Merkle tree (sum512)
	Leaf size: 128 [64]bytes
	Root size: 128 [64]bytes
	MergeLeaf size: 256 [128]bytes

	ML = MergeLeaf (Node hash)
	C = count
	L = Leaf
	mr merkle root
	H = Blake3 hasher sum 512
	++ = all of that type

	f(mr) = H SUM(ML++)
	ML = (LH[1] + LH[2])
	LH = H(L)

*/

type Hash [64]byte

type MerkleTrie interface {
	Get(key []byte) ([]byte, bool)
	Put(key []byte, value []byte)
	Del(key []byte, value []byte) bool
}

type treeBuilder struct {
	treeContents []TreeContent
	size         int
	height       int
	leavesTotal  int
	state        TreeStates
}

type TreeContent [][]byte

type Tree struct {
	MerkleRoot []byte
	TrunkRoots [][]byte
	TreeHash   []byte
	Branches   map[int]*branch
	Hashable
}

type branch struct {
	// a branch can have multiple branches attached to it before there are leaves
	hasChildren      bool
	hasParent        bool
	multipleBranches []*branch
	childCount       atomic.Bool
	leftLeaf         *TreeLeaf
	rightLeaf        *TreeLeaf
	// in real life that would be the length of the branch
	byteLen int
	Hashable
}

type TreeLeaf struct {
	parentBranch *branch

	content TreeContent
	context string
	Hashable
}

func (tc TreeContent) hash() Hash {
	b := tc.Bytes()
	return hashFn(b)
}

func (tc TreeContent) Bytes() []byte {
	b, _ := json.Marshal(tc)
	return b
}

type Hashable interface {
	hash() Hash
}

type EmptyLeaf struct {
}

func (e EmptyLeaf) hash() Hash {
	return hashFn(nil)
}

func (b *branch) hash() Hash {
	var l, r [64]byte
	l = b.leftLeaf.hash()
	r = b.rightLeaf.hash()
	return hashFn(append(l[:], r[:]...))
}

func (l *TreeLeaf) hash() Hash {
	return hashFn(l.content.Bytes())
}

func hashFn(data []byte) Hash {
	return blake3.Sum512(data)
}

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

func NewTree(contents []TreeContent) {

	tree := &Tree{}
	c := make([]TreeContent, len(contents))
	copy(c, contents)

	tree.Branches = map[int]*branch{}
	b := tree.newBuilder(c)
	b._prepare(tree)
	if b.state == DONE {
		fmt.Printf("merkle trie done: %v", tree)
	}
	hashes := tree.Walk()
	mr := tree.hashMerkleRoot(hashes)
	log.Printf("tree merkle root: %v", hex.EncodeToString(mr))

}

func (m *Tree) newBuilder(contents []TreeContent) *treeBuilder {
	builder := &treeBuilder{}
	builder.height = 0
	builder.state = 0
	builder.size = 0
	builder.leavesTotal = 0
	builder.treeContents = contents
	return builder
}

func (m *Tree) Build(contents []Hashable) {

}

func (m *Tree) Put(key []byte, value []byte) {}
func (m *Tree) Get(key []byte) ([]byte, bool) {
	return nil, false
}
func (m *Tree) Del(key []byte, value []byte) bool {
	return false
}

type TreeStates int

const (
	NONE TreeStates = iota
	WORKING
	DONE
	ERRORED
	IDLE
	VERIFIED
	ARCHIVED
	ARCHIVING
	SNAPSHOT
	WALKING
	HASHING
)

func (tb *treeBuilder) _prepare(m *Tree) {

	tb.state = IDLE
	//  contentLen := len(tb.treeContents)
	for i := range tb.treeContents {
		tb.state = WORKING
		tb._addSingleBranch(m)
		// now we get the height to add the leaves
		height := tb.height - 1
		if height < len(tb.treeContents)-1 {
			tb._addLeaves(m, height, tb.treeContents[i], tb.treeContents[i+1])
		} else {
			tb._addLeaves(m, height, tb.treeContents[i], tb.treeContents[i-1])
		}
	}
	tb.state = DONE
}

func (tb *treeBuilder) _addSingleBranch(m *Tree) {
	branch := &branch{}
	branch.hasChildren = false
	branch.hasParent = false
	branch.leftLeaf = new(TreeLeaf)
	branch.rightLeaf = new(TreeLeaf)

	m.Branches[tb.height] = branch
	tb.height++
}

func (tb *treeBuilder) _addLeaves(m *Tree, height int, left, right Hashable) {
	m.Branches[height].leftLeaf = new(TreeLeaf)
	m.Branches[height].leftLeaf.content = left.(TreeContent)
	m.Branches[height].rightLeaf = new(TreeLeaf)
	m.Branches[height].rightLeaf.content = right.(TreeContent)

}

func (t *Tree) Print() {

	merkleRoot := t.hash()
	fmt.Println(merkleRoot)

}

func (t *Tree) hash() Hash {
	return Hash{}
}

type treeWalkFunction = func(b *branch, h func([]byte))

var TreeReceipt []byte

func (t *Tree) Walk() []map[int][]string {
	branches := t.Branches
	state := IDLE
	//height := 0
	var hashes []map[int][]string
	hashes = make([]map[int][]string, len(branches))
	for i := 0; i < len(branches); i++ {
		state = WALKING
		//height = i
		branch := branches[i]

		state = HASHING
		left := branch.leftLeaf.hash()
		right := branch.rightLeaf.hash()

		leaves := map[int][]string{}
		leaves[0] = []string{left.String(), right.String()}
		hashes = append(hashes, leaves)
		/*hashes[height][0] = left.String()
		hashes[height][1] = right.String()*/

	}

	state = DONE

	if state == DONE {
		return hashes
	}
	return t.Walk()
	//return nil
}

func (t *Tree) hashMerkleRoot(hashes []map[int][]string) []byte {

	var mu sync.Mutex

	root := &struct {
		height    int
		timestamp int64
		hash      []byte
		String    string
	}{}

	hasher := blake3.New()
	//	hasher.Reset()

	root.height = 0
	root.timestamp = time.Now().UnixNano()
	hs := make([][]byte, len(hashes))
	mu.Lock()
	defer mu.Unlock()
	for i := 0; i < len(hashes)-1; i++ {

		toJoin := hashes[i][0]
		hsj := strings.Join(toJoin, "")

		if hsj != "" {
			spew.Dump(hsj)
			hs[i] = []byte(hsj)
			hasher.Write(hs[i])
			root.height += i
		}
	}
	rhash := blake3.Sum512(hasher.Sum(nil))
	root.hash = rhash[:]
	root.String = hex.EncodeToString(root.hash)

	t.MerkleRoot = root.hash
	return t.MerkleRoot
}
