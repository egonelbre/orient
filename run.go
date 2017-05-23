package main

import (
	"errors"
	"fmt"
	"strings"
)

type Status string

const (
	Correct      = Status("C")
	Disqualified = Status("DQ")
	DidNotFinish = Status("DNF")
)

type Time int64 // time in seconds, HH:MM:SS

const (
	InvalidTime = Time(-1 << 63)
	TimeInf     = Time(1<<63 - 1)
)

func (a Time) IsInvalid() bool { return a == InvalidTime }

func (a Time) Min(b Time) Time {
	if a < b {
		return a
	}
	return b
}

func (a Time) Max(b Time) Time {
	if a > b {
		return a
	}
	return b
}

func MaxTime(xs ...Time) Time {
	if len(xs) == 0 {
		return InvalidTime
	}
	r := xs[0]
	for _, x := range xs {
		if r < x {
			r = x
		}
	}
	return r
}

func AverageTimeN(xs []Time, n int) Time {
	if len(xs) > n {
		xs = xs[:n]
	}
	if len(xs) == 0 {
		return InvalidTime
	}
	t := Time(0)
	for _, x := range xs {
		t += x
	}
	return t / Time(len(xs))
}

func (a Time) Seconds() int { return int(a) }
func (t *Time) Unmarshal(s string) error {
	var hh, mm, ss int64
	_, err := fmt.Sscanf(s, "%d:%d:%d", &hh, &mm, &ss)
	*t = Time(hh*60*60 + mm*60 + ss)

	if t.String() != s {
		panic(s + " != " + t.String())
	}
	return err
}

func (t Time) String() string {
	var x = int64(t)
	var hh, mm, ss int64
	x, ss = x/60, x%60
	x, mm = x/60, x%60
	hh = x
	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}

type Run struct {
	Number string
	SICard string

	FirstName string
	LastName  string

	Info   string // club, organization, school or other
	Course string

	Result Time // total time spent
	Status Status

	Point string

	Start    Time // absolute start time
	Controls []Control
	Finish   Time // absolute finish time
}

type ControlStatus string

const (
	UndefinedControl = ControlStatus("?")
	WrongControl     = ControlStatus("-")
	CorrectControl   = ControlStatus("+")
)

type Control struct {
	ID     string
	Status ControlStatus
	Time   Time // absolute time
}

func ParseRun(line string) (*Run, error) {
	tokens := strings.Split(line, ";")
	r := &Run{}
	if len(tokens) < 11 {
		return nil, errors.New("invalid line")
	}
	if tokens[len(tokens)-1] == "" {
		tokens = tokens[:len(tokens)-1]
	}

	r.Number = tokens[0]
	r.SICard = tokens[1]
	r.FirstName = tokens[2]
	r.LastName = tokens[3]
	r.Info = tokens[4]
	r.Course = tokens[5]
	if err := r.Result.Unmarshal(tokens[6]); err != nil {
		return nil, err
	}
	r.Status = Status(tokens[7])
	r.Point = tokens[8]
	if err := r.Start.Unmarshal(tokens[9]); err != nil {
		return nil, err
	}
	if err := r.Finish.Unmarshal(tokens[len(tokens)-1]); err != nil {
		return nil, err
	}

	n := len(tokens) - 11
	if n%3 != 0 {
		return nil, errors.New("invalid number of control data-points")
	}

	r.Controls = make([]Control, n/3)
	for i := range r.Controls {
		c := &r.Controls[i]
		c.ID = tokens[10+i*3]
		c.Status = ControlStatus(tokens[10+i*3+1])
		if err := c.Time.Unmarshal(tokens[10+i*3+2]); err != nil {
			return nil, err
		}
	}

	return r, nil
}
