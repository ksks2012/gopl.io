// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 187.

// Sorting sorts a music playlist into a variety of orders.
// Practice 7.8: Implement a sort function that takes a variable number of fields to sort by.
// Practice 7.9: Implement a web server that serves a sortable HTML table of tracks.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"
)

// !+main
type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

//!-main

// !+printTracks
func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush() // calculate column widths and print table
}

//!-printTracks

// !+artistcode
type byArtist []*Track

func (x byArtist) Len() int           { return len(x) }
func (x byArtist) Less(i, j int) bool { return x[i].Artist < x[j].Artist }
func (x byArtist) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

//!-artistcode

// !+yearcode
type byYear []*Track

func (x byYear) Len() int           { return len(x) }
func (x byYear) Less(i, j int) bool { return x[i].Year < x[j].Year }
func (x byYear) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

//!-yearcode

/*
//!+artistoutput
Title       Artist          Album              Year  Length
-----       ------          -----              ----  ------
Go Ahead    Alicia Keys     As I Am            2007  4m36s
Go          Delilah         From the Roots Up  2012  3m38s
Ready 2 Go  Martin Solveig  Smash              2011  4m24s
Go          Moby            Moby               1992  3m37s
//!-artistoutput

//!+artistrevoutput
Title       Artist          Album              Year  Length
-----       ------          -----              ----  ------
Go          Moby            Moby               1992  3m37s
Ready 2 Go  Martin Solveig  Smash              2011  4m24s
Go          Delilah         From the Roots Up  2012  3m38s
Go Ahead    Alicia Keys     As I Am            2007  4m36s
//!-artistrevoutput

//!+yearoutput
Title       Artist          Album              Year  Length
-----       ------          -----              ----  ------
Go          Moby            Moby               1992  3m37s
Go Ahead    Alicia Keys     As I Am            2007  4m36s
Ready 2 Go  Martin Solveig  Smash              2011  4m24s
Go          Delilah         From the Roots Up  2012  3m38s
//!-yearoutput

//!+customout
Title       Artist          Album              Year  Length
-----       ------          -----              ----  ------
Go          Moby            Moby               1992  3m37s
Go          Delilah         From the Roots Up  2012  3m38s
Go Ahead    Alicia Keys     As I Am            2007  4m36s
Ready 2 Go  Martin Solveig  Smash              2011  4m24s
//!-customout
*/

// !+customcode
type customSort struct {
	t    []*Track
	less func(x, y *Track) bool
}

func (x customSort) Len() int           { return len(x.t) }
func (x customSort) Less(i, j int) bool { return x.less(x.t[i], x.t[j]) }
func (x customSort) Swap(i, j int)      { x.t[i], x.t[j] = x.t[j], x.t[i] }

//!-customcode

func init() {
	//!+ints
	values := []int{3, 1, 4, 1}
	fmt.Println(sort.IntsAreSorted(values)) // "false"
	sort.Ints(values)
	fmt.Println(values)                     // "[1 1 3 4]"
	fmt.Println(sort.IntsAreSorted(values)) // "true"
	sort.Sort(sort.Reverse(sort.IntSlice(values)))
	fmt.Println(values)                     // "[4 3 1 1]"
	fmt.Println(sort.IntsAreSorted(values)) // "false"
	//!-ints
}

// input: array of fields
// output: sorted tracks by those fields ordered by the specified fields
func sortTracks(fields ...string) {
	// Implement sorting logic based on the provided fields.
	// This is a placeholder for the actual implementation.
	fmt.Println("Sorting tracks by fields:", fields)
	// You would typically use sort.Sort with a custom sort function here.
	sort.Sort(customSort{tracks, func(x, y *Track) bool {
		for _, field := range fields {
			switch field {
			case "Title":
				if x.Title != y.Title {
					return x.Title < y.Title
				}
			case "Artist":
				if x.Artist != y.Artist {
					return x.Artist < y.Artist
				}
			case "Album":
				if x.Album != y.Album {
					return x.Album < y.Album
				}
			case "Year":
				if x.Year != y.Year {
					return x.Year < y.Year
				}
			case "Length":
				if x.Length != y.Length {
					return x.Length < y.Length
				}
			}
		}
		return false
	}})
}

var currentSortFields []string

const maxSortFields = 3 // For example, top 3 sort keys

var trackListTemplate = template.Must(template.New("tracklist").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Music Tracks</title>
    <style>
        body { font-family: sans-serif; margin: 20px; }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
            cursor: pointer;
            position: relative; /* For sort indicator */
        }
        th a {
            text-decoration: none;
            color: black;
            display: block;
            padding: 8px; /* Make the whole th area clickable */
            margin: -8px; /* Adjust margin to fill th area */
        }
        th a:hover {
            color: blue;
        }
        .sort-indicator {
            position: absolute;
            right: 5px;
            top: 50%;
            transform: translateY(-50%);
            font-size: 0.8em;
            color: gray;
        }
    </style>
</head>
<body>
    <h1>Music Tracks</h1>
    <p>Click on a column header to sort. The order of sort keys is indicated by numbers.</p>
    <table>
        <thead>
            <tr>
                <th><a href="?sort=Title">Title<span class="sort-indicator">{{index .SortIndicators "Title"}}</span></a></th>
                <th><a href="?sort=Artist">Artist<span class="sort-indicator">{{index .SortIndicators "Artist"}}</span></a></th>
                <th><a href="?sort=Album">Album<span class="sort-indicator">{{index .SortIndicators "Album"}}</span></a></th>
                <th><a href="?sort=Year">Year<span class="sort-indicator">{{index .SortIndicators "Year"}}</span></a></th>
                <th><a href="?sort=Length">Length<span class="sort-indicator">{{index .SortIndicators "Length"}}</span></a></th>
            </tr>
        </thead>
        <tbody>
            {{range .Tracks}}
            <tr>
                <td>{{.Title}}</td>
                <td>{{.Artist}}</td>
                <td>{{.Album}}</td>
                <td>{{.Year}}</td>
                <td>{{.Length}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
    <p>Current Sort Order: {{.SortOrderStr}}</p>
</body>
</html>`))

// trackListData is the data structure passed to the HTML template.
type trackListData struct {
	Tracks         []*Track
	SortOrderStr   string
	SortIndicators map[string]string // e.g., {"Title": "¹", "Artist": "²"}
}

// trackList HTTP handler.
func trackList(w http.ResponseWriter, r *http.Request) {
	sortField := r.URL.Query().Get("sort")

	if sortField != "" {
		tempFields := make([]string, 0, maxSortFields)
		tempFields = append(tempFields, sortField)
		for _, field := range currentSortFields {
			if field != sortField {
				tempFields = append(tempFields, field)
			}
		}
		if len(tempFields) > maxSortFields {
			currentSortFields = tempFields[:maxSortFields]
		} else {
			currentSortFields = tempFields
		}

		sortTracks(currentSortFields...)
	} else {
		if len(currentSortFields) == 0 {
			currentSortFields = []string{"Title"}
			sortTracks(currentSortFields...)
		}
	}

	// Prepare sort indicators for the template.
	sortIndicators := make(map[string]string)
	superscripts := []string{"¹", "²", "³", "⁴", "⁵"}
	for i, field := range currentSortFields {
		if i < len(superscripts) {
			sortIndicators[field] = superscripts[i]
		}
	}

	// Prepare data for the template.
	data := trackListData{
		Tracks:         tracks,
		SortOrderStr:   strings.Join(currentSortFields, " → "),
		SortIndicators: sortIndicators,
	}

	// Execute the HTML template with the sorted data.
	if err := trackListTemplate.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func main() {
	// Initialize default sort order on server start
	currentSortFields = []string{"Title"}
	sortTracks(currentSortFields...)

	// Set up the HTTP route handler
	http.HandleFunc("/", trackList)

	// Start the HTTP server
	port := "8000"
	fmt.Printf("Server listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
