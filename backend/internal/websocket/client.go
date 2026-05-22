package websocket

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Client struct {
	Room string
	Conn *websocket.Conn
	Send chan []byte
}

func HandleWebSocket(c *fiber.Ctx) error {
	room := c.Params("trackingNumber", "global")
	if websocket.IsWebSocketUpgrade(c) {
		return websocket.New(func(conn *websocket.Conn) {
			client := &Client{
				Room: room,
				Conn: conn,
				Send: make(chan []byte, 256),
			}
			DefaultHub.Register(client)
			defer DefaultHub.Unregister(client)

			go func() {
				for msg := range client.Send {
					if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
						return
					}
				}
			}()

			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					break
				}
			}
		})(c)
	}
	return c.SendStatus(400)
}
