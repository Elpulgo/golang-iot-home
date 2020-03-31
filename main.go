package main

import (
	"bufio"
	"fmt"
	"iot-home/endpoints"
	"iot-home/endpoints/greeting"
	"os"
)

func main() {
	endpoints.Serve()
	fmt.Println("Yippikayaj")
	fmt.Println("Enter text:")

	text := read_line()
	fmt.Println(text)

	fmt.Println(greeting.WelcomeText)
}

func read_line() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}
