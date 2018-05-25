package main

import (
	binary "encoding/binary"
	fmt "fmt"
	"io"
	"os"

	proto "github.com/gogo/protobuf/proto"
)

const sizeOfLength = 8

type length int64

var endianness = binary.LittleEndian

func lastNextChangeID(dbPath string) string {
	var id string

	walkResponses(dbPath, func(r *Response) error {
		id = r.GetNextChangeId()
		return nil
	})

	return id
}

func walkResponses(dbPath string, walkFn func(*Response) error) error {
	f, err := os.Open(dbPath)
	if err != nil {
		return fmt.Errorf("could not read %s: %v", dbPath, err)
	}

	for {
		var l length
		err := binary.Read(f, endianness, &l)

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return fmt.Errorf("could not decode message length: %v", err)
		}

		data := make([]byte, l)
		_, err = f.Read(data)

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		response := &Response{}
		if err := proto.Unmarshal(data, response); err != nil {
			return fmt.Errorf("could not read task: %v", err)
		}

		if err := walkFn(response); err != nil {
			return err
		}
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
