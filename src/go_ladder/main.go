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

	var fin sync.WaitGroup

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
		// Maybe get all sizedChans first, and then wg.Add(len(map))
		if !ok {
			sizedChan = make(chan string)
			chanMap[len(word)] = sizedChan
			metaChan <- WordResult{len(word), sizedChan}
			// fmt.Println("New Chan", len(word))
		}
		sizedChan <- word
	}

	totalCount := 0
	fin.Add(len(chanMap))
	for _, ch := range chanMap {
		close(ch)
	}
	go func() {
		for {
			num := <-finishedChan
			totalCount += num
			fin.Done()
		}
	}()

	fmt.Println("Waiting for shit to finish")
	fin.Wait()
	fmt.Println("Total!", totalCount)
}

func handleWord(size int, wordSizedChan chan string) int {
	count := 0
	if size == 10 {
		// Want to make sure even though the channel is closed, I can still do shit here
		time.Sleep(2 * time.Millisecond)
		fmt.Println("Just proving a point")
	}
	for _ = range wordSizedChan {
		count++
	}
	// fmt.Printf("Size %d words: %d\n", size, count)
	return count
}

func handleSizedChan(metaChan chan WordResult) chan int {
	finished := make(chan int)

	go func() {
		for {
			dd := <-metaChan
			go func() {
				count := handleWord(dd.size, dd.wordChan)
				finished <- count
			}()
		}
	}()

	return finished
}
