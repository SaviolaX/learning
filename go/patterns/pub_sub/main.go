package main

import (
	"fmt"
)

type Broker struct {
	topics map[string][]chan string
}

func main() {
	broker := NewBroker()

	aliceNews := broker.Subscribe("news")
	aliceWeather := broker.Subscribe("weather")

	bobNews := broker.Subscribe("news")
	bobSport := broker.Subscribe("sport")

	broker.Publish("news", "Go 1.23 released!")
	broker.Publish("sport", "Dynamo won!")
	broker.Publish("weather", "Sunny today!")

	fmt.Println("Alice news:", <-aliceNews)
	fmt.Println("Alice weather:", <-aliceWeather)
	fmt.Println("Bob news:", <-bobNews)
	fmt.Println("Bob sport:", <-bobSport)

	broker.Close()
}

func NewBroker() *Broker {
	return &Broker{
		topics: make(map[string][]chan string),
	}
}

func (b *Broker) Subscribe(topic string) <-chan string {
	newChannel := make(chan string, 1)
	b.topics[topic] = append(b.topics[topic], newChannel)
	return newChannel
}

func (b *Broker) Publish(topic string, entry string) {
	topicChannels := b.topics[topic]
	for _, x := range topicChannels {
		x <- entry
	}
}

func (b *Broker) Close() {
	topics := b.topics
	for _, x := range topics {
		for _, i := range x {
			close(i)
		}
	}
}
