package gateway

import (
	"github.com/Shopify/sarama"
	"log"
	"github.com/hyperdelta/nyduscanal/parser"
)

func StartGmarketNydusCanal(stop <-chan bool, address string) <-chan string {
	consumer, err := sarama.NewConsumer([]string{address}, nil)
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

				out <- parser.GmarketAddOrderParser(msg.Value)
			case isStop := <-stop:
				if(isStop) {
					break ConsumerLoop
				}
			}

		}
	}()

	return out
}