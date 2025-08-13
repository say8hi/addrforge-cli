package util

import (
	"fmt"
	"os"
)

func SaveResult(filename, result string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error saving to file: %v\n", err)
		fmt.Print(result) // Fallback to console
		return
	}
	defer file.Close()

	file.WriteString(result + "\n")
	fmt.Printf("Result saved to %s\n", filename)
}
