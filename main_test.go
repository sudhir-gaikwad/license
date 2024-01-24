package main

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestProcessCSVData(t *testing.T) {
	// Initialize the test data
	testData := []Application{
		{ComputerId: "1", UserId: "1", ApplicationId: "374", ComputerType: "DESKTOP", Comment: "Comment1"},
		{ComputerId: "2", UserId: "1", ApplicationId: "374", ComputerType: "LAPTOP", Comment: "Comment2"},
		{ComputerId: "3", UserId: "2", ApplicationId: "374", ComputerType: "DESKTOP", Comment: "Comment3"},
		{ComputerId: "4", UserId: "2", ApplicationId: "374", ComputerType: "DESKTOP", Comment: "Comment4"},
	}

	// Mocking the data channel
	dataChannel := make(chan []Application, 1)
	computerCounts = make(map[string]int)

	defer close(dataChannel)

	// Creating a wait group to synchronize the completion of the goroutine
	var wg sync.WaitGroup
	wg.Add(1)

	// Run the goroutine with test data
	go func() {
		defer wg.Done()
		processCSVData(dataChannel, &wg)
	}()

	// Send the test data to the data channel
	dataChannel <- testData

	// Close the data channel to signal the end of data
	close(dataChannel)

	// Wait for the goroutine to complete
	wg.Wait()

	// Assert the expected results based on the test data
	assert.Equal(t, 3, minimumCopies) // One copy for each user, as Desktop and Laptop come in pairs
	// Add more assertions based on your test cases
}

func TestUpdated(t *testing.T) {
	// Reset global variables for each test case
	minimumCopies = 0
	computerCounts = make(map[string]int)
	duplicates = []string{}

	// Test case 1: One Desktop and one Laptop, one copy expected
	app1 := Application{ComputerType: "DESKTOP", UserId: "1"}
	app2 := Application{ComputerType: "LAPTOP", UserId: "1"}
	updated(app1)
	updated(app2)
	assert.Equal(t, 1, minimumCopies)

	// Test case 2: Two Desktops and two Laptops, two copies expected (one copy for each pair)
	app3 := Application{ComputerType: "DESKTOP", UserId: "2"}
	app4 := Application{ComputerType: "DESKTOP", UserId: "2"}
	app5 := Application{ComputerType: "LAPTOP", UserId: "2"}
	app6 := Application{ComputerType: "LAPTOP", UserId: "2"}
	updated(app3)
	updated(app4)
	updated(app5)
	updated(app6)
	assert.Equal(t, 2, minimumCopies)

	// Test case 3: One Desktop and three Laptops, three copies expected (one copy for each Laptop)
	app7 := Application{ComputerType: "DESKTOP", UserId: "3"}
	app8 := Application{ComputerType: "LAPTOP", UserId: "3"}
	app9 := Application{ComputerType: "LAPTOP", UserId: "3"}
	app10 := Application{ComputerType: "LAPTOP", UserId: "3"}
	updated(app7)
	updated(app8)
	updated(app9)
	updated(app10)
	assert.Equal(t, 3, minimumCopies)
	// Add more test cases as needed
}
