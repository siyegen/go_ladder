package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type WordResult struct {
	size     int
	wordChan chan string
}

func main() {
	fmt.Println("Hi Words")

	chanMap := make(map[int]chan string)
	metaChan := make(chan WordResult)

	finishedChan := handleSizedChan(metaChan)

	file, err := os.Open("./words")
	if err != nil {
		log.Fatal("Can't open file", err)
	}

	scanner := bufio.NewScanner(file)
	// Read line by line, send each word on its sized channel
	for scanner.Scan() {
		var sizedChan chan string
		word := scanner.Text()
		// Try to get channel for n size from map, otherwise create
		sizedChan, ok := chanMap[len(word)]
		if !ok {
			sizedChan = make(chan string)
			chanMap[len(word)] = sizedChan
			metaChan <- WordResult{len(word), sizedChan}
			fmt.Println("New Chan", len(word))
		}
		sizedChan <- word
	}
	// Wait for shit to finish, maybe select with kill later
	for _, ch := range chanMap {
		close(ch)
	}
	fmt.Println("Waiting for shit to finish")
	<-finishedChan
}

func handleWord(size int, wordSizedChan chan string) int {
	count := 0
	for _ = range wordSizedChan {
		count++
	}
	if size == 10 {
		// Want to make sure even though the channel is closed, I can still do shit here
		time.Sleep(3 * time.Second)
		fmt.Println("Just proving a point")
	}
	fmt.Printf("Size %d words: %d\n", size, count)
	return count
}

func handleSizedChan(metaChan chan WordResult) chan struct{} {
	finished := make(chan struct{})
	var wg sync.WaitGroup
	total := 0
	go func() {
		for {
			dd := <-metaChan
			wg.Add(1)
			go func() {
				count := handleWord(dd.size, dd.wordChan)
				total += count
				wg.Done()
			}()
		}
	}()

	go func() {
		wg.Wait()
		fmt.Println("time to close chan")
		fmt.Println("Total!", total)
		close(finished)
	}()

	return finished
}
