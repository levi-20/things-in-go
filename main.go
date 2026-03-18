package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Answer struct {
	No     string
	Ans    int
	Score  int
	Result string
}

func main() {

	data, err := os.ReadFile(filepath.Join("quiz.csv"))

	if err != nil {
		log.Fatalf("Error while reding the csv file %v", err)
	}

	csvReader := csv.NewReader(bytes.NewReader(data))

	records, e := csvReader.ReadAll()

	if e != nil {
		log.Fatal("error while parsing csv: ", e)
	}

	reader := bufio.NewReader(os.Stdin)
	answers := make([]Answer, 0, len(records)-1)
	fmt.Println("Quiz\n====")
	fmt.Println("You have 30 seconds to answer all questions\n=============================================")

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
quizLoop:
	for _, record := range records[1:] {

		fmt.Printf("\nProblem %s: %s\nAnswer = ", record[0], record[1])

		inputChannel := make(chan string, 1)

		go func() {
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalln("error while reading", err)
			}
			inputChannel <- input

		}()

		select {
		case <-timer.C:
			fmt.Println("\n\n*********** Times is Up! Better luck next time***********")
			break quizLoop
		case input := <-inputChannel:
			var ans int
			input = strings.TrimSpace(input)
			if input != "" {
				ans, _ = strconv.Atoi(input)

				correct_answer, _ := strconv.Atoi(record[2])
				if ans != correct_answer {
					answers = append(answers, Answer{record[0], ans, 0, "Wrong"})
				} else {
					answers = append(answers, Answer{record[0], ans, 10, "Correct"})
				}
			}
		}

	}

	score := 0
	fmt.Println("Your Anwers\n===========")
	if len(answers) == 0 {
		fmt.Println("No answers provided")
	} else {
		for _, ans := range answers {
			score += ans.Score
			fmt.Printf("\n%s: Your answer: %d, Score: %d, Result: %s", ans.No, ans.Ans, ans.Score, ans.Result)
		}
	}

	var result string

	if score < 30 {
		result = "Failed"
	} else {
		result = "Passed"
	}

	fmt.Printf("\n================\nYour Score = %d\nResult: %s\n================", score, result)
}
