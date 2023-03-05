package main

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"net/http"
	_ "net/http/pprof"

	"github.com/LukaGiorgadze/bloXroute/configs"
	"github.com/LukaGiorgadze/bloXroute/internal/client"
	"github.com/LukaGiorgadze/bloXroute/internal/consumers"
	"github.com/LukaGiorgadze/bloXroute/internal/store"
	"github.com/nats-io/nats.go"
)

func main() {

	// NewConfig loads and parses the environment variables into structs.
	// It uses the default tag to set values, which can be overwritten by setting environment
	cfg, err := configs.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Initializes the message client by establishing a connection with the messaging system.
	// The msgClient is of the IMessageClient interface type and can be replaced with other implementations
	// of messaging systems like RabbitMQ, Kafka, etc. It can also be mocked during testing.
	var msgClient client.IMessageClient = client.NewNatsClient(cfg.NatsURL, []nats.Option{nats.UserInfo(cfg.NatsUser, cfg.NatsPass)})
	err = msgClient.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = msgClient.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// A mutex is used by the mutator and accessor goroutines.
	var lock = sync.RWMutex{}

	// A mutex is used by the file writer goroutine.
	var fileLock = sync.Mutex{}

	// The unsafeStore is the global store (database) kept in memory and is accessible and modifiable by goroutines.
	// However, direct access or modification is not recommended without using locks,
	// as it can lead to race conditions and other synchronization issues.
	// Therefore, it is advised to use the `unsafeStore.Lock()` method to control access to the store to ensure safe concurrent operations.
	var unsafeStore store.IStore = store.NewOrderedMap(&lock, &fileLock, cfg.OutputFilePath)

	// The consumers in this application contain handlers, which are the first callbacks in the subscribe method.
	// These handlers can be used to write additional logic, initialize routines,
	// and perform other tasks before the consumers start processing messages.
	//
	// The itemMutateConsumer is responsible for watching messages sent to the client.ItemMutateSubject subject,
	// which are used to modify the state of the global store. These messages can be either ADD or DELETE signals,
	// and the consumer is subscribed to the item.mutate.* wildcard subject to receive any messages sent
	// to item.mutate.add or item.mutate.delete.
	//
	// For more information visit https://docs.nats.io/nats-concepts/subjects.
	itemMutateConsumer := consumers.NewItemMutateHandler(&cfg, unsafeStore)
	err = msgClient.Subscribe(client.ItemMutateSubject, itemMutateConsumer.Handler())
	if err != nil {
		log.Panic(err)
	}
	defer msgClient.Unsubscribe(client.ItemMutateSubject)

	// itemAccessConsumer reads data requested by the client and communicates with the fileWriter,
	// which is responsible for writing read outputs to a file.
	itemAccessConsumer := consumers.NewItemAccessHandler(&cfg, unsafeStore)
	err = msgClient.Subscribe(client.ItemGetSubject, itemAccessConsumer.Handler())
	if err != nil {
		log.Panic(err)
	}
	defer msgClient.Unsubscribe(client.ItemGetSubject)

	// Run pprof to visualize and analyze profiling data.
	if cfg.Pprof {
		log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
	} else {
		// Wait for interrupt signal and then close the application
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, os.Kill)
		<-c
	}
}
