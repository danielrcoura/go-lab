package main

import "fmt"

// Master Go Programming With These Concurrency Patterns (in 40 minutes)
// https://www.youtube.com/watch?v=qyM8Pi1KiiM

func main() {
	numbers := []int{1, 2, 3, 4, 5}
	numberChannel := sliceToChannel(numbers)
	squaredChannel := square(numberChannel)
	printChannel(squaredChannel)
}

func sliceToChannel(numbers []int) chan int {
	out := make(chan int)
	go func() {
		for _, n := range numbers {
			fmt.Println("sliceToChannel...", n)
			out <- n
		}
		close(out)
	}()
	return out
}

func square(numbers chan int) chan int {
	out := make(chan int)
	go func() {
		for n := range numbers {
			fmt.Println("square...", n)
			out <- n * n
		}
		close(out)
	}()
	return out
}

func printChannel(numbers chan int) {
	for n := range numbers {
		fmt.Println("printChannel:", n)
	}
}
