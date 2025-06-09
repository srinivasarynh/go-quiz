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
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'questions,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit of the quiz is 30 seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("failed to open csv file: %s\n", *csvFilename))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("failed to parse the provided file")
	}
	
	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

  correct := 0 
	problemLoop: 
		for i, p := range problems {
			fmt.Printf("problem %d: %s = \n", i+1, p.q)
			answerCh := make(chan string)
			go func() {
				var answer string
		  	fmt.Scanf("%s\n", &answer)
				answerCh <- answer
			}()
			select {
	   		case <-timer.C:
				fmt.Println()
	  	  break problemLoop
		  case answer := <-answerCh: 
		    if answer == p.a {
		    	correct++
	  	  }
			}
	}

	fmt.Printf("you scored %d out of %d,\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}


func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
