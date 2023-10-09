package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Message struct {
	Type      string `json:"type"`
	Message   int    `json:"message"`
	MessageId *int   `json:"msg_id"`
}

type Topology struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

func main() {
	node := maelstrom.NewNode()
	var messages []int
	var topology map[string][]string

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var msgBody Message
		var msgResp = make(map[string]any)
		json.Unmarshal(msg.Body, &msgBody)
		msgResp["type"] = "broadcast_ok"
		msgResp["msg_id"] = &msgBody.MessageId

		var found bool = false
		for _, message := range messages {
			if message == msgBody.Message {
				found = true
				break
			}
		}

		if found == false {
			messages = append(messages, msgBody.Message)
			nodes := topology[node.ID()]
			var gossipResp = make(map[string]any)
			gossipResp["type"] = "broadcast"
			for i := 0; i < len(nodes); i++ {
				reciever := nodes[i]
				for i := 0; i < len(messages); i++ {
					gossipResp["messages"] = messages[i]
					node.Send(reciever, gossipResp)
				}
			}
		}

		if &msgBody.MessageId != nil {
			return node.Reply(msg, msgResp)
		}

		return nil
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		var msgBody map[string]any
		json.Unmarshal(msg.Body, &msgBody)
		msgBody["type"] = "read_ok"
		msgBody["messages"] = messages
		return node.Reply(msg, msgBody)
	})

	node.Handle("topology", func(msg maelstrom.Message) error {
		var msgBody Topology
		var msgResp = make(map[string]any)
		json.Unmarshal(msg.Body, &msgBody)
		msgResp["type"] = "topology_ok"
		topology = msgBody.Topology
		return node.Reply(msg, msgResp)
	})

	node.Handle("broadcast_ok", func(msg maelstrom.Message) error {
		var msgResp = make(map[string]any)
		json.Unmarshal(msg.Body, &msgResp)
		msgResp["type"] = "broadcast_ok"
		return node.Reply(msg, msgResp)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
