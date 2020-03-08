package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	low := 1
	high := 100
	fmt.Println("Please think of a number between", low, "and", high)
	fmt.Println("Press ENTER when ready")
	scanner.Scan()
	count := 0
	for {

		guess := (low + high) / 2
		fmt.Println("I guess the number is", guess)
		fmt.Println("Is that:")
		fmt.Println("(a) too high?")
		fmt.Println("(b) too low?")
		fmt.Println("(c) correct?")
		scanner.Scan()
		response := scanner.Text()
		if response == "a" {
			high = guess - 1
			count++
		} else if response == "b" {
			low = guess + 1
			count++
		} else if response == "c" {
			fmt.Println("I won!")
			fmt.Println("you have made", count, "attempts to win")
			break
		} else {
			fmt.Println("Invalid response, try again")
			count++
		}
	}
}
