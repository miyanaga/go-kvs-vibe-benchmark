package main

import (
	"fmt"
	"log"

	"github.com/miyanaga/go-kvs-vibe-benchmark/internal/benchmark"
	"github.com/miyanaga/go-kvs-vibe-benchmark/internal/kvs"
)

func main() {
	runner := benchmark.NewRunner("data")
	
	kvsLibraries := []kvs.KVS{
		kvs.NewLevelDB(),
		kvs.NewBBolt(),
		kvs.NewBadger(),
		kvs.NewPebble(),
		kvs.NewSQLite(),
	}
	
	var results []*benchmark.Result
	
	for _, kvsImpl := range kvsLibraries {
		fmt.Printf("Running benchmark for %s...\n", kvsImpl.Name())
		result, err := runner.RunBenchmark(kvsImpl)
		if err != nil {
			log.Printf("Error benchmarking %s: %v", kvsImpl.Name(), err)
			continue
		}
		results = append(results, result)
		fmt.Printf("%s benchmark completed\n\n", kvsImpl.Name())
	}
	
	fmt.Println("Benchmark Results:")
	fmt.Println("Library\tAppend(ms)\tUpdate(ms)\tGet(ms)\tFileSize(bytes)")
	for _, result := range results {
		fmt.Printf("%s\t%d\t%d\t%d\t%d\n",
			result.Library,
			result.AppendMS,
			result.UpdateMS,
			result.GetMS,
			result.FileSizeB,
		)
	}
}