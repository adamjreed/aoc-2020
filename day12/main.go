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

type Coords struct {
	posX int
	posY int
}

type Direction struct {
	Instruction string
	Amount      int
}

type Navigator struct {
	Directions []*Direction
	Ship       *Coords
	NavMethod  NavMethod
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
			n.NavMethod.Move(direction, n.Ship)
		case "R":
			fallthrough
		case "L":
			n.NavMethod.Turn(direction, n.Ship)
		case "F":
			n.NavMethod.Advance(direction, n.Ship)
		default:
			return 0, errors.New(fmt.Sprintf("unexpected instruction: %s", direction.Instruction))
		}
	}

	return int(math.Abs(float64(n.Ship.posX))) + int(math.Abs(float64(n.Ship.posY))), nil
}

type NavMethod interface {
	Move(dir *Direction, ship *Coords)
	Turn(dir *Direction, ship *Coords)
	Advance(dir *Direction, ship *Coords)
}

type HeadingNav struct {
	Heading int
}

func (h *HeadingNav) Move(dir *Direction, ship *Coords) {
	moveObj(ship, dir)
}

func (h *HeadingNav) Turn(dir *Direction, ship *Coords) {
	switch dir.Instruction {
	case "R":
		h.Heading = (h.Heading + dir.Amount) % 360
	case "L":
		if h.Heading < dir.Amount {
			h.Heading = (h.Heading + 360) - dir.Amount
		} else {
			h.Heading -= dir.Amount
		}
	}
}

func (h *HeadingNav) Advance(dir *Direction, ship *Coords) {
	if h.Heading < 90 {
		ship.posY += dir.Amount
	} else if h.Heading < 180 {
		ship.posX += dir.Amount
	} else if h.Heading < 270 {
		ship.posY -= dir.Amount
	} else if h.Heading < 360 {
		ship.posX -= dir.Amount
	}
}

type WaypointNav struct {
	Waypoint *Coords
}

func (w *WaypointNav) Move(dir *Direction, ship *Coords) {
	moveObj(w.Waypoint, dir)
}

func (w *WaypointNav) Turn(dir *Direction, ship *Coords) {
	switch dir.Instruction {
	case "R":
		angle := degToRad(float64(dir.Amount))
		w.RotateWaypoint(angle)
	case "L":
		angle := degToRad(float64(dir.Amount))
		w.RotateWaypoint(angle * -1)
	}
}

func (w *WaypointNav) Advance(dir *Direction, ship *Coords) {
	ship.posY += w.Waypoint.posY * dir.Amount
	ship.posX += w.Waypoint.posX * dir.Amount
}

func (w *WaypointNav) RotateWaypoint(angle float64) {
	x := int(float64(w.Waypoint.posX)*math.Cos(angle)) + int(float64(w.Waypoint.posY)*math.Sin(angle))
	y := int(float64(w.Waypoint.posY)*math.Cos(angle)) - int(float64(w.Waypoint.posX)*math.Sin(angle))

	w.Waypoint.posX = x
	w.Waypoint.posY = y
}

func moveObj(obj *Coords, dir *Direction) {
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
		NavMethod: &HeadingNav{
			Heading: 90,
		},
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
		NavMethod: &WaypointNav{
			Waypoint: &Coords{
				posX: 10,
				posY: 1,
			},
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
