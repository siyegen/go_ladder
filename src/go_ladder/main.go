package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Hi Words")

	chanMap := make(map[int]chan string)
	metaChan := make(chan chan string)

	go func() {
		for {
			dd := <-metaChan
			go func(wordSizedChan chan string) {
				count := 0
				for word := range wordSizedChan {
					fmt.Println(len(word))
					count++
				}
				fmt.Printf("Size %d words: %d", 2, count)
			}(dd)
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
			metaChan <- sizedChan
			fmt.Println("New Chan", len(word))
		}
		// Add wg here, remove it on processing side?
		sizedChan <- word
	}
	close(metaChan)
	// We could range over all channels and close them?
}

func handleWord() {

}
