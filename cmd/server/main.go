package main

import (
	"fmt"

	"github.com/BlackestDawn/learn-pub-sub-starter/internal/gamelogic"
	"github.com/BlackestDawn/learn-pub-sub-starter/internal/pubsub"
	"github.com/BlackestDawn/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	conn, err := amqp.Dial(pubsub.AmqpServer)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	pubCH, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		fmt.Sprintf("%s.*", routing.GameLogSlug),
		pubsub.QueueTypeDurable,
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Printf("Queue %v declared and bound.\n", queue.Name)

	gamelogic.PrintServerHelp()
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		command := input[0]
		switch command {
		case "pause":
			fmt.Println("Pausing game...")
			err = pubsub.PublishJSON(pubCH, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				fmt.Println(err)
			}
		case "resume":
			fmt.Println("Resuming game...")
			err = pubsub.PublishJSON(pubCH, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			if err != nil {
				fmt.Println(err)
			}
		case "quit":
			fmt.Println("Exiting...")
			return
		case "help":
			gamelogic.PrintServerHelp()
		default:
			fmt.Println("Unknown command:", command)
			continue
		} // end switch
	} // end for
}
