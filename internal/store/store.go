package store

import "sync"

// The IStore interface defines a set of methods that any
// memory storage/data structure implementation should implement.
// It provides an abstract instance that can be used for different data structures.
// By using this interface, we can switch between different data structures
// without changing the rest of the code that uses it, or just adding new ones.
type IStore interface {
	Add(string, string) bool
	Remove(string) bool
	Get(string) (string, bool)
	GetAll() []string
	Lock() *sync.RWMutex
	FileLock() *sync.Mutex
	GetOutputFilePath() string
}
