package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cartabinaria/unibo-go/timetable"
	"golang.org/x/exp/slices"
)

const TIMEFORMAT = "2006-01-02T15:04:05"

type Aula struct {
	Edificio  string `json:"des_edificio"`
	Piano     string `json:"des_piano"`
	Indirizzo string `json:"des_indirizzo"`
}

type LezioniTime time.Time

type OrarioLezioni struct {
	Title     string      `json:"title"`
	Time      string      `json:"time"`
	Aule      []Aula      `json:"aule"`
	Teams     string      `json:"teams"`
	StartTime LezioniTime `json:"start"`
	EndTime   LezioniTime `json:"end"`
}

func (t *LezioniTime) Format(format string) string {
	return (*time.Time)(t).Format(format)
}

func (t *LezioniTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsedTime, err := time.Parse(TIMEFORMAT, s)
	if err != nil {
		return err
	}
	*t = LezioniTime(parsedTime)
	return nil
}

// GetTimeTable returns an HTML string containing the timetable for the given
// course on the given date. Returns an empty string if there are no lessons.
func GetTimeTable(courseType, courseName string, curriculum string, year int, day time.Time) (string, error) {

	interval := &timetable.Interval{Start: day, End: day}
	events, err := timetable.FetchTimetable(courseType, courseName, curriculum, year, interval)
	if err != nil {
		log.Printf("Error getting timetable: %s\n", err)
		return "", err
	}

	// Sort the events by start time
	slices.SortFunc(events, func(a, b timetable.Event) int {
		return int(a.Start.Time.Sub(b.Start.Time).Nanoseconds())
	})

	b := strings.Builder{}
	for _, event := range events {
		b.WriteString(fmt.Sprintf(`  🕘 <b><a href="%s">%s</a></b>`, event.Teams, event.Title))
		b.WriteString("\n")
		b.WriteString(event.Start.Format("15:04") + " - " +
			event.End.Format("15:04"))
		b.WriteString("\n")
		if len(event.Classrooms) > 0 {
			b.WriteString(fmt.Sprintf("  🏢 %s - %s\n", event.Classrooms[0].BuildingDesc, event.Classrooms[0].FloorDesc))
			b.WriteString(fmt.Sprintf("  📍 %s\n", event.Classrooms[0].AddressDesc))
		}
	}

	return b.String(), nil
}
