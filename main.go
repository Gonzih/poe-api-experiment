package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var exitingTheLoopErr = fmt.Errorf("Exiting the loop, sorry!")

func abs(i float32) float32 {
	return float32(math.Abs(float64(i)))
}

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
		input, err := generateMLInput("data/responses.bin", 1000)
		if err != nil {
			if err == exitingTheLoopErr {
				log.Printf(`Ignoring error "%s"`, err)
			} else {
				log.Fatal(err)
			}
		}

		err = input.Save()
		checkErr(err)
	case "generate-csv":
		input, err := generateMLInput("data/responses.bin", 1000)
		if err != nil {
			if err == exitingTheLoopErr {
				log.Printf(`Ignoring error "%s"`, err)
			} else {
				log.Fatal(err)
			}
		}

		w := csv.NewWriter(os.Stdout)

		var stringRow []string = make([]string, len(input.Fields[0]))

		for _, row := range input.Fields {
			for i, f := range row {
				stringRow[i] = strconv.FormatFloat(f, 'f', 6, 64)
			}

			if err := w.Write(stringRow); err != nil {
				log.Fatal(err)
			}
		}

		checkErr(err)
	case "ml-main":
		input := &MLInput{}
		err := input.Load()
		checkErr(err)
		evalFn, err := linearRegression(input)
		checkErr(err)
		// inSize := len(input.Fields)

		log.Printf("Calculating accuracy")
		for i := 0; i < 10; i++ {
			sample := input.Fields[i]
			originPrice := sample[0]
			sampleResult, err := evalFn(sample)
			checkErr(err)
			errorRating := abs((sampleResult - originPrice) / originPrice * 100)
			log.Printf("%5.0f -> %5.0f, error %5.0f%%", originPrice, sampleResult, errorRating)
		}

	default:
		log.Fatalf("Uknown command %s", command)
	}
}
