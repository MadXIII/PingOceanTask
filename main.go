package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Store struct {
	Map map[string]int
	// sync.RWMutex
}

type Slice []string

func (s *Slice) String() string {
	return "my string"
}

func (s *Slice) Set(val string) error {
	*s = append(*s, val)
	return nil
}

func main() {
	var urls Slice
	flag.Var(&urls, "url", "flag can be -url || --url to send urls")

	str := flag.String("str", "pingocean", "flag can be -str || --str to send string")
	flag.Parse()

	if len(urls) < 1 {
		urls = []string{"https://pingocean.com/"}
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelFunc()

	err := test(ctx, urls, *str)
	if err != nil {
		log.Println(err)
	}
}

func test(ctx context.Context, urls []string, str string) error {
	m := &Store{
		Map: make(map[string]int),
	}

	for _, url := range urls {
		count, err := searchAndCount(ctx, url, str)
		if err != nil {
			return fmt.Errorf("... %w", err)
		}
		m.Map[url] = count

	}
	fmt.Println(m.Map)
	return nil
}

func searchAndCount(ctx context.Context, url, str string) (int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("error searchAndCount, NewRequest: %w", err)
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error searchAndCount, clientDo: %w", err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error searchAndCount, ReadAll: %w", err)
	}

	return strings.Count(string(bytes), str), nil
}
