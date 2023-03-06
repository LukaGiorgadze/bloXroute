package workers

import (
	"github.com/LukaGiorgadze/bloXroute/internal/client"
	"github.com/LukaGiorgadze/bloXroute/internal/models"
)

const (
	ADD_ITEM    = string(client.ItemMutateAddSubject)
	DELETE_ITEM = string(client.ItemMutateDeleteSubject)
)

// Once is a struct that represents a single worker that can process one item at a time.
type OnceMutator struct {

	// Queue is a channel that is used to send messages to the worker.
	Queue chan *models.Msg

	// workersConfig is a pointer to a WorkersConfig struct that is shared among all workers.
	workersConfig *WorkersConfig
}

func NewOnceMutator(cfg *WorkersConfig) *OnceMutator {
	return &OnceMutator{
		Queue:         make(chan *models.Msg, 1),
		workersConfig: cfg,
	}
}

// MutatorWorker is a function that listens for messages on the Queue channel and performs mutations on the workersConfig store.
// If the subject is ADD_ITEM, it adds the map item to the workersConfig store.
// If the subject is DELETE_ITEM, it removes the map item from the workersConfig store.
// The function is designed to run indefinitely, waiting for messages on the Queue channel.
func (o *OnceMutator) MutatorWorker() {
	for {
		item := <-o.Queue

		switch item.Subject {
		case ADD_ITEM:
			o.workersConfig.Store.Lock().Lock()
			_ = o.workersConfig.Store.Add(item.Key, item.Value)
			o.workersConfig.Store.Lock().Unlock()

		case DELETE_ITEM:
			o.workersConfig.Store.Lock().Lock()
			_ = o.workersConfig.Store.Remove(item.Key)
			o.workersConfig.Store.Lock().Unlock()

		}
	}
}
