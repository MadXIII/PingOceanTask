package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Store struct {
	Map map[string]int
	sync.RWMutex
}

func NewStore() *Store {
	return &Store{Map: make(map[string]int)}
}

func (s *Store) Set(url string, count int) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.Map[url]; ok {
		return errors.New("key is already exists")
	}
	s.Map[url] = count

	return nil
}

func (s *Store) Get() map[string]int {
	s.RLock()
	defer s.RUnlock()

	return s.Map
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

	err := goSender(ctx, urls, *str)
	if err != nil {
		log.Println(err)
	}
}

func goSender(ctx context.Context, urls []string, str string) error {
	m := NewStore()

	wg := sync.WaitGroup{}

	for _, url := range urls {
		select {
		case <-ctx.Done():
			return errors.New("timeout limit")
		default:
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				count, err := searchAndCount(ctx, url, str)
				if err != nil {
					fmt.Println(err)
					return
				}
				if err := m.Set(url, count); err != nil {
					fmt.Println(err)
					return
				}
			}(url)
		}
	}
	wg.Wait()

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

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error searchAndCount, ReadAll: %w", err)
	}

	return strings.Count(string(bytes), str), nil
}
