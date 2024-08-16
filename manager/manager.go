package manager

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Client struct {
	Conn  *websocket.Conn
	email string
}

var (
	upgrader = websocket.Upgrader{}
)
var (
	Clients = make(map[string]*Client)
	Mu      sync.Mutex
)

func Connect(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	email := c.Get("email").(string)
	client := Clients[email]
	if client == nil {
		AddClient(email, ws)
		SendNotificationForEmail(email)
	}
	defer func() {
		RemoveClient(email)
	}()
	for {

		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}

func AddClient(email string, conn *websocket.Conn) {
	Mu.Lock()
	defer Mu.Unlock()
	Clients[email] = &Client{Conn: conn, email: email}
}

func RemoveClient(email string) {
	log.Println("did we rly enter hereee? inside the RemoveClient")
	Mu.Lock()
	defer Mu.Unlock()
	delete(Clients, email)
}

func SendClientNotification(message string) error {
	log.Println("did wee enter here", Clients)
	Mu.Lock()
	defer Mu.Unlock()

	// Iterate over all clients in the map
	for _, client := range Clients {
		notification := []byte(message)
		log.Println("did wee enter here", client.email)
		err := client.Conn.WriteMessage(websocket.TextMessage, notification)
		if err != nil {
			log.Printf("Error sending notification: %v\n", err)
			return err
		} else {
			log.Printf("Notification sent to shop: %s\n", client.email)
			return nil
		}
	}
	return nil
}

func SendNotificationForEmail(email string) {
	Mu.Lock()
	defer Mu.Unlock()
	client := Clients[email]
	if client == nil {
		return
	}
	err := client.Conn.WriteMessage(websocket.TextMessage, []byte("refetch"))
	if err != nil {
		log.Printf("Error sending notification: %v\n", err)
	} else {
		log.Printf("Notification sent to shop: %s\n", client.email)
	}
}
