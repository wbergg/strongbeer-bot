package sheetservice

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/wbergg/strongbeer-bot/config"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetService struct {
	srv           *sheets.Service
	spreadsheetId string
}

var UserMap = make(map[string]string)

func InitSheetService(spreadsheetId, jsonKeyFile string) (*SheetService, error) {
	ctx := context.Background()

	// Read the JSON key file
	//jsonKeyFile := "creds/sheets.json"
	b, err := os.ReadFile(jsonKeyFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Parse the JSON key file
	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// Create a Sheets service client
	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Populate UserMap
	populateUserMap()

	return &SheetService{
		srv:           srv,
		spreadsheetId: spreadsheetId,
	}, nil
}

func populateUserMap() {
	// Populate usermap with data from config-file
	for key, value := range config.Loaded.Users {
		UserMap[key] = value
	}
}

func (ss *SheetService) GetSheetUserData(username string) (string, error) {
	row := UserMap[username]
	// Read specific cell
	readUserName := "H1!B" + row
	resp, err := ss.srv.Spreadsheets.Values.Get(ss.spreadsheetId, readUserName).Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve data from cell: %v", err)
	}

	var value string
	if len(resp.Values) > 0 && len(resp.Values[0]) > 0 {
		value = resp.Values[0][0].(string)
	} else {
		value = ""
	}

	// Read specific cell
	readPoints := "H1!AE" + row
	resp, err = ss.srv.Spreadsheets.Values.Get(ss.spreadsheetId, readPoints).Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve data from cell: %v", err)
	}

	var points string
	if len(resp.Values) > 0 && len(resp.Values[0]) > 0 {
		points = resp.Values[0][0].(string)
	} else {
		points = ""
	}

	// Read range
	readRange := "H1!C" + row + ":Z" + row
	resp, err = ss.srv.Spreadsheets.Values.Get(ss.spreadsheetId, readRange).Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve data from range: %v", err)
	}

	// Count occurrences of [xX] in the specified range
	count := 0
	if len(resp.Values) > 0 {
		for _, cell := range resp.Values[0] {
			if cell == "x" || cell == "X" {
				count++
			}
		}
	}

	message := "*" + value + " (" + username + ")*" + "\n\nYou have a total of " +
		strconv.Itoa(count) + " checkins out of 24" + "\n\nYou are at place " + points
	return message, nil
}

func (ss *SheetService) GetSheetUserCheckin(username string) (string, error) {
	row := UserMap[username]
	// Read specific cell
	readUserName := "H1!B" + row
	resp, err := ss.srv.Spreadsheets.Values.Get(ss.spreadsheetId, readUserName).Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve data from cell: %v", err)
	}

	var value string
	if len(resp.Values) > 0 && len(resp.Values[0]) > 0 {
		value = resp.Values[0][0].(string)
	} else {
		value = ""
	}

	// Read specific cell
	readPoints := "H1!AE" + row
	resp, err = ss.srv.Spreadsheets.Values.Get(ss.spreadsheetId, readPoints).Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve data from cell: %v", err)
	}

	var points string
	if len(resp.Values) > 0 && len(resp.Values[0]) > 0 {
		points = resp.Values[0][0].(string)
	} else {
		points = ""
	}

	// Read range
	readRange := "H1!C7:Z7"
	resp, err = ss.srv.Spreadsheets.Values.Get(ss.spreadsheetId, readRange).Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve data from range: %v", err)
	}

	// Figure out week / column relation
	tn := time.Now()
	_, week := tn.ISOWeek()
	count := 0
	//var column int
	if len(resp.Values) > 0 {
		for _, cell := range resp.Values[0] {
			cellStr, _ := cell.(string)
			fmt.Println(cell)
			cs, _ := strconv.Atoi(cellStr)
			if cs == week {
				//column = cell
				fmt.Println("det är nu vecka: ", week)
			}
		}
	}

	message := "*" + value + " (" + username + ")*" + "\n\nYou have a total of " +
		strconv.Itoa(count) + " checkins out of 24" + "\n\nYou are at place " + points
	return message, nil
}

func (ss *SheetService) GetSheetTopList(ll bool) (string, error) {

	// Read range for B9:B20 and AE9:AE20
	readRange := []string{"H1!B9:B20", "H1!AC9:AC20"}

	resp, err := ss.srv.Spreadsheets.Values.BatchGet(ss.spreadsheetId).Ranges(readRange...).Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve data from range: %v", err)
	}

	// Ensure the response contains both ranges
	if len(resp.ValueRanges) < 2 {
		return "", fmt.Errorf("expected two ranges in the response, but got %d", len(resp.ValueRanges))
	}

	// Initialize a slice to hold the row data
	var rows []map[string]string

	// Assuming both ranges have the same number of rows
	bValues := resp.ValueRanges[0].Values
	acValues := resp.ValueRanges[1].Values

	for i := 0; i < len(bValues); i++ {
		name := fmt.Sprintf("%v", bValues[i][0])
		poi := fmt.Sprintf("%v", acValues[i][0])
		row := map[string]string{"Name": name, "Points": poi}
		rows = append(rows, row)
	}

	// Sort the rows based on points value in descending order
	sort.Slice(rows, func(i, j int) bool {
		poiI, _ := strconv.Atoi(rows[i]["Points"])
		poiJ, _ := strconv.Atoi(rows[j]["Points"])
		return poiI > poiJ
	})

	var rowStrings []string
	i := 0
	var tr string

	for j, row := range rows {
		if ll && j == 3 {
			break
		}

		// Check if points compared to last loop match
		if tr != row["Points"] {
			i++
			tr = row["Points"]
		}

		// Convert rows to a string for the final message
		rowStrings = append(rowStrings, fmt.Sprintf("%s - %s (Points: %s)", strconv.Itoa(i), row["Name"], row["Points"]))
	}

	message := strings.Join(rowStrings, "\n")
	return message, nil
}
