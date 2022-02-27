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
	"time"
)

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
			tb._addLeaves(m, height, tb.treeContents[i-1], tb.treeContents[i-1])
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

func (t *Tree) Walk() map[int]map[int][]string {
	branches := t.Branches
	state := IDLE
	height := 0
	hashes := map[int]map[int][]string{}
	for i := 0; i < len(branches); i++ {
		state = WALKING
		height = i
		branch := branches[i]
		if branch.hasChildren {
			for j := 0; j < len(branch.multipleBranches); j++ {
				state = WALKING
				child := branch.multipleBranches[j]
				left := child.leftLeaf.hash()
				right := branch.leftLeaf.hash()
				hashes[height][j] = []string{left.String(), right.String()}
			}
		}
		state = HASHING
		left := branch.leftLeaf.hash()
		right := branch.rightLeaf.hash()
		hashes[height] = map[int][]string{}
		hashes[height][i] = []string{left.String(), right.String()}
	}

	state = DONE
	spew.Dump(hashes)
	if state == DONE {
		return hashes
	}
	return t.Walk()
	//return nil
}

func (t *Tree) hashMerkleRoot(hashes map[int]map[int][]string) []byte {

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
	for i := 0; i < len(hashes); i++ {
		if len(hashes[i]) > 1 {
			hs2 := make([][]byte, len(hashes[i]))
			for j := 0; j < len(hashes[i]); j++ {
				hslevel2 := hashes[i][j]
				hs2[j] = []byte(hslevel2[0])
				hasher.Write(hs2[j])
				root.height += j
			}
		}
		hslevel1 := hashes
		hsj := strings.Join(hslevel1[i][i], "")
		//spew.Dump(hsj)
		hs[i] = []byte(hsj)
		hasher.Write(hs[i])
		root.height += i
	}
	rhash := blake3.Sum512(hasher.Sum(nil))
	root.hash = rhash[:]
	root.String = hex.EncodeToString(root.hash)
	//spew.Dump(hashes)

	t.MerkleRoot = root.hash
	return t.MerkleRoot
}
