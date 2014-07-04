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
			// On creation we need to setup the chan to get data
			metaChan <- WordResult{len(word), sizedChan}
			fmt.Println("New Chan", len(word))
		}
		// Send the word over, our work is finished
		sizedChan <- word
	}
	// Wait for shit to finish, maybe select with kill later
	time.Sleep(3 * time.Second)
	for _, ch := range chanMap {
		close(ch)
	}
	fmt.Println("Waiting for shit to finish")
	<-finishedChan
}

func handleWord(size int, wordSizedChan chan string) {
	count := 0
	for word := range wordSizedChan {
		fmt.Println("handleWord", word)
		count++
	}
	fmt.Printf("Size %d words: %d\n", size, count)
}

func handleSizedChan(metaChan chan WordResult) chan struct{} {
	finished := make(chan struct{})
	var wg sync.WaitGroup
	go func() {
		for {
			dd := <-metaChan
			fmt.Println("adding wg")
			wg.Add(1)
			go func() {
				handleWord(dd.size, dd.wordChan)
				fmt.Println("wg.done")
				wg.Done()
			}()
		}
	}()

	go func() {
		wg.Wait()
		fmt.Println("time to close chan")
		close(finished)
	}()

	return finished
}
