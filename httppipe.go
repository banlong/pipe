package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type msg struct {
	Text string
}

func handleErr(err error) {
	if err != nil {
		log.Fatalf("%s\n", err)
	}
}

// use a io.Pipe to connect a JSON encoder to an HTTP POST: this way you do
// not need a temporary buffer to store the JSON bytes
func main() {
	pReader, pWriter := io.Pipe()

	// writing without a reader will deadlock so write in a goroutine
	go func() {
		// it is important to close the writer or reading from the other end of the
		// pipe will never finish
		defer pWriter.Close()

		m := msg{Text: "brought to you by io.Pipe()"}
		err := json.NewEncoder(pWriter).Encode(&m)
		handleErr(err)
		//the above codes write an encoded message into a pipe, which is connected to the ENCODER
	}()

	resp, err := http.Post("http://localhost:9092/pipe", "application/json", pReader)
	//resp, err := http.Post("https://httpbin.org/post", "application/json", pReader)
	handleErr(err)
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	handleErr(err)

	log.Printf("%s\n", b)
}
