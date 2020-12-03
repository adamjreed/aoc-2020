package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type PasswordParts struct {
	Low int
	High int
	Target string
	Password string
}

func parse(raw string) (*PasswordParts, error) {
	re := regexp.MustCompile(`(\d+)-(\d+) ([a-z]): (\w+)`)
	groups := re.FindStringSubmatch(raw)

	if len(groups) < 5 {
		return nil, errors.New("not enough matches in input row")
	}

	low, err := strconv.Atoi(groups[1])
	if err != nil {
		return nil, err
	}
	high, err := strconv.Atoi(groups[2])
	if err != nil {
		return nil, err
	}

	return &PasswordParts{
		Low:      low,
		High:     high,
		Target:   groups[3],
		Password: groups[4],
	}, nil
}

func part1(parts *PasswordParts) bool {
	count := strings.Count(parts.Password, parts.Target)

	if count >= parts.Low && count <= parts.High {
		return true
	}

	return false
}

func part2(parts *PasswordParts) bool {
	matchFirstPosition := parts.Password[parts.Low-1:parts.Low] == parts.Target
	matchSecondPosition := parts.Password[parts.High-1:parts.High] == parts.Target

	if matchFirstPosition != matchSecondPosition {
		return true
	}

	return false
}

func main() {
	var (
		valid1 int
		valid2 int
	)

	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln(fmt.Sprintf("Could not load input: %s", err))
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			log.Fatalln(fmt.Sprintf("Error scanning input: %s", err))
		}

		row, err := parse(scanner.Text())
		if err != nil {
			log.Fatalln(fmt.Sprintf("Could not parse input row: %s", scanner.Text()))
		}

		if part1(row) {
			valid1++
		}

		if part2(row) {
			valid2++
		}
	}

	log.Println(fmt.Sprintf("There were %d valid passwords for part 1", valid1))
	log.Println(fmt.Sprintf("There were %d valid passwords for part 2", valid2))
}