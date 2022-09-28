package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/gorilla/websocket"
)

var url = "wss://eth-mainnet.ws.alchemyapi.io/ws/demo"

var urlPolygon = "wss://polygon-mainnet.g.alchemy.com/v2/demo"

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		_, bMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// unmarshal the message
		var ethlog ethtypes.Log
		err = json.Unmarshal(bMsg, &ethlog)
		if err != nil {
			fmt.Println("Unmash:", err)
			return
		}

		// print out that message for clarity
		fmt.Println("Message:", bMsg)

	}
}

// createSchema creates database schema for ethtypes.Log
func createSchema(db *pg.DB) error {
	err := db.Model((*ethtypes.Log)(nil)).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {

	// web socket
	c, err := ethclient.Dial(urlPolygon)
	if err != nil {
		fmt.Println("Error start:", err)
		return
	}
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// raw request
	// rawReq := `{"jsonrpc":"2.0","id": 1, "method": "eth_subscribe", "params": ["logs"]}`

	// bReq := []byte(rawReq)

	// send request
	logChan := make(chan ethtypes.Log)

	sub, err := c.SubscribeFilterLogs(context.Background(), ethereum.FilterQuery{}, logChan)
	if err != nil {
		fmt.Println("Error subscription:", err)
		return
	}

	// connect to local db
	db := pg.Connect(&pg.Options{
		Network:  "tcp",
		Addr:     "127.0.0.1:5432",
		User:     "user",
		Password: "local-db",
		Database: "wsdb",
	})
	defer db.Close()

	// create log schema
	err = createSchema(db)
	if err != nil {
		fmt.Println("Error create:", err)
	}
	fmt.Println("Created:", db)

	/*
		vLog := ethtypes.Log{
			Address:     common.HexToAddress("0xecf8f87f810ecf450940c9f60066b4a7a501d6a7"),
			BlockHash:   common.HexToHash("0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056"),
			BlockNumber: 2019236,
			Data:        []byte{},
			Index:       2,
			TxIndex:     3,
			TxHash:      common.HexToHash("0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e"),
			Topics: []common.Hash{
				common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
			},
			Removed: true,
		}
		_, err = db.Model(&vLog).Insert()
		if err != nil {
			panic(err)
		}
		fmt.Println("Inserted:", db)
	*/

	// read the messages and insert the logs in DB
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
			return
		case vLog := <-logChan:
			_, err = db.Model(&vLog).Insert()
			if err != nil {
				panic(err)
			}
			fmt.Println(vLog) // pointer to event log
		}
	}

	// reader(c)

}
