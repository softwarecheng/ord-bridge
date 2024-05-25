package indexer

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"
	"time"

	"github.com/softwarecheng/ord-bridge/indexer/pb"
	"google.golang.org/protobuf/proto"
)

func TestDecode2(t *testing.T) {
	value1 := pb.Inscription{Number: 1, Sat: 2, GenesesAddress: "123467890"}

	//fmt.Printf("%v\n", value1)
	start := time.Now()
	var encodeBytes []byte
	var err error
	for i := 0; i < 1000; i++ {
		encodeBytes, err = proto.Marshal(&value1)
		if err != nil {
			t.Fatal(err)
		}
	}
	fmt.Printf("encode time: %vs\n", time.Since(start).Seconds()) // 0.5ms
	fmt.Printf("%d\n", len(encodeBytes))                          // 15

	start = time.Now()
	result1 := &pb.Inscription{}
	for i := 0; i < 1000; i++ {
		err = proto.Unmarshal(encodeBytes, result1)
		if err != nil {
			t.Fatal(err)
		}
	}
	fmt.Printf("decode time: %vs\n", time.Since(start).Seconds()) // 0.15ms
	//fmt.Printf("%v\n", result1)

	start = time.Now()
	for i := 0; i < 1000; i++ {
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(&value1); err != nil {
			t.Fatal(err)
		}
		encodeBytes = buf.Bytes()
	}
	fmt.Printf("encode time: %vs\n", time.Since(start).Seconds()) // 10ms
	fmt.Printf("%d\n", len(encodeBytes))                          // 93

	start = time.Now()
	result2 := &pb.Inscription{}
	for i := 0; i < 1000; i++ {
		buf := bytes.NewBuffer(encodeBytes)
		dec := gob.NewDecoder(buf)
		err := dec.Decode(result2)
		if err != nil {
			t.Fatal(err)
		}
	}
	fmt.Printf("decode time: %vs\n", time.Since(start).Seconds()) // 40ms
}
