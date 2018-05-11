package gohq

import (
	"testing"
	"fmt"
)

func TestTrivia(t *testing.T) {
	sock, _ := HQDebug()
	for {
		bytes, _ := sock.Read()
		if summary := sock.ParseQuestionSummary(bytes); summary != nil {
			fmt.Println(summary.EliminatedPlayersCount)
		} else if stats := sock.ParseStats(bytes); stats != nil {
			fmt.Println(stats.ViewerCounts.Playing)
		} else if question := sock.ParseQuestion(bytes); question != nil {
			fmt.Println(question.Answers)
		}
	}
}
