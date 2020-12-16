package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func parseSchedule() (departureTime int, buses []int, err error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return 0, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			return 0, nil, err
		}

		if departureTime == 0 {
			departureTime, err = strconv.Atoi(scanner.Text())
			if err != nil {
				return 0, nil, err
			}
		} else {
			busesRaw := strings.Split(scanner.Text(), ",")
			for _, bus := range busesRaw {
				if bus == "x" {
					buses = append(buses, -1)
					continue
				}

				busInt, err := strconv.Atoi(bus)
				if err != nil {
					return 0, nil, err
				}

				buses = append(buses, busInt)
			}
		}
	}

	return departureTime, buses, nil
}

func part1(departureTime int, buses []int) int {
	closestDeparture := -1
	closestBus := -1
	for _, bus := range buses {
		if bus == -1 {
			continue
		}

		nextDeparture := departureTime + (bus - (departureTime % bus))

		if closestDeparture == -1 || nextDeparture-departureTime < closestDeparture {
			closestDeparture = nextDeparture - departureTime
			closestBus = bus
		}
	}

	return closestDeparture * closestBus
}

func part2(buses []int) int {
	departure := 0
	increment := 1
	sortedBuses := [][]int{}

	for pos, id := range buses {
		if id == -1 {
			continue
		}

		sortedBuses = append(sortedBuses, []int{pos, id})
	}

	sort.Slice(sortedBuses, func(i, j int) bool { return sortedBuses[i][1] > sortedBuses[j][1] })

	for _, bus := range sortedBuses {
		for (departure+bus[0])%bus[1] != 0 {
			departure += increment
		}
		increment *= bus[1]
	}

	return departure
}

func main() {
	departureTime, buses, err := parseSchedule()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error parsing schedule: %s", err))
	}

	log.Println(fmt.Sprintf("The answer to part 1 is %d", part1(departureTime, buses)))
	log.Println(fmt.Sprintf("The answer to part 2 is %d", part2(buses)))
}
