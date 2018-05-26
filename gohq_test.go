package gohq

import (
	"testing"
	"log"
	"fmt"
	"strconv"
)

func TestHQ(t *testing.T) {
	game, err := DebugHQ()
	if err != nil {
		log.Fatal(err)
	}

	for {
		bytes, err := game.Read()
		if err != nil {
			log.Fatal(err)
		}

		if stats := game.ParseBroadcastStats(bytes); stats != nil {
			continue
		} else if message := game.ParseChatMessage(bytes); message != nil {
			continue
		} else if gameStatus := game.ParseGameStatus(bytes); gameStatus != nil {
			fmt.Println("You have joined, the prize is", gameStatus.Prize)
		} else if question := game.ParseQuestion(bytes); question != nil {
			fmt.Println("Question Incoming, you have", strconv.Itoa(question.TimeLeftMs)+"ms to answer it!")
		} else if questionClosed := game.ParseQuestionClosed(bytes); questionClosed != nil {
			fmt.Println("The question", questionClosed.QuestionID, "is over!")
		} else if questionSummary := game.ParseQuestionSummary(bytes); questionSummary != nil {
			fmt.Println(questionSummary.EliminatedPlayersCount, "players have been eliminated,", questionSummary.AdvancingPlayersCount, "players remain.")
		} else if questionFinished := game.ParseQuestionFinished(bytes); questionFinished != nil {
			fmt.Println("The question", questionFinished.QuestionID, "has finished.")
		} else {
			fmt.Println("This is an unknown message:", string(bytes))
		}
	}
}
