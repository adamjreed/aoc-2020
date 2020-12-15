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
	"strings"
)

type MaskSet struct {
	Mask     string
	Registry []Register
}

type Register struct {
	Address int
	Value   int
}

func reverse(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}

func parseInput() (maskSets []MaskSet, err error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	maskSets = []MaskSet{}

	var currentMask *MaskSet
	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			return nil, err
		}

		line := scanner.Text()
		if line[0:4] == "mask" {
			if currentMask != nil {
				maskSets = append(maskSets, *currentMask)
				currentMask = nil
			}

			currentMask = &MaskSet{
				Mask: reverse(line[7:]),
			}
		} else {
			re := regexp.MustCompile(`^mem\[(\d+)\] = (\d+)$`)
			matches := re.FindStringSubmatch(line)
			if len(matches) != 3 {
				return nil, errors.New("invalid input")
			}

			address, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, err
			}

			value, err := strconv.Atoi(matches[2])
			if err != nil {
				return nil, err
			}

			currentMask.Registry = append(currentMask.Registry, Register{
				Address: address,
				Value:   value,
			})
		}
	}

	if currentMask != nil {
		maskSets = append(maskSets, *currentMask)
	}

	return maskSets, nil
}

func calculateMasks(mask string) (masks []string) {
	re := regexp.MustCompile(`X`)

	matches := re.FindAllString(mask, -1)
	format := fmt.Sprintf("%%0%db", len(matches))
	for i := 0; i < int(math.Pow(float64(2), float64(len(matches)))); i++ {
		bin := reverse(fmt.Sprintf(format, i))
		newMask := mask
		newMask = strings.Replace(newMask, "0", "Y", -1)
		for _, char := range bin {
			newMask = strings.Replace(newMask, "X", string(char), 1)
		}

		newMask = strings.Replace(newMask, "Y", "X", -1)
		masks = append(masks, newMask)
	}

	return masks
}

func applyMask(mask string, value int) int {
	for pos, bit := range mask {
		switch string(bit) {
		case "0":
			value = clearBit(pos, value)
		case "1":
			value = setBit(pos, value)
		}
	}

	return value
}

func sumRegister(register map[int]int) (sum int) {
	for _, value := range register {
		sum += value
	}

	return sum
}

func setBit(pos int, value int) int {
	value |= (1 << pos)
	return value
}

func clearBit(pos int, value int) int {
	mask := ^(1 << pos)
	value &= mask
	return value
}

func part1(maskSets []MaskSet) int {
	register := map[int]int{}

	for _, set := range maskSets {
		for _, pos := range set.Registry {
			register[pos.Address] = applyMask(set.Mask, pos.Value)
		}
	}

	return sumRegister(register)
}

func part2(maskSets []MaskSet) int {
	register := map[int]int{}

	for _, set := range maskSets {
		for _, pos := range set.Registry {
			masks := calculateMasks(set.Mask)

			for _, mask := range masks {
				register[applyMask(mask, pos.Address)] = pos.Value
			}
		}
	}

	return sumRegister(register)
}

func main() {
	maskSets, err := parseInput()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error parsing input: %s", err))
	}

	log.Println(fmt.Sprintf("The sum in part 1 is %d", part1(maskSets)))
	log.Println(fmt.Sprintf("The sum in part 2 is %d", part2(maskSets)))
}
