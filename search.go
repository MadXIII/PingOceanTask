package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type Store struct {
	Map map[string]int
	sync.RWMutex
}

func NewStore() *Store {
	return &Store{Map: make(map[string]int)}
}

func (s *Store) setMap(url string, count int) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.Map[url]; ok {
		return errors.New("key is already exists")
	}
	s.Map[url] = count

	return nil
}

func (s *Store) getMap() map[string]int {
	s.RLock()
	defer s.RUnlock()

	return s.Map
}

func goSender(ctx context.Context, urls []string, str string) (map[string]int, error) {
	m := NewStore()

	wg := sync.WaitGroup{}

	for _, url := range urls {
		select {
		case <-ctx.Done():
			fmt.Println("Test")
			return nil, errors.New("timeout limit")
		default:
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				count, err := searchAndCount(ctx, url, str)
				if err != nil {
					fmt.Println(err)
					return
				}
				if err := m.setMap(url, count); err != nil {
					fmt.Println(err)
					return
				}
			}(url)
		}
	}
	wg.Wait()

	return m.getMap(), nil
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
