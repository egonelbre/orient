package main

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

type Competition struct {
	Runs []*Run
}

func ParseCompetition(data string) (*Competition, error) {
	lineno := 0
	sc := bufio.NewScanner(strings.NewReader(data))

	errs := []error{}
	fail := func(err error) bool {
		if err != nil {
			errs = append(errs, fmt.Errorf("%d: %v", lineno, err))
		}
		return false
	}

	comp := &Competition{}
	for sc.Scan() {
		lineno++
		line := sc.Text()
		if line == "" {
			continue
		}

		r, err := ParseRun(line)
		if fail(err) {
			continue
		}

		comp.Runs = append(comp.Runs, r)
	}

	if len(errs) == 0 {
		return comp, nil
	} else {
		text := errs[0].Error()
		for _, err := range errs[1:] {
			text += "\n" + err.Error()
		}
		return comp, errors.New(text)
	}
}
