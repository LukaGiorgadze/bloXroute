package consumers

import (
	"encoding/json"
	"log"

	"github.com/LukaGiorgadze/bloXroute/internal/models"
	"github.com/nats-io/nats.go"
)

// Unmarshal the input message and convert it into our defined model/struct.
// In addition, assign any necessary properties to the model/struct.
func msgToStruct(msg *nats.Msg) (item *models.Msg) {
	err := json.Unmarshal(msg.Data, &item)
	if err != nil {
		log.Println(err)
	}
	item.Subject = msg.Subject

	return
}
