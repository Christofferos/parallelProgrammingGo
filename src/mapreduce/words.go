package main

import (
	"fmt"
	"time"
	"io/ioutil"
	"strings"
	"sync"
)

const DataFile = "loremipsum.txt"

// Return the word frequencies of the text argument.
func WordCount(text string) map[string]int {
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, ".", "")
	text = strings.ReplaceAll(text, ",", "")
  words := strings.Fields(text)

	freqs := make(map[string]int)
	wg := new(sync.WaitGroup)
	wg2 := new(sync.WaitGroup)
  mapchan := make (chan map[string]int)

  numwords := len(words)
  processes := 10
	processlength := numwords/processes

	for startpoint := 0; startpoint < numwords; startpoint+=processlength {
		j := startpoint+processlength
		wg.Add(1)
		go func(startpoint, j int, wg *sync.WaitGroup, words []string) {
			innerfreqs := make(map[string]int)
			for i := startpoint; i < j && i<numwords; i++ {
					innerfreqs[words[i]]++
			}
			mapchan <- innerfreqs
			wg.Done()
		}(startpoint,j, wg, words)
	}

	//REDUCE
	go func(mapchan chan map[string]int, freqs map[string]int){
    wg2.Add(1)
    for i := 0; i < processes; i++ {
      //partialFreq := <- mapchan
			for a, b := range <- mapchan {
				freqs[a] += b
			}
		}
    wg2.Done()
	}(mapchan, freqs)

  wg.Wait()
	close(mapchan)
	wg2.Wait()
	return freqs
}

// Benchmark how long it takes to count word frequencies in text numRuns times.
//
// Return the total time elapsed.
func benchmark(text string, numRuns int) int64 {
	start := time.Now()
	for i := 0; i < numRuns; i++ {
		WordCount(text)
	}
	runtimeMillis := time.Since(start).Nanoseconds() / 1e6

	return runtimeMillis
}

// Print the results of a benchmark
func printResults(runtimeMillis int64, numRuns int) {
	fmt.Printf("amount of runs: %d\n", numRuns)
	fmt.Printf("total time: %d ms\n", runtimeMillis)
	average := float64(runtimeMillis) / float64(numRuns)
	fmt.Printf("average time/run: %.2f ms\n", average)
}

func main() {
	// read in DataFile as a string called data
	b, err := ioutil.ReadFile(DataFile)
    if err != nil {
        fmt.Print(err)
    }
	data := string(b) // convert file content to a string

	fmt.Printf("%#v", WordCount(string(data)))

	numRuns := 100
	runtimeMillis := benchmark(string(data), numRuns)
	printResults(runtimeMillis, numRuns)
}
