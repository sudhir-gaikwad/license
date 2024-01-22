package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

type Application struct {
	ComputerId    string
	UserId        string
	ApplicationId string
	ComputerType  string
	Comment       string
}

type Data struct {
	LaptopCount  int
	DesktopCount int
}

func main() {
	fmt.Println(time.Now().Format(time.RFC850))

	applicationId := "374"
	file := "sample-large.csv"

	data, err := readFile(file)
	if err != nil {
		fmt.Println("Error in processing request:", err)
		return
	}

	copiesRequired := calculateCopiesRequired(data, applicationId)

	fmt.Printf("%d copies required for application %s, Input file: %s \n",
		copiesRequired, applicationId, file)
	fmt.Println(time.Now().Format(time.RFC850))
}

func readFile(filename string) ([]Application, error) {
	fmt.Printf("Reading file: %s\n", filename)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []Application
	// Ignore header
	for _, line := range lines[1:] {
		app := Application{
			ComputerId:    line[0],
			UserId:        line[1],
			ApplicationId: line[2],
			ComputerType:  line[3],
			Comment:       line[4],
		}
		data = append(data, app)
	}
	return data, nil
}

func calculateCopiesRequired(applications []Application, applicationId string) int {

	// Key: UserId-ApplicationId -> Uniquely maintain the count
	// Value: Struct contains Laptop count and Desktop count
	appCountMap := buildAppCountMap(applications, applicationId)
	fmt.Println("buildAppCountMap--------------------------------")

	// Calculate the total number of licenses required
	totalLicenses := 0
	for _, data := range appCountMap {

		laptopCount := data.LaptopCount
		desktopCount := data.DesktopCount

		//fmt.Printf("Key: %s LaptopCount: %d DesktopCount: %d\n", key, laptopCount, desktopCount)

		if desktopCount == laptopCount {
			// If count of laptop and desktop are same then consider only one of them
			totalLicenses = totalLicenses + desktopCount
		} else {

			total := desktopCount + laptopCount
			if desktopCount > laptopCount {
				// Reduce laptop count from the total count
				totalLicenses = totalLicenses + (total - laptopCount)
			} else {
				// Reduce desktop count from the total count
				totalLicenses = totalLicenses + (total - desktopCount)
			}
		}
	}

	return totalLicenses
}

func buildAppCountMap(applications []Application, applicationId string) map[string]Data {
	// Key: UserId-ApplicationId -> Uniquely maintain the count
	// Value: Struct contains Laptop count and Desktop count
	appCountMap := make(map[string]Data)

	// Used to identify duplicate records
	var duplicates []string

	for _, app := range applications {

		// Ignore the records which do not match with the given application id
		if app.ApplicationId != applicationId {
			continue
		}

		app.ComputerType = strings.ToUpper(app.ComputerType)

		// Create a unique key for each user and application
		key := fmt.Sprintf("%s-%s", app.UserId, app.ApplicationId)
		//fmt.Printf("key------------: %s\n", key)

		// Do not process duplicate records
		str := fmt.Sprintf("%s-%s-%s", app.ComputerId, app.UserId, app.ApplicationId)
		if slices.Contains(duplicates, str) {
			continue
		}
		duplicates = append(duplicates, str)

		if app.ComputerType == "LAPTOP" {

			previousData, found := appCountMap[key]
			if found {
				previousData.LaptopCount++
				appCountMap[key] = previousData
			} else {
				result := Data{
					LaptopCount:  1,
					DesktopCount: 0,
				}
				appCountMap[key] = result
			}
		} else {
			previousData, found := appCountMap[key]
			if found {
				previousData.DesktopCount++
				appCountMap[key] = previousData
			} else {
				result := Data{
					LaptopCount:  0,
					DesktopCount: 1,
				}
				appCountMap[key] = result
			}
		}
	}

	return appCountMap
}
