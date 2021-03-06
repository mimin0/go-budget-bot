package spreadSheerReader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

func Reader(spreadsheetId string) string {
	b, err := ioutil.ReadFile("spreadSheet/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	// scope list https://developers.google.com/sheets/api/guides/authorizing
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	readRange := "Sheet!A2:D"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	bData := []string{}
	var stringData string
	// var CountStr string

	if len(resp.Values) == 0 {
		stringData = "No data found."
		return stringData
	} else {
		fmt.Println("Date, Expencis, Amount, Type:")
		for _, row := range resp.Values {
			// Print columns A and E, which correspond to indices 0 and 4.
			//fmt.Printf("%s, %s, %s, %s\n", row[0], row[1], row[2], row[3])
			bData = append(bData, fmt.Sprintf("%v", row[0]), fmt.Sprintf("%v", row[1]), fmt.Sprintf("%v", row[2]), fmt.Sprintf("%v\n", row[3]))
		}
	}
	stringData = strings.Join(bData, ",")
	return stringData
}

func Writer(spreadsheetId string, messageDate int, textMessage string) string {
	b, err := ioutil.ReadFile("spreadSheet/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	//parsing time of adding new items
	i, err := strconv.ParseInt(fmt.Sprintf("%v", messageDate), 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	goodTimeFormat := tm.Format("2006/01/02 15:04:05")

	// If modifying these scopes, delete your previously saved token.json.
	// scope list https://developers.google.com/sheets/api/guides/authorizing
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	type body struct {
		Data struct {
			Range  string     `json:"range"`
			Values [][]string `json:"values"`
		} `json:"data"`
		ValueInputOption string `json:"valueInputOption"`
	}

	// get the last filled row in a spread sheet
	currentRange := "Sheet!A2:D"
	var stringL []string
	Resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, currentRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from  sheet: %v", err)
	}
	index := len(Resp.Values)
	stringL = []string{"Sheet!A", fmt.Sprintf("%v", index), ":D", fmt.Sprintf("%v", index)}

	readRange := strings.Join(stringL, "")
	text := strings.Split(textMessage, " ")
	values := [][]interface{}{{goodTimeFormat, text[1], text[2], text[3]}}
	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
	}
	rb.Data = append(rb.Data, &sheets.ValueRange{
		Range:  readRange,
		Values: values,
	})
	_, err = srv.Spreadsheets.Values.BatchUpdate(spreadsheetId, rb).Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")
	done := "done!"
	return done
}
