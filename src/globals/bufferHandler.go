package globals

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type connNode struct {
	conn     *websocket.Conn
	prevNode *connNode
	nextNode *connNode
}

var buffer strings.Builder
var updateBuffer []Update
var ticker *time.Ticker
var wsConns *connNode = nil

func AppendConnection(_conn *websocket.Conn) *connNode {
	newNode := connNode{
		prevNode: nil,
		nextNode: nil,
		conn:     _conn,
	}

	if wsConns == nil {
		wsConns = &newNode
		return &newNode
	}

	nextConn := wsConns

	for nextConn.nextNode != nil {
		nextConn = nextConn.nextNode
	}

	newNode.prevNode = nextConn
	nextConn.nextNode = &newNode

	return &newNode
}

func DeleteConnection(_node *connNode) {
	prevNode := _node.prevNode
	nextNode := _node.nextNode
	if prevNode != nil {
		prevNode.nextNode = nextNode
	} else {
		wsConns = nextNode
	}

	if nextNode != nil {
		nextNode.prevNode = prevNode
	}

	_node.conn = nil
	_node.nextNode = nil
	_node.prevNode = nil
}

func AppendBuffer(update string) {
	buffer.WriteString(update)
}

func GetBuffer() string {
	return buffer.String()
}

func AddUpdate(update Update) {
	updateBuffer = append(updateBuffer, update)
}

func clearUpdates() {
	updateBuffer = make([]Update, 0)
}

func GetUpdates() []Update {
	return updateBuffer
}

func StartLogging() {
	ticker = time.NewTicker(500 * time.Millisecond)

	for range ticker.C {
		if buffer.Len() > 0 {

			var chunk string
			if buffer.Len() > 2000 {
				fullBuffer := buffer.String()
				chunk = fullBuffer[:2000]

				split := strings.Split(chunk, "\n")
				if len(split) > 1 {
					split = split[:len(split)-1]
				}
				chunk = strings.Join(split, "\n")

				buffer.Reset()
				buffer.WriteString(fullBuffer[len(chunk):])
			} else {
				chunk = buffer.String()
				buffer.Reset()
			}

			bot := GetDiscord()

			channel := os.Getenv("CONSOLE_CHANNEL")

			if channel == "" {
				log.Fatal("CONSOLE CHANNEL EMPTY IN .env")
			}

			node := wsConns
			for connsActive := node != nil; connsActive; connsActive = node != nil {
				if node.conn != nil {
					err := node.conn.WriteMessage(websocket.TextMessage, []byte(chunk))
					if err != nil {
						brokenNode := node
						node = node.nextNode
						DeleteConnection(brokenNode)
						continue
					}
					node = node.nextNode
				} else {
					brokenNode := node
					node = node.nextNode
					DeleteConnection(brokenNode)
				}
			}

			bot.ChannelMessageSend(channel, chunk)
		}
	}
}
