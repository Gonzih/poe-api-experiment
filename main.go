package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var exitingTheLoopErr = fmt.Errorf("Exiting the loop, sorry!")

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	command := os.Args[1]
	switch command {
	case "pull":
		checkErr(pull("data/responses.bin"))

	case "list":
		checkErr(walkResponses("data/responses.bin", func(r *Response) error {
			fmt.Printf("%s\n", r.GetNextChangeId())
			return nil
		}))
	case "last-id":
		log.Printf("Last ID: %s", lastNextChangeID("data/responses.bin"))
	case "generate-fields":
		checkErr(generateFields("data/responses.bin"))
	case "generate-input":
		input, err := generateMLInput("data/responses.bin", 10000)
		if err != nil {
			if err == exitingTheLoopErr {
				log.Printf(`Ignoring error "%s"`, err)
			} else {
				log.Fatal(err)
			}
		}

		err = input.Save()
		checkErr(err)
	case "ml-main":
		input := &MLInput{}
		err := input.Load()
		checkErr(err)
		evalFn, err := linearRegression(input)
		checkErr(err)

		c := 0
		limit := 10
		for {
			c++
			if c > limit {
				break
			}

			sample := input.Fields[rand.Intn(len(input.Fields))]
			sampleResult, err := evalFn(sample)
			checkErr(err)

			log.Printf("For input %v y = %3.3f", sample[0], sampleResult)
		}
	default:
		log.Fatalf("Uknown command %s", command)
	}
}
