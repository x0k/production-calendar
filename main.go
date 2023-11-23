package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	WEEKEND     = 1
	HOLIDAY     = 2
	PRE_HOLIDAY = 3
)

func getDayType(dayText string) (int, error) {
	lowered := strings.ToLower(dayText)
	switch {
	case strings.HasPrefix(lowered, "праздничный"):
		return HOLIDAY, nil
	case strings.HasPrefix(lowered, "выходной"):
		return WEEKEND, nil
	case strings.Contains(lowered, "сокращен"):
		return PRE_HOLIDAY, nil
	default:
		return 0, fmt.Errorf("unknown day type: %s", dayText)
	}
}

var CALENDAR_YEAR_REG_EXP = regexp.MustCompile(`\d+`)

var MONTHS = []string{
	"yanvar",
	"fevral",
	"mart",
	"aprel",
	"maj",
	"iyun",
	"iyul",
	"avgust",
	"sentyabr",
	"oktbyar",
	"noyabr",
	"dekabr",
}

func makeProductionCalendar(doc *goquery.Document, out *map[string]int) {
	calendarTitleElement := doc.Find("h1").First()
	if calendarTitleElement.Length() == 0 {
		log.Fatal("calendar title not found")
	}
	calendarYear := CALENDAR_YEAR_REG_EXP.FindString(calendarTitleElement.Text())
	if len(calendarYear) == 0 {
		log.Fatal("calendar year not found")
	}
	calendarRootElement := calendarTitleElement.Parent().Parent().Parent()
	if calendarRootElement.Length() == 0 {
		log.Fatal("calendar root element not found")
	}
	for i, month := range MONTHS {
		monthRootElement := calendarRootElement.Find(fmt.Sprintf("[href$=\"/%s/\"]", month)).First().Parent().Parent()
		if monthRootElement.Length() == 0 {
			log.Fatalf("month root element not found: %s", month)
		}
		monthNumber := fmt.Sprintf("%02d", i+1)
		monthRootElement.Find("._1c_LS,._1YS-8").Each(func(i int, dayElement *goquery.Selection) {
			dayNumberText := dayElement.Find("[role=\"button\"]").Text()
			dayNumber, err := strconv.Atoi(dayNumberText)
			if err != nil {
				log.Fatal(err)
			}
			dayText := dayElement.Find("[role=\"tooltip\"]").Children().First().Text()
			dayType, err := getDayType(dayText)
			if err != nil {
				log.Fatal(err)
			}
			(*out)[fmt.Sprintf("%s-%s-%s", calendarYear, monthNumber, fmt.Sprintf("%02d", dayNumber))] = dayType
		})
	}
}

func main() {
	doc, err := goquery.NewDocumentFromReader(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	result := make(map[string]int)
	makeProductionCalendar(doc, &result)
	data, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(data)
}
