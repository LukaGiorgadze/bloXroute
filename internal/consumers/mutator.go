package consumers

import (
	"github.com/LukaGiorgadze/bloXroute/configs"
	"github.com/LukaGiorgadze/bloXroute/internal/client"
	"github.com/LukaGiorgadze/bloXroute/internal/store"
	"github.com/LukaGiorgadze/bloXroute/internal/workers"
	"github.com/nats-io/nats.go"
)

// ItemMutateHandler holds some necessary instances in order to run jobs and communicate effectively.
type ItemMutateHandler struct {

	// Global configuration (envs) of app/server configuration.
	cfg *configs.Config

	// Store where all the data are stored during mutation.
	store store.IStore

	// Once is a pattern to run mutation worker only once, since we want to keep maintain,
	// ordering of items added/delete.
	onceMutator *workers.OnceMutator
}

func NewItemMutateHandler(cfg *configs.Config, store store.IStore) *ItemMutateHandler {

	// Configure worker, set store where it should store data.
	workersConfig := &workers.WorkersConfig{
		Store: store,
	}

	// Inizialize worker and assign it to the ItemMutateHandler struct,
	// so it can be used later in handler or consumer.
	onceMutator := workers.NewOnceMutator(workersConfig)

	return &ItemMutateHandler{
		cfg,
		store,
		onceMutator,
	}
}

// Handler runs MutatorWorker routine and returns consumer for reading messages form subscription.
func (ih *ItemMutateHandler) Handler() func(*nats.Msg) {

	// Run MutatorWorker and wait for the messages in another "thread".
	go ih.onceMutator.MutatorWorker()

	return ih.consumer()
}

// consumer reads messages from subscription and sends it to `onceMutator.Queue` if channel is not blocked.
// `MutatorWorker()` receives this message and processes mutation, so after that `onceMutator.Queue`
// becomes unblocked and available for the next cycle.
// The idea is to have 1 processing at the time to keep ordering of insertion/deletion in the store.
func (ih *ItemMutateHandler) consumer() func(msg *nats.Msg) {

	const (
		ADD_ITEM    = string(client.ItemMutateAddSubject)
		DELETE_ITEM = string(client.ItemMutateDeleteSubject)
	)

	return func(msg *nats.Msg) {
		// There might be chance that msg.Subject does not contain any of them,
		// if so - we skip it.
		if msg.Subject != ADD_ITEM && msg.Subject != DELETE_ITEM {
			return
		}

		// We are sending converted message to the onceMutator.Queue channel.
		// The mutator worker will then retrieve the message from the queue and process it.
		// The onceMutator.Queue is a buffered channel with a capacity of 1,
		// meaning that any new incoming messages from the subscription will be blocked until the mutator
		// worker has finished processing the current message. So we can maintain ordering of insertion/deletion.
		m := msgToStruct(msg)
		ih.onceMutator.Queue <- m
	}
}
