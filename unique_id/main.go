package main

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	node := maelstrom.NewNode()
	node.Handle("generate", func(msg maelstrom.Message) error {
		var msgBody map[string]any
		json.Unmarshal(msg.Body, &msgBody)
		uuid := uuid.New().String()
		msgBody["id"] = uuid
		msgBody["type"] = "generate_ok"
		return node.Reply(msg, msgBody)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
