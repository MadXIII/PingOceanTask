package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"
)

type Slice []string

func (s *Slice) String() string {
	return "my string"
}

func (s *Slice) Set(val string) error {
	*s = append(*s, val)
	return nil
}

func main() {
	start := time.Now()
	var urls Slice
	flag.Var(&urls, "url", "flag can be -url || --url to send urls")
	str := flag.String("str", "", "flag can be -str || --str to send string")
	flag.Parse()

	if len(urls) < 1 {
		log.Printf("No URLs")
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancelFunc()

	m := worker(ctx, urls, *str)

	fmt.Println(m)
	for key, val := range m {
		fmt.Printf("%s %d\n", key, val)
	}
	fmt.Println(time.Since(start))
}
