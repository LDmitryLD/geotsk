package controller

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

/**
 * Ниже код не редактировать
 */

func (c *CourierController) Websocket(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	handleConnection(conn, c.MoveCourier)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type webSocketMessage struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type CourierMove struct {
	Direction int `json:"direction"`
	Zoom      int `json:"zoom"`
}

func deserializeMessage(msg []byte) (webSocketMessage, error) {
	var m webSocketMessage
	err := json.Unmarshal(msg, &m)
	if err != nil {
		return m, err
	}
	m.Data, err = json.Marshal(m.Data)

	return m, err
}

func handleConnection(conn *websocket.Conn, f func(webSocketMessage)) {
	defer conn.Close()
	for {
		messageType, rawMessage, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if messageType != websocket.TextMessage {
			continue
		}
		err = handleMessage(rawMessage, f)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func handleMessage(rawMessage []byte, f func(webSocketMessage)) error {
	m, err := deserializeMessage(rawMessage)
	if err != nil {
		return err
	}

	f(m)

	return nil
}
