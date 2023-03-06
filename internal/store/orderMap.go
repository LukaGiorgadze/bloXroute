package store

import (
	"fmt"
	"sync"
)

type item struct {
	key   string
	value string
	next  *item
	prev  *item
}

// LinkedList
// The OrderedMap struct holds a head pointer to the first item in the map,
// a tail pointer to the last item in the map, and a map named "items" that holds pointers
// to all the items in the map for quick access. It also holds a size value to keep
// track of the number of items in the map, a pointer to a sync.RWMutex used for locking access to the map,
// and a string that represents the path to the output file.
type OrderedMap struct {
	head  *item
	tail  *item
	size  int
	items map[string]*item
	// Lock method returns the sync.RWMutex used to lock access to the ordered map data structure.
	lock *sync.RWMutex
	// Lock method returns the sync.Mutex used to lock access to the output file data.
	fileLock *sync.Mutex
	// Output file path where data will be saved.
	// should be set during initialization of OrderedMap (NewOrderedMap).
	outputFilPath string
}

func NewOrderedMap(mu *sync.RWMutex, mu2 *sync.Mutex, outputFilPath string) *OrderedMap {
	return &OrderedMap{
		lock:          mu,
		fileLock:      mu2,
		items:         make(map[string]*item),
		outputFilPath: outputFilPath,
	}
}

func (om *OrderedMap) Add(key string, value string) (ok bool) {
	if _, exists := om.items[key]; exists {
		return
	}

	newItem := &item{key: key, value: value}

	if om.head == nil {
		om.head = newItem
		om.tail = newItem
	} else {
		om.tail.next = newItem
		newItem.prev = om.tail
		om.tail = newItem
	}

	om.items[key] = newItem
	om.size++

	return !ok
}

func (om *OrderedMap) Remove(key string) (ok bool) {
	item, exists := om.items[key]
	if !exists {
		return
	}

	if item == om.head {
		om.head = item.next
		if om.head != nil {
			om.head.prev = nil
		}
	} else if item == om.tail {
		om.tail = item.prev
		if om.tail != nil {
			om.tail.next = nil
		}
	} else {
		item.prev.next = item.next
		item.next.prev = item.prev
	}

	delete(om.items, key)
	om.size--

	return !ok
}

func (om *OrderedMap) Get(key string) (value string, ok bool) {
	item, ok := om.items[key]
	if !ok {
		return
	}

	return item.value, ok
}

func (om *OrderedMap) GetAll() []string {
	result := make([]string, om.size)
	index := 0
	for item := om.head; item != nil; item = item.next {
		result[index] = fmt.Sprintf("(%s=%s)", item.key, item.value)
		index++
	}
	return result
}

func (om *OrderedMap) Lock() *sync.RWMutex {
	return om.lock
}

func (om *OrderedMap) FileLock() *sync.Mutex {
	return om.fileLock
}

func (om *OrderedMap) GetOutputFilePath() string {
	return om.outputFilPath
}
