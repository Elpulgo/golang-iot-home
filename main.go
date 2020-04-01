package main

import (
	"bufio"
	"fmt"
	"iot-home/endpoints"
	"iot-home/endpoints/greeting"
	"iot-home/netatmo"
	"log"
	"net/http"
	"os"
)

func main() {
	netatmoRest := new(netatmo.Rest)
	endpoints.Init(netatmo.New(netatmoRest))

	fmt.Println(netatmoRest.GetCurrent())
	fmt.Println("Yippikayajdd")

	// fmt.Println("Enter text:")

	// text := read_line()
	// fmt.Println(text)

	fmt.Println(greeting.WelcomeText)

	log.Println("Listening on :3001...")
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func read_line() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}
