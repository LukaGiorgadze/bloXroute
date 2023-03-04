package models

// Item model represents a key-value pair in a JSON format.
type Item struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

// Msg model is used for communication on a messaging system.
// It contains an Item field and a Subject field, which specifies the subject/topic
// of the message being sent or received.
// This model can be used to serialize and deserialize data for messaging between systems.
type Msg struct {
	Item
	Subject string
}
