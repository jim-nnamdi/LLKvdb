package model

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// HeapItem represents a single key-value pair with the source file index
type HeapItem struct {
	Key     int64  // The key of the current item
	Value   string // The value of the current item
	FileIdx int    // The index of the file this item came from
}

// MinHeap implements the heap.Interface for a min-heap of HeapItem
type MinHeap []HeapItem

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].Key < h[j].Key }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(HeapItem))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *MinHeap) Init() {
	n := len(*h)
	for i := n/2 - 1; i >= 0; i-- { // Start from the first non-leaf node
		h.siftDown(i, n)
	}
}

func (h *MinHeap) siftDown(i, n int) {
	smallest := i
	left := 2*i + 1
	right := 2*i + 2

	if left < n && (*h)[left].Key < (*h)[smallest].Key {
		smallest = left
	}
	if right < n && (*h)[right].Key < (*h)[smallest].Key {
		smallest = right
	}

	if smallest != i {
		h.Swap(i, smallest)
		h.siftDown(smallest, n)
	}
}

// MergeSSTables merges multiple sorted SSTable files into one sorted output file
func (h *MinHeap) MergeSSTables(sstableFiles []string, outputFile string) error {
	// Open all input files
	readers := make([]*os.File, len(sstableFiles))
	scanners := make([]*bufio.Scanner, len(sstableFiles))
	for i, file := range sstableFiles {
		f, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", file, err)
		}
		readers[i] = f
		scanners[i] = bufio.NewScanner(f)
		scanners[i].Split(bufio.ScanLines)
	}

	// Create the output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)

	// Initialize the Min-Heap
	heap := &MinHeap{}
	for i, scanner := range scanners {
		if scanner.Scan() {
			line := scanner.Text()
			key, value, err := parseKeyValue(line)
			if err != nil {
				return fmt.Errorf("failed to parse key-value: %w", err)
			}
			heap.Push(HeapItem{Key: key, Value: value, FileIdx: i})
		}
	}
	heap.Init() // Initialize the heap structure

	// Perform k-way merge
	for heap.Len() > 0 {
		// Extract the smallest item
		item := heap.Pop().(HeapItem)

		// Write the item to the output file
		_, err := writer.WriteString(fmt.Sprintf("%d:%s\n", item.Key, item.Value))
		if err != nil {
			return fmt.Errorf("failed to write to output file: %w", err)
		}

		// Push the next key-value pair from the same file into the heap
		scanner := scanners[item.FileIdx]
		if scanner.Scan() {
			line := scanner.Text()
			key, value, err := parseKeyValue(line)
			if err != nil {
				return fmt.Errorf("failed to parse key-value: %w", err)
			}
			heap.Push(HeapItem{Key: key, Value: value, FileIdx: item.FileIdx})
		}
	}

	// Flush the writer and close files
	writer.Flush()
	for _, f := range readers {
		f.Close()
	}
	return nil
}

// Helper function to parse a key-value pair from a string
func parseKeyValue(line string) (int64, string, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid key-value format: %s", line)
	}
	key, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid key: %s", parts[0])
	}
	return key, parts[1], nil
}
