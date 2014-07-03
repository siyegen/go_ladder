package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

type WordResult struct {
	size     int
	wordChan chan string
}

func main() {
	fmt.Println("Hi Words")

	chanMap := make(map[int]chan string)
	metaChan := make(chan WordResult)
	doneChan := make(chan struct{})
	var wg sync.WaitGroup

	go func() {
		for {
			dd := <-metaChan
			go handleWord(dd.size, doneChan, dd.wordChan, &wg)
		}
	}()

	file, err := os.Open("./words")
	if err != nil {
		log.Fatal("Can't open file", err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var sizedChan chan string
		word := scanner.Text()
		sizedChan, ok := chanMap[len(word)]
		if !ok {
			sizedChan = make(chan string)
			chanMap[len(word)] = sizedChan
			metaChan <- WordResult{len(word), sizedChan}
			fmt.Println("New Chan", len(word))
			wg.Add(1)
		}
		// Add wg here, remove it on processing side?
		sizedChan <- word
	}
	close(doneChan)
	for size, ch := range chanMap {
		fmt.Println("Closing channel for", size)
		close(ch)
	}
	// wg.Wait()
	// close(metaChan)
	// We could range over all channels and close them?
}

func handleWord(size int, done <-chan struct{}, wordSizedChan chan string, wg *sync.WaitGroup) {
	count := 0
LOOP:
	for {
		select {
		case <-done:
			close(wordSizedChan)
			log.Println("Closing wordSizedChan")
			break LOOP
		case word := <-wordSizedChan:
			fmt.Println(word, size)
			count++
		}
	}
	fmt.Printf("Size %d words: %d\n", size, count)
	// wg.Done()
}
