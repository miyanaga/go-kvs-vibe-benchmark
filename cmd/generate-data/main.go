package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	
	file, err := os.Create("data/keys.tsv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for i := 0; i < 100000; i++ {
		randomValue := rand.Int()
		hash := sha256.Sum256([]byte(fmt.Sprintf("%d", i)))
		hashStr := fmt.Sprintf("%x", hash)
		
		_, err := fmt.Fprintf(file, "%s\t%d\n", hashStr, randomValue)
		if err != nil {
			panic(err)
		}
	}
	
	fmt.Println("Generated 100,000 test records in data/keys.tsv")
}