package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// type Request struct {
// 	JSONRPC string `json:"jsonrpc"`
// 	Method  string `json:"method"`
// 	ID      jsonrpc.ID     `json:"id"`
// 	Params  jsonrpc.Params `json:"params"`
// }

// type Response struct {
// 	JSONRPC string `json:"jsonrpc"`
// 	Method  string `json:"method"`
// 	ID      jsonrpc.ID     `json:"id"`
// 	Params  jsonrpc.Params `json:"params"`
// }

var url = "wss://eth-mainnet.ws.alchemyapi.io/ws/demo"

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		_, bMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// unmarshal the message
		// var ethlog ethtypes.Log
		// err = json.Unmarshal(bMsg, &ethlog)
		// if err != nil {
		// 	fmt.Println("Unmash:", err)
		// 	return
		// }
		
		// print out that message for clarity
		fmt.Println("Message:", bMsg)

	}
}

func main() {

	// web socket
	c, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// request
	rawReq := `{"jsonrpc":"2.0","id": 1, "method": "eth_subscribe", "params": ["logs"]}`

	bReq := []byte(rawReq)

	// send request

	err = c.WriteMessage(1, bReq)

	if err != nil {
		fmt.Println("Error write:", err)
		return
	}

	reader(c)

}
