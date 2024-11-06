package main

import (
	"bufio"
	"fmt"
	"log"
	"mts2024golang/seminar_8_06_11/internal/admin"
	"mts2024golang/seminar_8_06_11/internal/consumer"
	"mts2024golang/seminar_8_06_11/internal/producer"
	"os"
	"strings"
)

func main() {
	topic := "chat-room"
	brokers := []string{"localhost:19092"}

	a, err := admin.New(brokers)
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()
	ok, err := a.TopicExists(topic)
	if !ok {
		err = a.CreateTopic(topic)
		if err != nil {
			log.Fatal(err)
		}
	}

	userName := ""
	fmt.Print("Enter your userName: ")
	_, err = fmt.Scanln(&userName)
	if err != nil {
		log.Fatal(err)
	}

	p, err := producer.New(brokers, topic)
	if err != nil {
		log.Fatalln("Failed to create producer:", err)
	}
	defer p.Close()

	c, err := consumer.New(brokers, topic)
	if err != nil {
		log.Fatalln("Failed to create consumer:", err)
	}
	defer c.Close()

	go c.PrintMessages()
	fmt.Println("Connected. Press Ctrl+C to exit")
	reader := bufio.NewReader(os.Stdin)
	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		p.SendMessage(userName, message)
	}
}
