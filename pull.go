package main

import (
	fmt "fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gogo/protobuf/jsonpb"
)

func httpGetBody(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return resp.Body, err
	}

	return resp.Body, err
}

func pull(dbPath string) error {
	nextChangeID := lastNextChangeID(dbPath)

	log.Printf("Restarting from ID %s", nextChangeID)

	counter := 0

	for {
		url := fmt.Sprintf("http://www.pathofexile.com/api/public-stash-tabs?id=%s", nextChangeID)

		data := &Response{}

		body, err := httpGetBody(url)
		if err != nil {
			return err
		}

		defer body.Close()
		unmarshaller := jsonpb.Unmarshaler{AllowUnknownFields: true}
		err = unmarshaller.Unmarshal(body, data)
		if err != nil {
			return err
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

	return nil
}
