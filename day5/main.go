package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
)

var (
	highestSeatId1 int
	seatIds        map[int]struct{}
	mySeatId       int
)

func parseSeats() error {
	f, err := os.Open("input.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			return err
		}

		line := scanner.Text()
		rows := line[0:7]
		row, err := findSeat(rows, 0, 127)
		if err != nil {
			return err
		}

		cols := line[7:]
		col, err := findSeat(cols, 0, 7)
		if err != nil {
			return err
		}

		seatId := calculateSeatId(row, col)
		if seatId > highestSeatId1 {
			highestSeatId1 = seatId
		}

		seatIds[seatId] = struct{}{}
	}

	return nil
}

func calculateSeatId(row int, col int) int {
	return row*8 + col
}

func findSeat(rows string, low int, high int) (int, error) {
	section := rows[0:1]
	diff := high - low

	switch section {
	case "F":
		fallthrough
	case "L":
		if len(rows) == 1 {
			return low, nil
		}
		return findSeat(rows[1:], low, high-(diff/2)-1)
	case "B":
		fallthrough
	case "R":
		if len(rows) == 1 {
			return high, nil
		}
		return findSeat(rows[1:], low+(diff/2)+1, high)
	}

	return -1, errors.New("invalid section")
}

func main() {
	seatIds = map[int]struct{}{}
	err := parseSeats()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error finding seats: %s", err))
	}

	for left, _ := range seatIds {
		_, middle := seatIds[left+1]
		_, right := seatIds[left+2]

		if middle == false && right == true {
			mySeatId = left + 1
		}
	}

	fmt.Println(fmt.Sprintf("The highest seat ID is: %d", highestSeatId1))
	fmt.Println(fmt.Sprintf("My seat ID is: %d", mySeatId))
}
