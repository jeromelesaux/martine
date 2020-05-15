package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	count := 0
	for scanner.Scan() {
		if count%8 == 0 {
			fmt.Printf("\ndb ")
		}
		fmt.Printf("%s", scanner.Text())
		if (count+1)%8 != 0 {
			fmt.Printf(", ")
		}
		count++
	}
}
