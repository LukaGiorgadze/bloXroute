package store

import (
	"fmt"
	_ "net/http/pprof"
	"sync"
)

type item2 struct {
	key  string
	val  string
	next *item2
	prev *item2
}

type LinkedList struct {
	head          *item2
	tile          *item2
	size          int
	lock          *sync.RWMutex
	fileLock      *sync.Mutex
	outputFilPath string
}

func NewLinkedList(mu *sync.RWMutex, mu2 *sync.Mutex, outputFilPath string) *LinkedList {
	return &LinkedList{
		lock:          mu,
		fileLock:      mu2,
		outputFilPath: outputFilPath,
	}
}

func (ll *LinkedList) Add(key, val string) bool {

	new := &item2{
		key: key,
		val: val,
	}

	if ll.head == nil {
		ll.head = new
		ll.tile = new
	} else {
		new.prev = ll.tile
		ll.tile.next = new
		ll.tile = new
	}

	ll.size++

	return true
}

func (ll *LinkedList) Get(key string) (string, bool) {

	current := ll.head

	for ; current != nil; current = current.next {
		if current.key == key {
			return current.val, true
		}
	}

	return "", false
}

func (ll *LinkedList) Remove(key string) bool {

	current := ll.head

	for ; current != nil; current = current.next {
		if current.key == key {
			if current.prev != nil {
				current.prev.next = current.next
			} else {
				ll.head = current.next
			}

			if current.next != nil {
				current.next.prev = current.prev
			} else {
				ll.tile = current.prev
			}

			ll.size--
			return true
		}
	}

	return false

}

func (ll *LinkedList) GetAll() []string {

	current := ll.head
	result := make([]string, ll.size)

	for ; current != nil; current = current.next {
		result = append(result, fmt.Sprintf("%s: %s \n", current.key, current.val))
	}

	return result
}

func (ll *LinkedList) Lock() *sync.RWMutex {
	return ll.lock
}

func (ll *LinkedList) FileLock() *sync.Mutex {
	return ll.fileLock
}

func (ll *LinkedList) GetOutputFilePath() string {
	return ll.outputFilPath
}
