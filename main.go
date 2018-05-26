package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gogo/protobuf/jsonpb"
)

func httpGetBody(url string) io.ReadCloser {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	return resp.Body
}

func pull(dbPath string) {
	nextChangeID := lastNextChangeID(dbPath)

	log.Printf("Restarting from ID %s", nextChangeID)

	counter := 0

	for {
		url := fmt.Sprintf("http://www.pathofexile.com/api/public-stash-tabs?id=%s", nextChangeID)

		data := &Response{}

		body := httpGetBody(url)
		defer body.Close()
		unmarshaller := jsonpb.Unmarshaler{AllowUnknownFields: true}
		err := unmarshaller.Unmarshal(body, data)

		if err != nil {
			panic(err)
		}

		if len(data.Stashes) == 0 {
			log.Println("Got empty stashes in response")
			break
		}

		err = appendToFile(data, dbPath)

		log.Printf("Finished fetching '%s', counter is %d", url, counter)

		counter++
		nextChangeID = data.GetNextChangeId()

		time.Sleep(time.Millisecond * 500)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	command := os.Args[1]
	switch command {
	case "pull":
		pull("data.bin")

	case "list":
		walkResponses("data.bin", func(r *Response) error {
			fmt.Printf("%s\n", r.GetNextChangeId())
			return nil
		})
	case "last-id":
		log.Printf("Last ID: %s", lastNextChangeID("data.bin"))
	case "generate-input":
		generateMLInput("data.bin")
	case "ml-main":
		input := generateMLInput("data.bin")
		evalFn, err := linearRegression(input)
		if err != nil {
			log.Fatal(err)
		}

		sample := input[rand.Intn(len(input))]
		sampleResult, err := evalFn(sample)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("For input %v y = %3.3f", sample, sampleResult)
	default:
		log.Fatalf("Uknown command %s", command)
	}
}
