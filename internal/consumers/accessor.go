package consumers

import (
	"github.com/LukaGiorgadze/bloXroute/configs"
	"github.com/LukaGiorgadze/bloXroute/internal/client"
	"github.com/LukaGiorgadze/bloXroute/internal/store"
	"github.com/LukaGiorgadze/bloXroute/internal/workers"
	"github.com/nats-io/nats.go"
)

// ItemAccessHandler holds some necessary instances in order to run jobs and communicate effectively.
type ItemAccessHandler struct {
	// Global configuration (envs) of app/server configuration.
	cfg *configs.Config

	// Store where all the data are stored during mutation.
	store store.IStore

	semaphoreReader *workers.SemaphoreReader
	fileWriter      *workers.FileWriter
}

func NewItemAccessHandler(cfg *configs.Config, store store.IStore) *ItemAccessHandler {

	// Configure worker, set store where it should store data.
	workersConfig := &workers.WorkersConfig{
		Store: store,
	}

	// Inizialize workers and assign it to the ItemMutateHandler struct,
	// so it can be used later in handler or consumer.
	semaphoreReader := workers.NewSemaphoreReader(cfg.SemaphoreReadMaxGoroutines, workersConfig)

	// It's recommended to have SemaphoreReaders number capacity in Data channel to not keep
	// reader goroutines blocked until one FileWriter gouroutine reads the data.
	fileWriter := workers.NewFileWriter(cfg.SemaphoreReadMaxGoroutines, workersConfig)

	return &ItemAccessHandler{
		cfg,
		store,
		semaphoreReader,
		fileWriter,
	}
}

// Handler runs FileWriterWorker routine and returns consumer for reading messages form subscription.
func (ih *ItemAccessHandler) Handler() func(*nats.Msg) {

	// Run FileWriterWorker and wait for the messages in another "thread".
	if ih.store.GetOutputFilePath() != "" {
		go ih.fileWriter.FileWriterWorker()
	}

	return ih.consumer()
}

// The consumer function reads messages from subscription and acquires a limited resource from a channel.
// If the channel is blocked, the semaphoreReader worker won't be run,
// which implements the Semaphore concurrency pattern.
// This pattern is used to limit the number of goroutines running in parallel to a certain number,
// preventing them from overwhelming the system.
func (ih *ItemAccessHandler) consumer() func(msg *nats.Msg) {

	const (
		GET_ITEM  = string(client.ItemGetSubject)
		ITEM_LIST = string(client.ItemGetListSubject)
	)

	return func(msg *nats.Msg) {
		switch msg.Subject {

		case GET_ITEM:
			m := msgToStruct(msg)
			ih.semaphoreReader.Acquire()
			go ih.semaphoreReader.ReadOne(m, ih.fileWriter.Data)

		case ITEM_LIST:
			ih.semaphoreReader.Acquire()
			go ih.semaphoreReader.ReadAll(ih.fileWriter.Data)
		}
	}
}
