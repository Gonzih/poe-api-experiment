package main

import (
	"fmt"
	"io"
	"log"
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

		err = appendToFile(data, dbPath)

		log.Printf("Finished fetching '%s', counter is %d", url, counter)

		counter++
		nextChangeID = data.GetNextChangeId()

		if len(data.Stashes) == 0 {
			break
		}

		time.Sleep(time.Second * 2)
	}
}

func main() {
	command := os.Args[1]
	switch command {
	case "pull":
		pull("data.bin")

	case "list":
		list("data.bin")
	default:
		log.Fatalf("Uknown command %s", command)
	}
}
