package main

import (
	"bufio"
	"container/list"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
	"unbabel-challenge/internal/models"
)

func main() {

	inputFile, windowSize := getConsoleArguments()

	file, err := openFile(inputFile)
	if err != nil {
		return
	}
	defer file.Close()

	queue := list.New()
	var sum float64
	var currentTime time.Time
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event models.Translation
		err := json.Unmarshal(scanner.Bytes(), &event)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		timestamp, _ := parseTimestamp(event.Timestamp)
		if currentTime.IsZero() {
			currentTime = timestamp
		}

		// Output the average delivery time for each minute within the window
		for currentTime.Before(timestamp) || currentTime == timestamp {

			removed := removeElementsOutsideTimeWindow(queue, currentTime, windowSize)
			for _, elem := range removed {
				sum -= elem.Duration
			}
			printResultLine(queue, sum, currentTime)

			currentTime = currentTime.Add(time.Minute)
		}

		queue.PushBack(event)
		sum += event.Duration
	}

	printResultLine(queue, sum, currentTime)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}

func printResultLine(queue *list.List, sum float64, currentTime time.Time) {
	averageDeliveryTime := float32(0)
	if queue.Len() != 0 {
		averageDeliveryTime = float32(sum) / float32(queue.Len())
	}

	output := models.Output{
		Timestamp:           currentTime.Format("2006-01-02 15:04:05"),
		AverageDeliveryTime: averageDeliveryTime,
	}
	jsonOutput, _ := json.Marshal(output)
	fmt.Println(string(jsonOutput))
}

// Remove events that are out of the window from the front of the queue(oldest)
func removeElementsOutsideTimeWindow(queue *list.List, currentTime time.Time, windowSize int) []models.Translation {
	var removed []models.Translation

	for queue.Len() > 0 {
		front := queue.Front().Value.(models.Translation)

		frontTimestamp, _ := parseTimestamp(front.Timestamp)
		if currentTime.Sub(frontTimestamp) > time.Duration(windowSize)*time.Minute {
			queue.Remove(queue.Front())
			removed = append(removed, front)

		} else {
			break
		}
	}

	return removed
}

func openFile(inputFile string) (*os.File, error) {
	file, err := os.Open(inputFile)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	return file, nil
}

func parseTimestamp(s string) (time.Time, error) {
	timestamp, err := time.Parse("2006-01-02 15:04:05.999999", s)
	if err != nil {
		return time.Time{}, err
	}
	return timestamp.Truncate(time.Minute), nil
}

func getConsoleArguments() (string, int) {
	inputFile := flag.String("input_file", "events.txt", "Path to the input file")
	windowSize := flag.Int("window_size", 10, "Window size in minutes")
	flag.Parse()

	// Check if required flags are provided
	if *inputFile == "" {
		fmt.Println("Error: input_file flag is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *windowSize <= 0 {
		fmt.Println("Error: window_size should be greater than 0")
		flag.PrintDefaults()
		os.Exit(1)
	}

	return *inputFile, *windowSize
}
