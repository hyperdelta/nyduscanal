package main

import (
	"github.com/Shopify/sarama"
	"log"
)

func buildNyduscanal(stop <-chan bool) <-chan string {
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		panic(err)
	}

	partitionConsumer, err := consumer.ConsumePartition("test", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}



	out := make(chan string)
	go func() {
		defer func() {
			close(out)
			if err := partitionConsumer.Close(); err != nil {
				log.Fatalln(err)
			}

		}()
		ConsumerLoop:
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				// todo : 이 부분에 단순이 값이 아니라 이쁜 json으로 넘겨야함
				out <- parseData(msg.Value)
			case isStop := <-stop:
				if(isStop) {
					break ConsumerLoop
				}
			}

		}
	}()

	return out
}