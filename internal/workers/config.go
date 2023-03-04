package workers

import "github.com/LukaGiorgadze/bloXroute/internal/store"

// Global/shared configuration for workers
type WorkersConfig struct {
	Store store.IStore
}
