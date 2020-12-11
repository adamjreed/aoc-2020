package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
)

type Direction struct {
	Row  int
	Seat int
}

var seatMovement = []Direction{
	{Row: -1, Seat: -1},
	{Row: -1, Seat: 0},
	{Row: -1, Seat: 1},
	{Row: 0, Seat: -1},
	{Row: 0, Seat: 1},
	{Row: 1, Seat: -1},
	{Row: 1, Seat: 0},
	{Row: 1, Seat: 1},
}

type Rules struct {
	ContinueLooking bool
	SitOrStand      func(occupied int, occupiedAdjacent int) int
}

func getLayout() (numbers [][]int, err error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	layout := [][]int{}

	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			return nil, err
		}

		line := scanner.Text()
		seats := []int{}

		for _, c := range line {
			switch string(c) {
			case ".":
				seats = append(seats, -1)
			case "L":
				seats = append(seats, 0)
			case "#":
				seats = append(seats, 1)
			}
		}

		layout = append(layout, seats)
	}

	return layout, nil
}

func occupiedSeats(layout [][]int) int {
	occ := 0
	for _, seats := range layout {
		for _, seat := range seats {
			if seat > 0 {
				occ++
			}
		}
	}

	return occ
}

func iterateLayout(layout [][]int, seatingRules Rules) int {
	for {
		newLayout := [][]int{}
		for row, seats := range layout {
			newRow := []int{}
			for seat, _ := range seats {
				newRow = append(newRow, iterateRow(row, seat, layout, seatingRules))
			}
			newLayout = append(newLayout, newRow)
		}

		if reflect.DeepEqual(layout, newLayout) {
			break
		}

		layout = newLayout
	}

	return occupiedSeats(layout)
}

func iterateRow(row int, seat int, layout [][]int, seatingRules Rules) int {
	if layout[row][seat] == -1 {
		return -1
	}

	occupiedAdj := 0
out:
	for _, adj := range seatMovement {
		adjRow := row + adj.Row
		adjSeat := seat + adj.Seat

		for {
			if adjRow < 0 || adjRow >= len(layout) || adjSeat < 0 || adjSeat >= len(layout[row]) {
				continue out
			}

			if layout[adjRow][adjSeat] != -1 || !seatingRules.ContinueLooking {
				break
			}

			adjRow = adjRow + adj.Row
			adjSeat = adjSeat + adj.Seat
		}

		if layout[adjRow][adjSeat] > 0 {
			occupiedAdj++
		}
	}

	return seatingRules.SitOrStand(layout[row][seat], occupiedAdj)
}

func sitOrStandPart1(occupied int, occupiedAdj int) int {
	switch occupied {
	case 0:
		if occupiedAdj == 0 {
			return 1
		}
	case 1:
		if occupiedAdj >= 4 {
			return 0
		}
	}

	return occupied
}

func sitOrStandPart2(occupied int, occupiedAdj int) int {
	switch occupied {
	case 0:
		if occupiedAdj == 0 {
			return 1
		}
	case 1:
		if occupiedAdj >= 5 {
			return 0
		}
	}

	return occupied
}

func part1(layout [][]int) int {
	return iterateLayout(layout, Rules{
		ContinueLooking: false,
		SitOrStand:      sitOrStandPart1,
	})
}

func part2(layout [][]int) int {
	return iterateLayout(layout, Rules{
		ContinueLooking: true,
		SitOrStand:      sitOrStandPart2,
	})
}

func main() {
	layout, err := getLayout()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error fetching layout: %s", err))
	}

	log.Println(fmt.Sprintf("There were %d seats occupied in part 1", part1(layout)))
	log.Println(fmt.Sprintf("There were %d seats occupied in part 2", part2(layout)))
}
