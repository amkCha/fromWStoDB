package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/websocket"
)

// type Request struct {
// 	JSONRPC string `json:"jsonrpc"`
// 	Method  string `json:"method"`
// 	ID      jsonrpc.ID     `json:"id"`
// 	Params  jsonrpc.Params `json:"params"`
// }

// type Paparams struct {
// 	Subscription string `json:"subscription"`
// 	Result *ethtypes.Log `json:"result"`
// }

// type Response struct {
// 	JSONRPC string `json:"jsonrpc"`
// 	Method  string `json:"method"`
// 	Params Paparams `json:"params"`
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
	c, err := ethclient.Dial(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// request
	// rawReq := `{"jsonrpc":"2.0","id": 1, "method": "eth_subscribe", "params": ["logs"]}`

	// bReq := []byte(rawReq)

	// send request

	logChan := make(chan ethtypes.Log)

	sub, err := c.SubscribeFilterLogs(context.Background(), ethereum.FilterQuery{}, logChan)
	if err != nil {
		fmt.Println("Error subscription:", err)
		return
	}

	for {
		select {
			case err := <-sub.Err():
				log.Fatal(err)
				return 
			case vLog := <-logChan:
				fmt.Println(vLog) // pointer to event log
		}
	}

	// reader(c)

}
