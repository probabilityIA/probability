package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	filename := ".env.local"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading: %v", err)
	}
	fmt.Printf("File content (len=%d):\n%s\n", len(content), string(content))

	// Check BOM
	if len(content) > 3 && content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		fmt.Println("File has UTF-8 BOM")
	}
}
