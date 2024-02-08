package main

import (
	"encoding/csv"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strings"
	"unicode"
)

type Row struct {
	rank      string
	nick      string
	firstName string
	category  string
	followers string
	country   string
	engAuth   string
	engAvg    string
}

func main() {
	// init slice structs
	var Rows []Row

	c := colly.NewCollector()

	// valid user-agent
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/111.0.0.0 Safari/537.36"

	//HTML elements

	c.OnHTML(".table .row[data-v-2e6a30b8]", func(element *colly.HTMLElement) {

		rank := element.ChildText(".row-cell.rank span[data-v-2e6a30b8]")
		nick := element.ChildText(".contributor__name-content[data-v-c5a99f5a]")
		firstName := element.ChildText(".contributor__title[data-v-c5a99f5a]")
		category := element.ChildText(".row-cell.category[data-v-2e6a30b8] .ellipsis")
		followers := element.ChildText(".row-cell.subscribers[data-v-2e6a30b8]")
		country := element.ChildText(".row-cell.audience[data-v-e1ea9c14]")
		engAuth := element.ChildText(".row-cell.authentic[data-v-e1ea9c14]")
		engAvg := element.ChildText(".row-cell.engagement[data-v-e1ea9c14]")

		// add space between categories
		var categoryString strings.Builder

		for idx, val := range category {
			if idx > 0 {
				if unicode.IsUpper(val) && !unicode.IsUpper(rune(category[idx-1])) && unicode.IsLetter(val) {
					categoryString.WriteString(" ")
				}
			}
			categoryString.WriteRune(val)
		}

		// full result struct
		row := Row{
			rank:      rank,
			nick:      nick,
			firstName: firstName,
			category:  categoryString.String(),
			followers: followers,
			country:   country,
			engAuth:   engAuth,
			engAvg:    engAvg,
		}

		Rows = append(Rows, row)

	})

	c.Visit("https://hypeauditor.com/top-instagram-all-russia/")

	// --- export to CSV ---

	// open the output CSV file
	csvFile, csvErr := os.Create("Instagram .csv")
	// if the file creation fails
	if csvErr != nil {
		log.Fatalln("Failed to create the output CSV file", csvErr)
	}
	// release the resource allocated to handle
	// the file before ending the execution
	defer csvFile.Close()

	// create a CSV file writer
	writer := csv.NewWriter(csvFile)
	// release the resources associated with the
	// file writer before ending the execution
	defer writer.Flush()

	// add the header row to the CSV
	headers := []string{
		"rank",
		"nick",
		"Name",
		"category",
		"followers",
		"country",
		"engAuth",
		"engAvg",
	}
	writer.Write(headers)

	// store each Industry product in the
	// output CSV file
	for _, row := range Rows {
		// convert the Industry instance to
		// a slice of strings
		record := []string{
			row.rank,
			row.nick,
			row.firstName,
			row.category,
			row.followers,
			row.country,
			row.engAuth,
			row.engAvg,
		}

		// add a new CSV record
		writer.Write(record)
	}

}
