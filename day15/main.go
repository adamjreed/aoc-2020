package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func getStartingNumbers() (nums []int, err error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	nums = []int{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			return nil, err
		}

		split := strings.Split(scanner.Text(), ",")
		for _, char := range split {
			num, err := strconv.Atoi(char)
			if err != nil {
				return nil, err
			}

			nums = append(nums, num)
		}
	}

	return nums, nil
}

func iterate(nums []int, limit int) int {
	lastSpoken := -1
	numsSpoken := map[int]int{}
	for i := 1; i <= limit; i++ {
		if i <= len(nums) {
			lastSpoken = nums[i-1]
			numsSpoken[nums[i-1]] = i
			continue
		}

		if turnSpoken, ok := numsSpoken[lastSpoken]; ok {
			if turnSpoken == i-1 {
				lastSpoken = 0
				continue
			}

			numsSpoken[lastSpoken] = i - 1
			lastSpoken = i - 1 - turnSpoken
		} else {
			numsSpoken[lastSpoken] = i - 1
			lastSpoken = 0
		}
	}

	return lastSpoken
}

func main() {
	nums, err := getStartingNumbers()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error getting starting numbers: %s", err))
	}

	log.Println(fmt.Sprintf("The 2020th number in part 1 is %d", iterate(nums, 2020)))
	log.Println(fmt.Sprintf("The 30000000th number in part 2 is %d", iterate(nums, 30000000)))
}
