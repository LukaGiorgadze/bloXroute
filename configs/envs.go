package configs

type Config struct {
	NatsURL                    string `env:"NATS_URL" envDefault:"0.0.0.0:4222"`
	NatsUser                   string `env:"NATS_USER" envDefault:"dummy"`
	NatsPass                   string `env:"NATS_PASS" envDefault:"password"`
	SemaphoreReadMaxGoroutines string `env:"SEM_READ_MAX_GR" envDefault:"10"`
	OutputFilePath             string `env:"OUTPUT_FILE_PATH" envDefault:"./output/items.log"`
}
