package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	node := maelstrom.NewNode()
	node.Handle("echo", func(msg maelstrom.Message) error {
		var msgBody map[string]any
		json.Unmarshal(msg.Body, &msgBody)
		msgBody["type"] = "echo_ok"
		return node.Reply(msg, msgBody)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
