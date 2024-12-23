package main

import (
	"fmt"

	"github.com/BlackestDawn/learn-pub-sub-starter/internal/gamelogic"
	"github.com/BlackestDawn/learn-pub-sub-starter/internal/pubsub"
	"github.com/BlackestDawn/learn-pub-sub-starter/internal/routing"
)

func handlerGameLog() func(routing.GameLog) pubsub.AckType {
	return func(msg routing.GameLog) pubsub.AckType {
		defer fmt.Print("> ")
		err := gamelogic.WriteLog(msg)
		if err != nil {
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}
