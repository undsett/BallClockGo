package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	oneMinute    Track
	fiveMinute   Track
	oneHour      Track
	mainTrack    Track
	defaultValue Stack
	minutes      int
	cycles       = 1
	running      = true
	start        time.Time
	duration     = time.Duration(0)
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter 'mode1' or 'mode2' for mode selection: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimRight(text, "\n")
	if text != "mode1" && text != "mode2" {
		fmt.Printf("Invalid mode input")
		return
	}

	fmt.Print("Enter the number of balls. Valid numbers of balls are in the range 27 to 127 : ")
	number, _ := reader.ReadString('\n')
	number = strings.TrimRight(number, "\n")
	balls, err := strconv.Atoi(number)
	if err != nil && balls < 27 && balls > 127 {
		fmt.Println("Valid numbers of balls are in the range 27 to 127")
		return
	}

	switch text {
	case "mode1":
		newBallClock(balls, 0)
		runClock()
	case "mode2":
		fmt.Print("Enter the number of minutes : ")
		number, _ := reader.ReadString('\n')
		number = strings.TrimRight(number, "\n")
		minutes, err := strconv.Atoi(number)
		if err != nil && minutes > 1{
			fmt.Println("Number of minutes must be higher than 1")
			return
		}
		newBallClock(balls, minutes)
		runClock()
	}
}

func newBallClock(balls int, min int) {
	minutes = min
	mainTrack = NewTrack(balls, nil)
	oneHour = NewTrack(11, &mainTrack)
	fiveMinute = NewTrack(11, &oneHour)
	oneMinute = NewTrack(4, &fiveMinute)
	mainTrack.nextTrack = &oneMinute
	defaultValue = mainTrack.fillTrack(balls)
}

func runClock() {
	start = time.Now()
	defer printTime()
	if minutes == 0 {
		for running {
			mainTrack.move(&mainTrack)
			cycles++
			}
	} else {
		for cycles < minutes +1 {
			mainTrack.move(&mainTrack)
			cycles++
		}
		duration = time.Since(start)
		printJson()
	}
}

type Track struct {
	rail Stack
	capacity int
	nextTrack *Track
}

func NewTrack(capacity int, nextTrack *Track) Track {
	return Track{capacity: capacity, nextTrack: nextTrack}
}

func (tr *Track) isFull() bool {
	return len(tr.rail) == tr.capacity
}

type Mover interface {
	move(*Track)
}

func (tr *Track) move(mainTrack *Track) {
	pushOrRelease(mainTrack.rail.Pull(), tr, mainTrack)
}

func pushOrRelease(ball uint8, track *Track, mainTrack *Track) {
	if !track.nextTrack.isFull() {
		track.nextTrack.rail.Push(ball)
	} else {
		releaseTrack(ball, track.nextTrack, mainTrack)
	}
}

func releaseTrack(ball uint8, track *Track, mainTrack *Track) {
Loop:
	for {
		if value, ok := track.rail.Pop(); ok {
			mainTrack.rail.Push(value)
		}else {
			break Loop
		}
	}
	pushOrRelease(ball, track, mainTrack)
		if track.nextTrack == mainTrack {
			if Equal(mainTrack.rail, defaultValue) {
				fmt.Printf("%v  balls cycle after %v  days.\n", mainTrack.capacity, cycles/1440)
				duration = time.Since(start)
				running = false
			}
		}


}

func (tr *Track) fillTrack(capacity int) Stack{
	var i uint8 = 1
	for ; i < uint8(capacity + 1); i++ {
		tr.rail.Push(i)
	}
	return tr.rail
}

func Equal(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func printJson()  {
	response := make(map[string][]string)
	response["Min"] = intToString(oneMinute.rail)
	response["FiveMin"] = intToString(fiveMinute.rail)
	response["Hour"] = intToString(oneHour.rail)
	response["Main"] = intToString(mainTrack.rail)
	jsonResp, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(jsonResp))
}

func intToString (stack Stack) []string {
	var result []string
	for _, v := range stack {
		result = append(result, strconv.Itoa(int(v)))
	}
	return result
}
func printTime() {
	fmt.Println("Completed in", duration)
}

type Stack []uint8

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(element uint8) {
	*s = append(*s, element)
}

func (s *Stack) Pop() (uint8, bool) {
	if s.IsEmpty() {
		return 0, false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}

func (s *Stack) Pull() uint8 {
	if s.IsEmpty() {
		return 0
	} else {
		element := (*s)[0]
		*s = (*s)[1:]
		return element
	}
}