package workers

import (
	"fmt"
	"strings"

	"github.com/LukaGiorgadze/bloXroute/internal/models"
)

// SemaphoreReader is a structure that limits the maximum number of concurrent readers.
type SemaphoreReader struct {
	queue         chan struct{}
	workersConfig *WorkersConfig
}

func NewSemaphoreReader(max uint8, cfg *WorkersConfig) *SemaphoreReader {
	return &SemaphoreReader{
		queue:         make(chan struct{}, max),
		workersConfig: cfg,
	}
}

// Acquire acquires a resource from the channel.
func (s *SemaphoreReader) Acquire() {
	s.queue <- struct{}{}
}

// Release releases a resource back to the channel.
func (s *SemaphoreReader) Release() {
	<-s.queue
}

// ReadAll reads all items safely in the store using RLock.
func (s *SemaphoreReader) ReadAll(fileWriterCh chan<- string) {

	defer s.Release()

	s.workersConfig.Store.Lock().RLock()
	items := s.workersConfig.Store.GetAll()
	s.workersConfig.Store.Lock().RUnlock()

	str := strings.Join(items, ",")

	// Print data in the server's stdout
	fmt.Println(str)

	// Send data to the file writer channel, so FileWriter worker can start it's job.
	fileWriterCh <- str

}

// ReadOne reads one item from the store by key.
func (s *SemaphoreReader) ReadOne(item *models.Msg, fileWriterCh chan<- string) {

	defer s.Release()

	s.workersConfig.Store.Lock().RLock()
	val, ok := s.workersConfig.Store.Get(item.Key)
	s.workersConfig.Store.Lock().RUnlock()
	if !ok {
		fmt.Println(item.Key, "= no data")
		return
	}

	// Build the string to be sent to the channel.
	// The reason of using strings.Builder instead of string concatenation is
	// that string is immutable and concatenation allocates memory each time,
	// which may lead to performance issues. Using strings.Builder, we can efficiently build the
	// string by appending each piece to the buffer without creating intermediate copies of the string.
	var sb strings.Builder
	sb.WriteString(item.Key)
	sb.WriteString("=")
	sb.WriteString(val)

	// Print data in the server's stdout
	fmt.Println(sb.String())

	// Send data to the file writer channel, so FileWriter worker can start it's job.
	fileWriterCh <- sb.String()

}
