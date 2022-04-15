package main

import (
	"context"
	"fmt"
	"io/ioutil"
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

func (s *Store) setMap(url string, count int) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.Map[url]; ok {
		return fmt.Errorf("%s url is already exists", url)
	}
	s.Map[url] = count

	return nil
}

func (s *Store) getMap() map[string]int {
	s.RLock()
	defer s.RUnlock()

	return s.Map
}

func worker(ctx context.Context, urls []string, str string) map[string]int {
	store := NewStore()

	workerSize := 20
	mainCh := make(chan string, workerSize)
	errCh := make(chan error, len(urls))

	wg := sync.WaitGroup{}
	for i := 0; i < workerSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range mainCh {
				store.searchAndCount(ctx, url, str, errCh)
			}
		}()
	}

	for _, url := range urls {
		mainCh <- url
	}

	close(mainCh)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case err := <-errCh:
				fmt.Println(err)
			case <-time.After(10 * time.Second):
				return
			}
		}
	}()

	defer close(errCh)

	wg.Wait()

	return store.getMap()
}

func (s *Store) searchAndCount(ctxParent context.Context, url, str string, errCh chan<- error) {
	ctx, cancelFunc := context.WithTimeout(ctxParent, 10*time.Second)
	defer cancelFunc()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		errCh <- fmt.Errorf("error due URL: %s, searchAndCount, NewRequest: %v", url, err)
		return
	}

	client := http.Client{}

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		errCh <- fmt.Errorf("error due URL: %s, searchAndCount, clientDo: %v", url, err)
		return
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errCh <- fmt.Errorf("error due URL: %s, searchAndCount, ReadAll: %v", url, err)
		return
	}

	count := strings.Count(string(bytes), str)
	if err := s.setMap(url, count); err != nil {
		errCh <- fmt.Errorf("error due URL: %s, searchAndCount, setMap: %v", url, err)
		return
	}
}

// ./pingoceantask -url https://github.com/ -url https://pingocean.com/ -url https://pingocean.com/a -url https://pingocean.com/b -url https://pingocean.com/c -url https://pingocean.com/d -url https://pingocean.com/e -url https://pingocean.com/f -url https://pingocean.com/g -url https://pingocean.com/h -url https://pingocean.com/i -url https://pingocean.com/j -url https://pingocean.com/k -url https://pingocean.com/l -url https://pingocean.com/m -url https://pingocean.com/n -url https://pingocean.com/o -url https://pingocean.com/p -url https://pingocean.com/q -url https://pingocean.com/r -url https://pingocean.com/s -url https://pingocean.com/t -url https://pingocean.com/u -url https://pingocean.com/v -url https://pingocean.com/w -url https://pingocean.com/x -url https://pingocean.com/y -url https://pingocean.com/z -str pingocean
