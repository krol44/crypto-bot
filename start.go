package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":3434", nil))
}

func send(chatId string, text string) {
	resp, _ := http.Get("https://api.telegram.org/xxxxxxxx/sendMessage?chat_id=" + chatId + "&text=" + text)
	fmt.Printf("%s", resp)
}

func handler(w http.ResponseWriter, r *http.Request) {
	stringSlice := strings.Split(r.URL.Path[1:], "/")

	go runObserver(stringSlice)

	fmt.Fprintf(w, "ok")
}

func runObserver(data []string) {

	fmt.Println(data)
	fmt.Println(fmt.Sprintf("wss://stream.binance.com/ws/%s@trade", data[0]))

	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://stream.binance.com/ws/%s@trade", data[0]), nil)

	//fmt.Println(conn)
	//fmt.Println(res)
	fmt.Println(err)

	type ResultJson struct {
		Price string `json:"p"`
	}

	for {
		_, message, readErr := conn.ReadMessage()
		if readErr != nil {
			fmt.Println(readErr)
			return
		}

		//fmt.Printf("%s", message)

		var dataFromJson ResultJson

		json.Unmarshal(message, &dataFromJson)

		//fmt.Println(dataFromJson.Price)

		floatOut, _ := strconv.ParseFloat(dataFromJson.Price, 42)
		floatUser, _ := strconv.ParseFloat(data[2], 42)

		if floatOut >= floatUser {
			send(data[3], "Binance notify: "+data[0]+" up to "+data[2])

			fmt.Println("send")
			return
		}

		//fmt.Print(".")
	}
}

func errorLog(errorInfo error) {
	if errorInfo != nil {
		fmt.Println(errorInfo)
	}
}
