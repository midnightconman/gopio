package schema

import (
	"fmt"
	"strings"
)

type Direction int32
type State int32
type Pull int32

const (
	Input Direction = iota
	Output
)
const _Direction_name = "InputOutput"

var _Direction_index = [...]int32{0, 5, 11}

func (i Direction) String() string {
	if i >= Direction(len(_Direction_index)-1) {
		return fmt.Sprintf("Direction(%d)", i)
	}
	return _Direction_name[_Direction_index[i]:_Direction_index[i+1]]
}

func ParseDirection(direction string) Direction {
	switch strings.ToLower(direction) {
	case "input":
		return Input
	default:
		return Output
	}
}

const (
	Low State = iota
	High
)
const _State_name = "LowHigh"

var _State_index = [...]int32{0, 3, 7}

func (i State) String() string {
	if i >= State(len(_State_index)-1) {
		return fmt.Sprintf("State(%d)", i)
	}
	return _State_name[_State_index[i]:_State_index[i+1]]
}

func ParseState(state string) State {
	switch strings.ToLower(state) {
	case "high":
		return High
	default:
		return Low
	}
}

const (
	PullOff Pull = iota
	PullDown
	PullUp
)
const _Pull_name = "PullOffPullDownPullUp"

var _Pull_index = [...]int32{0, 7, 15, 21}

func (i Pull) String() string {
	if i >= Pull(len(_Pull_index)-1) {
		return fmt.Sprintf("Pull(%d)", i)
	}
	return _Pull_name[_Pull_index[i]:_Pull_index[i+1]]
}

func ParsePull(pull string) Pull {
	switch strings.ToLower(pull) {
	case "pullup":
		return PullUp
	case "pulldown":
		return PullDown
	default:
		return PullOff
	}
}
