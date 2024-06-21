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
	"github.com/wbergg/strongbeer-bot/helper"
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

func (ss *SheetService) GetStatus(username string) (string, error) {

	// Read cell
	value, err := ss.GetCell(username, "B")
	if err != nil {
		return "", fmt.Errorf("unable to retrieve data from cell: %v", err)
	}

	// Read cell
	points, err := ss.GetCell(username, "AE")
	if err != nil {
		return "", fmt.Errorf("unable to retrieve data from cell: %v", err)
	}

	// Read range
	resp, err := ss.GetRange(username, "C", "Z")
	if err != nil {
		return "", fmt.Errorf("unable to retrieve data from cellrange: %v", err)
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

func (ss *SheetService) GetCheckin(username string) (string, error) {
	// Read cell
	row, ok := UserMap[username]
	if !ok {
		return "", fmt.Errorf("cannot find username in config: %v", ok)
	}

	value, err := ss.GetCell(row, "B")
	if err != nil {
		return "", fmt.Errorf("unable to retrieve data from cell: %v", err)
	}

	// Read range for row 7 to get weeks
	resp, err := ss.GetRange("7", "C", "AB")
	if err != nil {
		return "", fmt.Errorf("unable to retrieve data from cellrange: %v", err)
	}

	// Figure out week / column relation
	tn := time.Now()
	_, week := tn.ISOWeek()
	var checkin string

	if len(resp.Values) > 0 {
		for column, cell := range resp.Values[0] {
			cellStr, _ := cell.(string)
			cs, _ := strconv.Atoi(cellStr)
			if cs == week {
				c := helper.GetColumnLetter(column + 2)
				checkin, err = ss.GetCell(row, c)
				if err != nil {
					return "", fmt.Errorf("unable to retrieve data from cell: %v", err)
				}
			}
		}
	}

	message := "*" + value + " (" + username + ")*" + "\n\n" +
		"Vecka: " + strconv.Itoa(week) + "\n\n"
	if checkin == "x" {
		message = message + "You are *CHECKED IN*"
	} else {
		message = message + "You are *NOT CHECKED IN*"
	}

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

func (ss *SheetService) GetCell(row string, cell string) (string, error) {

	// Read specific cell
	readCell := config.Loaded.Sheetname + "!" + cell + row
	resp, err := ss.srv.Spreadsheets.Values.Get(ss.spreadsheetId, readCell).Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve data from cell: %v", err)
	}

	var value string
	if len(resp.Values) > 0 && len(resp.Values[0]) > 0 {
		value = resp.Values[0][0].(string)
	} else {
		value = ""
	}

	return value, err
}

func (ss *SheetService) GetRange(row string, cellfrom string, cellto string) (*sheets.ValueRange, error) {

	// Read range

	readRange := config.Loaded.Sheetname + "!" + cellfrom + row + ":" + cellto + row
	resp, err := ss.srv.Spreadsheets.Values.Get(ss.spreadsheetId, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from range: %v", err)
	}

	return resp, err
}
