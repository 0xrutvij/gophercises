package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	csv_filename := flag.String("csv", "problems.csv", "A csv file of the format 'question,answer'")
	time_limit := flag.Int("limit", 20, "Time limit for the quiz in seconds.")

	flag.Parse()

	file, err := os.Open(*csv_filename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open CSV file %s", *csv_filename))
	}

	quiz_input := csv.NewReader(file)

	lines, err := quiz_input.ReadAll()

	if err != nil {
		exit("Failed to parse the CSV file provided.")
	}

	problems := parse_lines(lines)

	var cor_ans int
	timer := time.NewTimer(time.Duration(*time_limit) * time.Second)

problem_loop:
	for i, problem := range problems {

		ans := make(chan string)
		go prompt_question(i, problem.ques, ans)

		select {
		case <-timer.C:
			fmt.Println()
			break problem_loop
		case answer := <-ans:
			if answer == problem.ans {
				cor_ans += 1
			}
		}

	}

	fmt.Printf("\nYour score: %d/%d \n", cor_ans, len(problems))

}

func prompt_question(qno int, ques string, ans chan string) {

	var input string

	fmt.Printf("Problem #%d: %s = ", qno+1, ques)
	_, err := fmt.Scanf("%s", &input)

	if err != nil {
		ans <- ""
	}

	ans <- input
}

func parse_lines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			ques: line[0],
			ans:  strings.TrimSpace(line[1]),
		}
	}

	return ret
}

type problem struct {
	ques string
	ans  string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
