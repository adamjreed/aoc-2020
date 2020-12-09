package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

func getNumbers() (numbers []int, err error) {
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

		number, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, number)
	}

	return numbers, nil
}

func part1(numbers []int) (int, error) {
	buffer := numbers[0:25]
	for i := 25; i < len(numbers)-25; i++ {
		valid := checkAddends(numbers[i], buffer)
		if !valid {
			return numbers[i], nil
		}

		buffer = append(buffer[1:], numbers[i])
	}

	return 0, errors.New("no numbers failed the encryption check")
}

func checkAddends(sum int, addends []int) bool {
	for _, addend1 := range addends {
		difference := sum - addend1
		if difference == addend1 {
			continue
		}
		for _, addend2 := range addends {
			if difference == addend2 {
				return true
			}
		}
	}

	return false
}

func part2(target int, numbers []int) (int, error) {
	buffer := []int{}
	for i := 0; i < len(numbers); i++ {
		sum := numbers[i]

		for j := 0; sum < target && j < len(numbers); j++ {
			buffer = append(buffer, numbers[i+j])
			sum = calculateSum(buffer)
		}

		if sum == target {
			lowest := buffer[0]
			highest := buffer[0]
			for _, i := range buffer {
				if i < lowest {
					lowest = i
				}

				if i > highest {
					highest = i
				}
			}

			return lowest + highest, nil
		}

		buffer = []int{}
	}

	return 0, errors.New(fmt.Sprintf("got to the end of input with no valid sum for %d", target))
}

func calculateSum(numbers []int) (sum int) {
	for _, number := range numbers {
		sum = sum + number
	}

	return sum
}

func main() {
	numbers, err := getNumbers()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error fetching instructions: %s", err))
	}

	notSum, err := part1(numbers)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error finding the weakness in part 1: %s", err))
	}

	log.Println(fmt.Sprintf("The first number which fails the encryption check is %d in part 1", notSum))

	weakness, err := part2(notSum, numbers)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error finding the weakness in part 2: %s", err))
	}
	log.Println(fmt.Sprintf("The encryption weakness is %d in part 2", weakness))
}
