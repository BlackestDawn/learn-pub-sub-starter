package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/BlackestDawn/learn-pub-sub-starter/internal/gamelogic"
	"github.com/BlackestDawn/learn-pub-sub-starter/internal/pubsub"
	"github.com/BlackestDawn/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	conn, err := amqp.Dial(pubsub.AmqpServer)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	pubCH, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		panic(err)
	}

	gameState := gamelogic.NewGameState(username)
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.QueueTypeTransient,
		handlerPause(gameState),
	)
	if err != nil {
		panic(err)
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+username,
		routing.ArmyMovesPrefix+".*",
		pubsub.QueueTypeTransient,
		handlerMove(gameState, pubCH),
	)
	if err != nil {
		panic(err)
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.WarRecognitionsPrefix,
		routing.WarRecognitionsPrefix+".*",
		pubsub.QueueTypeDurable,
		handlerWar(gameState, pubCH),
	)
	if err != nil {
		panic(err)
	}

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		command := input[0]
		switch command {
		case "spawn":
			err = gameState.CommandSpawn(input)
			if err != nil {
				fmt.Println(err)
			}
		case "move":
			move, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = pubsub.PublishJSON(
				pubCH,
				routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+"."+username,
				move,
			)
			if err != nil {
				fmt.Println(err)
			}
		case "quit":
			gamelogic.PrintQuit()
			return
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			if len(input) < 2 {
				fmt.Println("usage: spam <number>")
				continue
			}
			nr, err := strconv.Atoi(input[1])
			if err != nil {
				fmt.Println("unable to parse number:", err)
				continue
			}
			for i := 0; i < nr; i++ {
				pubsub.PublishGob(
					pubCH,
					routing.ExchangePerilTopic,
					routing.GameLogSlug+"."+username,
					routing.GameLog{
						CurrentTime: time.Now(),
						Username:    username,
						Message:     gamelogic.GetMaliciousLog(),
					},
				)
			}
		default:
			fmt.Println("Unknown command:", command)
			continue
		} // end switch
	} // end for
}
