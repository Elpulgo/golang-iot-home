package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var file *os.File
var mu sync.Mutex

func write(text string, severity string) {
	mu.Lock()
	file, err := os.OpenFile(
		"log.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644)

	if err != nil {
		log.Println(err)
	}

	logger := log.New(file, severity+": ", log.LstdFlags)
	logger.Println(text)

	defer file.Close()
	defer mu.Unlock()
}

func Info(text string) {
	fmt.Println("INFO: " + text)
	write(text, "INFO")
}

func Error(text string) {
	fmt.Println("ERROR: " + text)
	write(text, "ERROR")
}

func Fatal(text string) {
	fmt.Println("FATAL: " + text)
	write(text, "FATAL")
	os.Exit(1)
}
