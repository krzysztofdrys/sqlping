package pinger

import (
	"fmt"
	"time"
)

type State struct {
	Error     string
	StartedAt time.Time
	EndedAt   time.Time
}

func (s State) OnPing(err error) State {
	if s.Error == "" && err == nil {
		return s
	}
	now := time.Now().UTC()

	if s.Error != "" && err == nil {
		return State{StartedAt: now}
	}
	if s.Error == "" && err != nil {
		return State{StartedAt: now, Error: err.Error()}
	}
	if s.Error == err.Error() {
		return s
	}
	return State{StartedAt: s.StartedAt, Error: err.Error()}
}

func (s State) String() string {
	state := "ok"
	verb := "was"
	if s.Error != "" {
		state = "down"
	}
	var duration time.Duration
	if s.EndedAt.IsZero() {
		verb = "has been"
		duration = time.Since(s.StartedAt)
	} else {
		duration = s.EndedAt.Sub(s.StartedAt)
	}

	str := fmt.Sprintf("Connection %s %s for %s", verb, state, duration)
	if s.Error != "" {
		str = str + " Last error: " + s.Error
	}
	return str
}
