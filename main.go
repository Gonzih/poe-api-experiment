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

func abs(i float64) float64 {
	return float64(math.Abs(float64(i)))
}

func must(err error) {
	if err != nil {
		if err == exitingTheLoopErr {
			log.Printf(`Ignoring error "%s"`, err)
		} else {
			log.Fatal(err)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	command := os.Args[1]
	switch command {
	case "pull":
		must(pull("data/responses.bin"))

	case "list":
		must(walkResponses("data/responses.bin", func(r *Response) error {
			fmt.Printf("%s\n", r.GetNextChangeId())
			return nil
		}))
	case "last-id":
		log.Printf("Last ID: %s", lastNextChangeID("data/responses.bin"))
	case "generate-fields":
		must(generateFields("data/responses.bin"))
	case "generate-input":
		input, err := generateMLInput("data/responses.bin", 1000)
		must(err)

		err = input.Save()
		must(err)
	case "generate-csv":

		f, err := os.OpenFile("data/data.csv", os.O_RDWR|os.O_CREATE, 0755)
		defer f.Close()
		must(err)

		w := csv.NewWriter(f)

		fields, err := loadFieldsConfiguration()

		must(err)

		allFields := []string{
			"price",
			"item level",
			"num of sockets",
			"num of linked sockets",
			"corrupted",
			"is accessory",
			"is armour",
			"is jweel",
			"is weapon",
			"is gem",
			"is flask",
			"is map",
			"is currency",
			"is card",
		}

		allFields = append(allFields, fields.Properties...)
		allFields = append(allFields, fields.ImplicitMods...)
		allFields = append(allFields, fields.ExplicitMods...)

		err = w.Write(allFields)
		must(err)

		var stringRow []string = make([]string, len(allFields))

		err = walkFeatures("data/responses.bin", 100000000,
			func(features []float64) error {
				for i, f := range features {
					stringRow[i] = strconv.FormatFloat(f, 'f', 6, 64)
				}

				err = w.Write(stringRow)
				return err
			})
		must(err)

	case "ml-main":
		input := &MLInput{}
		err := input.Load()
		must(err)
		evalFn, err := linearRegression(input)
		must(err)
		// inSize := len(input.Fields)

		log.Printf("Calculating accuracy")
		for i := 0; i < 10; i++ {
			sample := input.Fields[i]
			originPrice := sample[0]
			sampleResult, err := evalFn(sample)
			must(err)
			errorRating := abs((sampleResult - originPrice) / originPrice * 100)
			log.Printf("%5.0f -> %5.0f, error %5.0f%%", originPrice, sampleResult, errorRating)
		}

	default:
		log.Fatalf("Uknown command %s", command)
	}
}
