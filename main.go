package main

import (
	"encoding/csv"
	"fmt"
	"go.uber.org/zap"
	"license/log"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Application struct {
	ComputerId    string
	UserId        string
	ApplicationId string
	ComputerType  string
	Comment       string
}

const (
	filePath      = "sample-large.csv"
	applicationId = "374"

	// Set the number of goroutines (adjust based on system's capabilities)
	numGoroutines = 4
)

var (
	logger *zap.Logger
	mu     sync.Mutex

	// Variable to store the total minimum copies
	minimumCopies int

	// Map to store Laptop and Desktop count
	computerCounts map[string]int

	// Used to identify duplicate records
	duplicates []string
)

func main() {
	// Initialize structured logging.
	logger = log.InitLogger("info")
	defer logger.Sync()

	startTime := time.Now()

	computerCounts = make(map[string]int)

	dataChannel := make(chan []Application, numGoroutines)
	done := make(chan bool)
	var wg sync.WaitGroup

	// Process the csv data concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		logger.Debug("Stared goroutine", zap.Int("Count", i))
		go processCSVData(dataChannel, &wg)
	}

	// Read the CSV file concurrently
	readCSV(filePath, dataChannel, numGoroutines, done)

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		done <- true
	}()

	<-done

	logger.Info(strconv.Itoa(minimumCopies) + " copies required for " + applicationId + " for file " + filePath + " time taken: " + time.Since(startTime).String())
}

func readCSV(filePath string, dataChannel chan<- []Application, numGoroutines int, done chan<- bool) {
	defer close(dataChannel)
	var wgCsv sync.WaitGroup

	file, err := os.Open(filePath)
	if err != nil {
		logger.Error("Error opening file", zap.Error(err))
		done <- false
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		logger.Error("Error reading CSV", zap.Error(err))
		done <- false
		return
	}

	// Ignore header
	lines = lines[1:]

	// Split the lines among goroutines
	chunkSize := len(lines) / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		startIndex := i * chunkSize
		endIndex := (i + 1) * chunkSize
		if i == numGoroutines-1 {
			endIndex = len(lines)
		}

		wgCsv.Add(1)
		logger.Info("Add")
		go func(start, end int) {
			processChunk(lines[start:end], dataChannel)
			wgCsv.Done()
		}(startIndex, endIndex)
	}

	wgCsv.Wait()
}

func processCSVData(dataChannel <-chan []Application, wg *sync.WaitGroup) {

	for apps := range dataChannel {
		for _, app := range apps {
			if app.ApplicationId == applicationId {
				app.ComputerType = strings.ToUpper(app.ComputerType)
				updated(app)
			}
		}
	}

	wg.Done()
}

func updated(app Application) {
	mu.Lock()
	defer mu.Unlock()

	// Do not process duplicate records
	str := fmt.Sprintf("%s", app.ComputerId)
	if slices.Contains(duplicates, str) {
		return
	}
	duplicates = append(duplicates, str)

	desktopKey := fmt.Sprintf("%s-%s", app.UserId, "DESKTOP")
	laptopKey := fmt.Sprintf("%s-%s", app.UserId, "LAPTOP")

	if app.ComputerType == "DESKTOP" {

		laptopCount := computerCounts[laptopKey]

		if laptopCount == 0 {
			minimumCopies++
			computerCounts[desktopKey]++
		} else {
			computerCounts[laptopKey]--
		}

	} else if app.ComputerType == "LAPTOP" {

		desktopCount := computerCounts[desktopKey]

		if desktopCount == 0 {
			minimumCopies++
			computerCounts[laptopKey]++
		} else {
			computerCounts[desktopKey]--
		}
	}
}

func processChunk(chunk [][]string, dataChannel chan<- []Application) {
	var apps []Application
	for _, line := range chunk {
		app := Application{
			ComputerId:    line[0],
			UserId:        line[1],
			ApplicationId: line[2],
			ComputerType:  line[3],
			Comment:       line[4],
		}
		apps = append(apps, app)
	}

	dataChannel <- apps
}
