package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"unicode/utf8"
)

type WordResult struct {
	size     int
	wordChan chan string
}

type WordNode struct {
	value     string
	connected []*WordNode
}

type WordGraph struct {
	nodes map[string]*WordNode
}

func (w *WordGraph) add(node *WordNode) {
	for _, currNode := range w.nodes {
		diff := 0
		for i, width := 0, 0; i < len(currNode.value); i += width {
			runeVal, size := utf8.DecodeRuneInString(currNode.value[i:])
			otherRuneVal, _ := utf8.DecodeRuneInString(node.value[i:])
			if runeVal != otherRuneVal {
				diff++
			}
			if diff > 1 {
				break
			}
			width = size
		}
		currNode.connected = append(currNode.connected, node)
	}
}

func main() {

	fmt.Println("Hi Words")

	chanMap := make(map[int]chan string)
	metaChan := make(chan WordResult)

	finishedChan := handleSizedChan(metaChan)

	var fin sync.WaitGroup

	file, err := os.Open("./small_words")
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

func handleWord(size int, wordSizedChan chan string) (int, *WordGraph) {
	count := 0
	wordGraph := &WordGraph{make(map[string]*WordNode)}
	for currWord := range wordSizedChan {
		wordNode := &WordNode{
			value:     currWord,
			connected: make([]*WordNode, 0),
		}
		wordGraph.nodes[currWord] = wordNode
		count++
	}

	return count, wordGraph
}

func handleSizedChan(metaChan chan WordResult) chan int {
	finished := make(chan int)

	go func() {
		for {
			dd := <-metaChan
			go func() {
				count, gg := handleWord(dd.size, dd.wordChan)
				fmt.Println(gg)
				finished <- count
			}()
		}
	}()

	return finished
}
