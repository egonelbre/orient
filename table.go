package main

import "sort"

type TimeTable struct {
	Start   Time
	Control map[string]*ControlTime
}

type ControlTime struct {
	SplitFrom map[string][]Time
	Relative  map[string][]Time  // runner -> relative time from start
	Split     map[[2]string]Time // split time {from, runner}
}

func NewTimeTable() *TimeTable {
	return &TimeTable{
		Start:   TimeInf,
		Control: make(map[string]*ControlTime),
	}
}

func (tt *TimeTable) AddRun(r *Run) {
	start := r.Start
	if start < tt.Start {
		tt.Start = start
	}

	last, lasttime := "", Time(0)
	add := func(check string, time Time) {
		c, ok := tt.Control[check]
		if !ok {
			c = &ControlTime{
				SplitFrom: make(map[string][]Time),
				Relative:  make(map[string][]Time),
				Split:     make(map[[2]string]Time),
			}
			tt.Control[check] = c
		}

		if last != "" {
			splittime := time - lasttime
			c.SplitFrom[last] = append(c.SplitFrom[last], splittime)
			c.Split[[2]string{last, r.SICard}] = splittime
		}
		last, lasttime = check, time

		c.Relative[r.SICard] = append(c.Relative[r.SICard], time-start)
	}

	add("START", r.Start)
	for _, c := range r.Controls {
		add(c.ID, c.Time)
	}
	add("FINISH", r.Finish)
}

func (tt *TimeTable) AddCompetition(c *Competition) {
	for _, r := range c.Runs {
		tt.AddRun(r)
	}
}

func (tt *TimeTable) Times(sicard string, controls []string) []Time {
	times := []Time{}
	for _, cid := range controls {
		c, ok := tt.Control[cid]
		if !ok {
			times = append(times, InvalidTime)
			continue
		}

		xs, ok := c.Relative[sicard]
		if !ok || len(xs) == 0 {
			times = append(times, InvalidTime)
			continue
		}

		x := xs[0]
		for _, b := range xs[1:] {
			x = x.Max(b)
		}

		times = append(times, x)
	}

	return times
}

func (tt *TimeTable) Splits(run *Run, controls []string) []Time {
	times := []Time{}
	for _, cid := range controls {
		c, ok := tt.Control[cid]
		if !ok {
			times = append(times, InvalidTime)
			continue
		}

		times = append(times, MaxTime(c.Relative[run.SICard]...))
	}

	return times
}

func (tt *TimeTable) Race(run *Run, controls []string) []Time {
	times := []Time{}

	for _, cid := range controls {
		c, ok := tt.Control[cid]
		if !ok {
			times = append(times, InvalidTime)
			continue
		}

		x := MaxTime(c.Relative[run.SICard]...)
		if x == InvalidTime {
			times = append(times, InvalidTime)
		} else {
			times = append(times, run.Start+x-tt.Start)
		}
	}

	return times
}

func (tt *TimeTable) Delta(run *Run, controls []string) []Time {
	times := []Time{}

	prev := ""
	for _, cid := range controls {
		c, ok := tt.Control[cid]
		if prev == "" {
			prev = cid
			times = append(times, InvalidTime)
			continue
		}

		split, ok := c.Split[[2]string{prev, run.SICard}]
		if !ok {
			times = append(times, InvalidTime)
			prev = cid
			continue
		}

		target := AverageTimeN(c.SplitFrom[prev], 6)
		times = append(times, target-split)
		prev = cid
	}

	return times
}

func (tt *TimeTable) Sort() {
	for _, c := range tt.Control {
		for from, xs := range c.SplitFrom {
			sort.Slice(xs, func(i, j int) bool {
				return xs[i] < xs[j]
			})
			c.SplitFrom[from] = xs
		}
	}
}
