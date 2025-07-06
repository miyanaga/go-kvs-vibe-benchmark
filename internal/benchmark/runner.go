package benchmark

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/miyanaga/go-kvs-vibe-benchmark/internal/kvs"
)

type Result struct {
	Library    string
	AppendMS   int64
	UpdateMS   int64
	GetMS      int64
	FileSizeB  int64
}

type Runner struct {
	dataDir string
}

func NewRunner(dataDir string) *Runner {
	return &Runner{dataDir: dataDir}
}

func (r *Runner) RunBenchmark(kvsImpl kvs.KVS) (*Result, error) {
	result := &Result{Library: kvsImpl.Name()}
	
	dbPath := filepath.Join(r.dataDir, kvsImpl.Name())
	
	if err := os.RemoveAll(dbPath); err != nil {
		return nil, fmt.Errorf("failed to remove existing db: %w", err)
	}
	
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db dir: %w", err)
	}
	
	if err := kvsImpl.Open(dbPath); err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	defer kvsImpl.Close()
	
	testData, err := r.loadTestData()
	if err != nil {
		return nil, fmt.Errorf("failed to load test data: %w", err)
	}
	
	start := time.Now()
	for _, item := range testData {
		value := &kvs.Value{Single: &item.Value}
		if err := kvsImpl.Set(item.Key, value); err != nil {
			return nil, fmt.Errorf("append failed: %w", err)
		}
	}
	result.AppendMS = time.Since(start).Milliseconds()
	fmt.Printf("%s append: %dms\n", kvsImpl.Name(), result.AppendMS)
	
	start = time.Now()
	for _, item := range testData {
		doubleValue := item.Value * 2
		value := &kvs.Value{Double: &doubleValue}
		if err := kvsImpl.Set(item.Key, value); err != nil {
			return nil, fmt.Errorf("update failed: %w", err)
		}
	}
	result.UpdateMS = time.Since(start).Milliseconds()
	fmt.Printf("%s update: %dms\n", kvsImpl.Name(), result.UpdateMS)
	
	start = time.Now()
	for _, item := range testData {
		value, err := kvsImpl.Get(item.Key)
		if err != nil {
			return nil, fmt.Errorf("get failed: %w", err)
		}
		
		expectedDouble := item.Value * 2
		if value.Double == nil || *value.Double != expectedDouble {
			return nil, fmt.Errorf("validation failed for key %s: expected double=%d, got %v", 
				item.Key, expectedDouble, value.Double)
		}
	}
	result.GetMS = time.Since(start).Milliseconds()
	fmt.Printf("%s get: %dms\n", kvsImpl.Name(), result.GetMS)
	
	fileSize, err := r.calculateDirectorySize(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate file size: %w", err)
	}
	result.FileSizeB = fileSize
	fmt.Printf("%s file size: %d bytes\n", kvsImpl.Name(), result.FileSizeB)
	
	return result, nil
}

type TestItem struct {
	Key   string
	Value int
}

func (r *Runner) loadTestData() ([]TestItem, error) {
	file, err := os.Open(filepath.Join(r.dataDir, "keys.tsv"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	var items []TestItem
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			continue
		}
		
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		
		items = append(items, TestItem{
			Key:   parts[0],
			Value: value,
		})
	}
	
	return items, scanner.Err()
}

func (r *Runner) calculateDirectorySize(dir string) (int64, error) {
	var size int64
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}