package workers

import (
	"log"
	"os"
	"strings"
)

type FileWriter struct {
	// Data received by reader channels
	Data          chan string
	workersConfig *WorkersConfig
}

func NewFileWriter(buf uint8, cfg *WorkersConfig) *FileWriter {
	return &FileWriter{
		// It's recommended to have SemaphoreReaders number capacity in Data channel to not keep
		// reader goroutines blocked until one FileWriter gouroutine reads the data.
		Data:          make(chan string, buf),
		workersConfig: cfg,
	}
}

// FileWriter is a worker that writes (by appending) data into log, after receiving msgs from the reader channels.
func (s *FileWriter) FileWriterWorker() {

	defer close(s.Data)

	for str := range s.Data {

		// s.workersConfig.Store.FileLock().Lock()
		// Write by appending
		f, err := os.OpenFile(s.workersConfig.Store.GetOutputFilePath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
			return
		}

		// The reason of using strings.Builder instead of string concatenation is
		// that string is immutable and concatenation allocates memory each time,
		// which may lead to performance issues. Using strings.Builder, we can efficiently build the
		// string by appending each piece to the buffer without creating intermediate copies of the string.
		var sb strings.Builder
		sb.WriteString(str)
		sb.WriteString("\n")

		_, err = f.Write([]byte(sb.String()))
		if err != nil {
			log.Fatal(err)
			return
		}
		// s.workersConfig.Store.FileLock().Unlock()

		f.Close()
	}
}
