package ws

type MessageType string

const (
	MsgBroadcast MessageType = "broadcast"
	MsgHello     MessageType = "hello"
	MsgGoodbyte  MessageType = "goodbye"
)

type Message struct {
	Type    MessageType `json:"type"`
	Content string      `json:"content"`
}
