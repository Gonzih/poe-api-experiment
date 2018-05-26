package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

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
		checkErr(pull("data.bin"))

	case "list":
		checkErr(walkResponses("data.bin", func(r *Response) error {
			fmt.Printf("%s\n", r.GetNextChangeId())
			return nil
		}))
	case "last-id":
		log.Printf("Last ID: %s", lastNextChangeID("data.bin"))
	case "generate-input":
		_, err := generateMLInput("data.bin")
		checkErr(err)
	case "generate-fields":
		checkErr(generateFields("data.bin"))
	case "ml-main":
		input, err := generateMLInput("data.bin")
		checkErr(err)
		evalFn, err := linearRegression(input)
		checkErr(err)

		c := 0
		for {
			c++
			if c > 5 {
				break
			}

			sample := input[rand.Intn(len(input))]
			sampleResult, err := evalFn(sample)
			checkErr(err)

			log.Printf("For input %v y = %3.3f", sample, sampleResult)
		}
	default:
		log.Fatalf("Uknown command %s", command)
	}
}
