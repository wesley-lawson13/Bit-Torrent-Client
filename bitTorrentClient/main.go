package main

import (
	"bufio"
	"strconv"
	"fmt"
	"os"
)

func printCoins(scanner *bufio.Scanner) {
    total := 0
    count := 0
    var max int
    var min int
    for {
        fmt.Print("> ")
        if scanner.Scan(); scanner.Err() != nil {
            fmt.Print("Scanner failed")
            return
        }
        input := scanner.Text()
        fmt.Printf("\n inputted '%s'\n", input)
        if input == "$" {
            break
        }
        intVal, err := strconv.Atoi(input)
        if err != nil {
            fmt.Println("Input is not an integer value.")
        } else {
            if count == 0 || intVal > max {
                max = intVal
            }
            if count == 0 || intVal < min {
                min = intVal
            }
            count += 1
            total += intVal
        }
    }
    mean := float64(total)/ float64(count)
    fmt.Printf("Min: %d, Max: %d, Mean: %f\n", min, max, mean)
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    printCoins(scanner)
}
