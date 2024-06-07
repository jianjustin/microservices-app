package x

import (
	"io"
	"log"
	"net/http"
	"testing"
)

func TestForGet(t *testing.T) {
	resp, err := http.Get("https://baidu.com")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}
