package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

type Navigator struct {
	Directions []*Direction
	Ship       *Coords
	Waypoint   *Coords
	Heading    int
}

type Coords struct {
	posX int
	posY int
}

type Direction struct {
	Instruction string
	Amount      int
}

func (n *Navigator) Run() (int, error) {
	for _, direction := range n.Directions {
		switch direction.Instruction {
		case "N":
			fallthrough
		case "E":
			fallthrough
		case "S":
			fallthrough
		case "W":
			n.Move(direction)
		case "R":
			fallthrough
		case "L":
			n.Turn(direction)
		case "F":
			n.Advance(direction)
		default:
			return 0, errors.New(fmt.Sprintf("unexpected instruction: %s", direction.Instruction))
		}
	}

	return int(math.Abs(float64(n.Ship.posX))) + int(math.Abs(float64(n.Ship.posY))), nil
}

func (n *Navigator) Move(dir *Direction) {
	var obj *Coords
	if n.Waypoint != nil {
		obj = n.Waypoint
	} else {
		obj = n.Ship
	}

	switch dir.Instruction {
	case "N":
		obj.posY += dir.Amount
	case "E":
		obj.posX += dir.Amount
	case "S":
		obj.posY -= dir.Amount
	case "W":
		obj.posX -= dir.Amount
	}
}

func (n *Navigator) Advance(dir *Direction) {
	if n.Waypoint != nil {
		n.Ship.posY += n.Waypoint.posY * dir.Amount
		n.Ship.posX += n.Waypoint.posX * dir.Amount
	} else {
		if n.Heading < 90 {
			n.Ship.posY += dir.Amount
		} else if n.Heading < 180 {
			n.Ship.posX += dir.Amount
		} else if n.Heading < 270 {
			n.Ship.posY -= dir.Amount
		} else if n.Heading < 360 {
			n.Ship.posX -= dir.Amount
		}
	}
}

func (n *Navigator) Turn(dir *Direction) {
	switch dir.Instruction {
	case "R":
		if n.Waypoint != nil {
			angle := degToRad(float64(dir.Amount))
			n.RotateWaypoint(angle)
		} else {
			n.Heading = (n.Heading + dir.Amount) % 360
		}
	case "L":
		if n.Waypoint != nil {
			angle := degToRad(float64(dir.Amount))
			n.RotateWaypoint(angle * -1)
		} else {
			if n.Heading < dir.Amount {
				n.Heading = (n.Heading + 360) - dir.Amount
			} else {
				n.Heading -= dir.Amount
			}
		}
	}
}

func (n *Navigator) RotateWaypoint(angle float64) {
	x := int(float64(n.Waypoint.posX)*math.Cos(angle)) + int(float64(n.Waypoint.posY)*math.Sin(angle))
	y := int(float64(n.Waypoint.posY)*math.Cos(angle)) - int(float64(n.Waypoint.posX)*math.Sin(angle))

	n.Waypoint.posX = x
	n.Waypoint.posY = y
}

func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func getDirections() ([]*Direction, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var directions []*Direction

	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			return nil, err
		}

		dir, err := getDirection(scanner.Text())
		if err != nil {
			return nil, err
		}

		directions = append(directions, dir)
	}

	return directions, nil
}

func getDirection(input string) (*Direction, error) {
	re := regexp.MustCompile(`^(\w{1})(\d+)$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 3 {
		return nil, errors.New(fmt.Sprintf("unexpected format for line: %s", input))
	}

	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, err
	}

	return &Direction{
		Instruction: matches[1],
		Amount:      amount,
	}, nil
}

func part1(directions []*Direction) (int, error) {
	nav := &Navigator{
		Directions: directions,
		Ship: &Coords{
			posX: 0,
			posY: 0,
		},
		Heading: 90,
	}

	mDistance, err := nav.Run()
	if err != nil {
		return 0, err
	}

	return mDistance, nil
}

func part2(directions []*Direction) (int, error) {
	nav := &Navigator{
		Directions: directions,
		Ship: &Coords{
			posX: 0,
			posY: 0,
		},
		Waypoint: &Coords{
			posX: 10,
			posY: 1,
		},
	}

	mDistance, err := nav.Run()
	if err != nil {
		return 0, err
	}

	return mDistance, nil
}

func main() {
	directions, err := getDirections()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error fetching layout: %s", err))
	}

	mDistance, err := part1(directions)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error calculating Manhattan distance for part 1: %s", err))
	}

	log.Println(fmt.Sprintf("The Manhattan distance is %d in part 1", mDistance))

	mDistance, err = part2(directions)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error calculating Manhattan distance for part 2: %s", err))
	}

	log.Println(fmt.Sprintf("The Manhattan distance is %d in part 2", mDistance))
}
