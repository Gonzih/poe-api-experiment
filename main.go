package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gogo/protobuf/jsonpb"
	proto "github.com/gogo/protobuf/proto"
)

func httpGetBody(url string) io.ReadCloser {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	return resp.Body
}

const sizeOfLength = 8

type length int64

var endianness = binary.LittleEndian

func list(dbPath string) error {
	b, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return fmt.Errorf("could not read %s: %v", dbPath, err)
	}

	for {
		if len(b) == 0 {
			return nil
		} else if len(b) < sizeOfLength {
			return fmt.Errorf("remaining odd %d bytes, what to do?", len(b))
		}

		var l length
		if err := binary.Read(bytes.NewReader(b[:sizeOfLength]), endianness, &l); err != nil {
			return fmt.Errorf("could not decode message length: %v", err)
		}
		b = b[sizeOfLength:]

		data := &Response{}
		if err := proto.Unmarshal(b[:l], data); err != nil {
			return fmt.Errorf("could not read task: %v", err)
		}
		b = b[l:]

		fmt.Printf(" %s\n", data.GetNextChangeId())
	}
}

func appendToFile(data *Response, dbPath string) error {
	b, err := proto.Marshal(data)

	if err != nil {
		return fmt.Errorf("could not encode data: %v", err)
	}

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("could not open %s: %v", dbPath, err)
	}

	if err := binary.Write(f, endianness, length(len(b))); err != nil {
		return fmt.Errorf("could not encode length of message: %v", err)
	}
	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("could not write task to file: %v", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("could not close file %s: %v", dbPath, err)
	}
	return nil
}

func pull(dbPath string) {
	nextChangeID := ""
	counter := 0

	for {
		url := fmt.Sprintf("http://www.pathofexile.com/api/public-stash-tabs?id=%s", nextChangeID)

		data := &Response{}

		body := httpGetBody(url)
		defer body.Close()
		unmarshaller := jsonpb.Unmarshaler{AllowUnknownFields: false}
		err := unmarshaller.Unmarshal(body, data)

		if err != nil {
			panic(err)
		}

		err = appendToFile(data, dbPath)

		log.Printf("Finished fetching '%s', counter is %d", url, counter)

		counter++
		nextChangeID = data.GetNextChangeId()

		if len(nextChangeID) == 0 {
			break
		}
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
		log.Fatal("Uknown command %s", command)
	}
}
