package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	validPassports1 int
	validPassports2 int
)

type Passport struct {
	Ecl string
	Pid string
	Eyr string
	Hcl string
	Byr string
	Iyr string
	Cid string
	Hgt string
}

func (p *Passport) IsValid(strict bool) bool {
	if p.Ecl == "" || p.Pid == "" || p.Eyr == "" || p.Hcl == "" || p.Byr == "" || p.Iyr == "" || p.Hgt == "" {
		return false
	}

	if strict {
		if !p.validateBirthday() {
			log.Println(fmt.Sprintf("Birthday invalid: %s", p.Byr))
			return false
		}
		if !p.validateIssued() {
			log.Println(fmt.Sprintf("Issue Year invalid: %s", p.Iyr))
			return false
		}
		if !p.validateExpires() {
			log.Println(fmt.Sprintf("Expires Year invalid: %s", p.Eyr))
			return false
		}
		if !p.validateHeight() {
			log.Println(fmt.Sprintf("Height invalid: %s", p.Hgt))
			return false
		}
		if !p.validateHairColor() {
			log.Println(fmt.Sprintf("Hair Color invalid: %s", p.Hcl))
			return false
		}
		if !p.validateEyeColor() {
			log.Println(fmt.Sprintf("Eye Color invalid: %s", p.Hcl))
			return false
		}
		if !p.validatePassportId() {
			log.Println(fmt.Sprintf("Passport ID invalid: %s", p.Pid))
			return false
		}
	}

	return true
}

func (p *Passport) validateBirthday() bool {
	return validateYear(p.Byr, 1920, 2002)
}

func (p *Passport) validateIssued() bool {
	return validateYear(p.Iyr, 2010, 2020)
}

func (p *Passport) validateExpires() bool {
	return validateYear(p.Eyr, 2020, 2030)
}

func (p *Passport) validateHeight() bool {
	return validateHeight(p.Hgt, "cm", 150, 193) != validateHeight(p.Hgt, "in", 59, 76)
}

func (p *Passport) validateHairColor() bool {
	re := regexp.MustCompile(`^#[0-9a-f]{6}$`)
	return re.MatchString(p.Hcl)
}

var validEyeColors = map[string]struct{}{
	"amb": {},
	"blu": {},
	"brn": {},
	"gry": {},
	"grn": {},
	"hzl": {},
	"oth": {},
}

func (p *Passport) validateEyeColor() bool {
	_, ok := validEyeColors[p.Ecl]
	return ok
}

func (p *Passport) validatePassportId() bool {
	re := regexp.MustCompile(`^[0-9]{9}$`)
	return re.MatchString(p.Pid)
}

func validateYear(year string, low int, high int) bool {
	compare, err := strconv.Atoi(year)
	if err != nil {
		return false
	}

	if compare < low || compare > high {
		return false
	}

	return true
}

func validateHeight(input string, unit string, low int, high int) bool {
	if strings.Contains(input, unit) {
		height, err := strconv.Atoi(strings.Replace(input, unit, "", -1))
		if err != nil {
			return false
		}
		if height < low || height > high {
			return false
		}

		return true
	}

	return false
}

func parsePassportData(info []string) (*Passport, error) {
	pairs := map[string]string{}
	for _, line := range info {
		groups := strings.Split(line, " ")
		for _, group := range groups {
			split := strings.Split(group, ":")
			pairs[split[0]] = split[1]
		}
	}

	asJSON, err := json.Marshal(pairs)
	if err != nil {
		return nil, err
	}

	var passport *Passport
	err = json.Unmarshal(asJSON, &passport)
	if err != nil {
		return nil, err
	}

	return passport, nil
}

func checkPassports() error {
	f, err := os.Open("input.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var buffer []string
	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			return err
		}

		line := scanner.Text()

		if len(line) == 0 {
			err = checkPassport(buffer)
			if err != nil {
				return err
			}
			buffer = []string{}
			continue
		}

		buffer = append(buffer, scanner.Text())
	}

	if len(buffer) > 0 {
		err = checkPassport(buffer)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkPassport(passportInfo []string) error {
	passport, err := parsePassportData(passportInfo)
	if err != nil {
		return err
	}
	if passport.IsValid(false) {
		validPassports1++
	}
	if passport.IsValid(true) {
		validPassports2++
	}

	return nil
}

func main() {
	err := checkPassports()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error parsing passports: %s", err))
	}

	log.Println(fmt.Sprintf("There are %d valid passports in part 1", validPassports1))
	log.Println(fmt.Sprintf("There are %d valid passports in part 2", validPassports2))
}
