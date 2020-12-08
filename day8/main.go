package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Instruction struct {
	Code   string
	Amount int
}

func getInstructions() (instructions []*Instruction, err error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			return nil, err
		}

		line := scanner.Text()
		re := regexp.MustCompile(`(\w{3}) ([+-]\d+)`)
		matches := re.FindStringSubmatch(line)
		amount, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, &Instruction{
			Code:   matches[1],
			Amount: amount,
		})
	}

	return instructions, nil
}

func part1(instructions []*Instruction) (acc int) {
	acc, _ = run(instructions)
	return acc
}

func part2(instructions []*Instruction) (acc int, err error) {
	for _, instruction := range instructions {
		if instruction.Code == "acc" {
			continue
		}

		flipInstruction(instruction)
		acc, finished := run(instructions)
		if finished {
			return acc, nil
		} else {
			flipInstruction(instruction)
		}
	}

	return 0, errors.New("flipped all possible instructions and program still failed")
}

func flipInstruction(instruction *Instruction) {
	switch instruction.Code {
	case "nop":
		instruction.Code = "jmp"
	case "jmp":
		instruction.Code = "nop"
	}
}

func run(instructions []*Instruction) (acc int, finished bool) {
	i := 0
	called := map[int]int{}
	for i < len(instructions) {
		if called[i] > 0 {
			return acc, false
		}
		called[i]++

		switch instructions[i].Code {
		case "acc":
			acc = acc + instructions[i].Amount
			i++
		case "jmp":
			i = i + instructions[i].Amount
		case "nop":
			i++
		}
	}

	return acc, true
}

func main() {
	instructions, err := getInstructions()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error fetching instructions: %s", err))
	}

	log.Println(fmt.Sprintf("The accumulator has a value of %d in part 1", part1(instructions)))
	acc, err := part2(instructions)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error running part 2: %s", err))
	}
	log.Println(fmt.Sprintf("The accumulator has a value of %d in part 2", acc))
}
