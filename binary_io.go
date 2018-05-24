package main

import (
	"bytes"
	binary "encoding/binary"
	fmt "fmt"
	"io/ioutil"
	"os"

	proto "github.com/gogo/protobuf/proto"
)

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
