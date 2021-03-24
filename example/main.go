package main

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/boltdb/bolt"
)

const pageHeaderSize = int(unsafe.Offsetof(((*page)(nil)).ptr))

const branchPageElementSize = int(unsafe.Sizeof(branchPageElement{}))
const leafPageElementSize = int(unsafe.Sizeof(leafPageElement{}))

type pgid uint64

type page struct {
	// 页id
	id pgid
	// flags：页类型，可以是分支，叶子节点，元信息，空闲列表
	flags uint16
	// 个数
	count uint16
	//
	overflow uint32
	//
	ptr uintptr
}

// branchPageElement represents a node on a branch page.
type branchPageElement struct {
	pos   uint32
	ksize uint32
	pgid  pgid
}

// leafPageElement represents a node on a leaf page.
// 叶子节点既存储key，也存储value
type leafPageElement struct {
	flags uint32
	pos   uint32
	ksize uint32
	vsize uint32
}

const maxAllocSize = 0x7FFFFFFF

// 叶子节点的key
// key returns a byte slice of the node key.
func (n *leafPageElement) key() []byte {
	buf := (*[maxAllocSize]byte)(unsafe.Pointer(n))
	// pos~ksize
	fmt.Println(buf[:unsafe.Sizeof(leafPageElement{})],len(buf))
	return (*[maxAllocSize]byte)(unsafe.Pointer(&buf[n.pos]))[:n.ksize:n.ksize]
}

// 叶子节点的value
// value returns a byte slice of the node value.
func (n *leafPageElement) value() []byte {
	buf := (*[maxAllocSize]byte)(unsafe.Pointer(n))
	fmt.Println(buf[:unsafe.Sizeof(leafPageElement{})],len(buf))
	// key:pos~ksize
	// value:pos+ksize~pos+ksize+vsize
	return (*[maxAllocSize]byte)(unsafe.Pointer(&buf[n.pos+n.ksize]))[:n.vsize:n.vsize]
}

func main() {
	// testBoltdb()

	// testUnsafeMethod()

	// fmt.Println(bolt.bucketHeaderSize)
	// fmt.Println(unsafe.Sizeof(bolt.Bucket{}))
}

func testUnsafeMethod() {
	fmt.Println("pageHeaderSize:", pageHeaderSize)
	fmt.Println("branchPageElementSize:", branchPageElementSize)
	fmt.Println("leafPageElementSize:", leafPageElementSize)
	fmt.Println("id:", unsafe.Sizeof((*page)(nil).id), unsafe.Offsetof((*page)(nil).id))
	fmt.Println("flags:", unsafe.Sizeof((*page)(nil).flags), unsafe.Offsetof((*page)(nil).flags))
	fmt.Println("count:", unsafe.Sizeof((*page)(nil).count), unsafe.Offsetof((*page)(nil).count))
	fmt.Println("overflow:", unsafe.Sizeof((*page)(nil).overflow), unsafe.Offsetof((*page)(nil).overflow))
	fmt.Println("ptr:", unsafe.Sizeof((*page)(nil).ptr), unsafe.Offsetof((*page)(nil).ptr))

	leaf := leafPageElement{
		flags: 0,
		pos:   10,
		ksize: 20,
		vsize: 5,
	}
	fmt.Println("key:", leaf.key(), "val:", leaf.value())
}

func testBoltdb() {
	db, err := bolt.Open("./my.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("user"))
		if err != nil {
			log.Fatalf("CreateBucketIfNotExists err:%s", err.Error())
			return err
		}
		if err = bucket.Put([]byte("hello"), []byte("world")); err != nil {
			log.Fatalf("bucket Put err:%s", err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalf("db.Update err:%s", err.Error())
	}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("user"))
		val := bucket.Get([]byte("hello"))
		log.Printf("the get val:%s", val)
		val = bucket.Get([]byte("hello2"))
		log.Printf("the get val2:%s", val)
		return nil
	})
	if err != nil {
		log.Fatalf("db.View err:%s", err.Error())
	}
}
