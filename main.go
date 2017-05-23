package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "USAGE:\n")
		fmt.Fprintf(os.Stderr, "\torient split data.csv\n")
		fmt.Fprintf(os.Stderr, "\torient race  data.csv\n")
		fmt.Fprintf(os.Stderr, "\torient delta data.csv\n")
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	comp, err := ParseCompetition(string(data))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	sort.Slice(comp.Runs, func(i, j int) bool {
		return comp.Runs[i].Result < comp.Runs[j].Result
	})

	tt := NewTimeTable()
	tt.AddCompetition(comp)
	tt.Sort()

	controls := []string{}
	controls = append(controls, "START")
	for _, run := range comp.Runs {
		if run.Status != Correct {
			continue
		}
		for _, c := range run.Controls {
			controls = append(controls, c.ID)
		}
		break
	}
	controls = append(controls, "FINISH")

	mode := os.Args[1]
	switch mode {
	case "split", "race", "delta":
		fmt.Print("Name;")
		fmt.Print("Course;")
		for _, cid := range controls {
			fmt.Print(cid, ";")
		}
		fmt.Println()
		for _, run := range comp.Runs {
			var times []Time
			switch mode {
			case "split":
				times = tt.Splits(run, controls)
			case "race":
				times = tt.Race(run, controls)
			case "delta":
				times = tt.Delta(run, controls)
			}

			fmt.Print(run.FirstName, " ", run.LastName, ";")
			fmt.Print(run.Course, ";")
			for _, t := range times {
				if !t.IsInvalid() {
					fmt.Print(t.Seconds())
				}
				fmt.Print(";")
			}
			fmt.Println()
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown mode %v\n", mode)
		os.Exit(1)
	}

}
