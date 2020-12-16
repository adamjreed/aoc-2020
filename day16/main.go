package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type TicketInput struct {
	Rules         []*Rule
	YourTicket    *Ticket
	NearbyTickets []*Ticket
}

type Rule struct {
	Name   string
	Ranges [][]int
}

type Ticket struct {
	Fields []int
}

var invalidTickets map[int]struct{}

func getTicketInput() (*TicketInput, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var (
		rules         []*Rule
		yourTicket    *Ticket
		nearbyTickets []*Ticket
	)

	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			return nil, err
		}

		line := scanner.Text()

		if line == "" {
			continue
		}

		if line == "your ticket:" {
			yourTicket = &Ticket{}
			continue
		}

		if line == "nearby tickets:" {
			nearbyTickets = []*Ticket{}
			continue
		}

		if nearbyTickets != nil {
			ticket, err := parseTicket(line)
			if err != nil {
				return nil, err
			}
			nearbyTickets = append(nearbyTickets, ticket)
		} else if yourTicket != nil {
			yourTicket, err = parseTicket(line)
			if err != nil {
				return nil, err
			}
		} else {
			parts := strings.Split(line, ": ")
			rule := &Rule{
				Name: parts[0],
			}
			ranges := strings.Split(parts[1], " or ")
			for _, r := range ranges {
				bounds := strings.Split(r, "-")
				low, err := strconv.Atoi(bounds[0])
				if err != nil {
					return nil, err
				}
				high, err := strconv.Atoi(bounds[1])
				if err != nil {
					return nil, err
				}

				rule.Ranges = append(rule.Ranges, []int{low, high})
			}
			rules = append(rules, rule)
		}
	}

	return &TicketInput{
		Rules:         rules,
		YourTicket:    yourTicket,
		NearbyTickets: nearbyTickets,
	}, nil
}

func parseTicket(input string) (*Ticket, error) {
	fields := strings.Split(input, ",")
	ticket := &Ticket{}
	for _, field := range fields {
		fieldInt, err := strconv.Atoi(field)
		if err != nil {
			return nil, err
		}
		ticket.Fields = append(ticket.Fields, fieldInt)
	}

	return ticket, nil
}

func getRule(rule map[string]struct{}) string {
	for k, _ := range rule {
		return k
	}
	return ""
}

func reduceRules(ruleCandidates map[int]map[string]struct{}) map[string]int {
	assignedRules := map[string]int{}

	for len(ruleCandidates) > 0 {
		for field, rules := range ruleCandidates {
			if len(rules) == 1 {
				assignedRules[getRule(rules)] = field
				delete(ruleCandidates, field)
				break
			}

			for rule, _ := range rules {
				if _, ok := assignedRules[rule]; ok {
					delete(rules, rule)
				}
			}
		}
	}

	return assignedRules
}

func checkRules(rules1 map[string]struct{}, rules2 map[string]struct{}) map[string]struct{} {
	if len(rules1) == 0 {
		return rules2
	}

	possibleRules := map[string]struct{}{}

	for rule, _ := range rules2 {
		if _, ok := rules1[rule]; ok {
			possibleRules[rule] = struct{}{}
		}
	}

	return possibleRules
}

func part1(input *TicketInput) int {
	var errorRate int
	invalidTickets = map[int]struct{}{}

	for i, ticket := range input.NearbyTickets {
		for _, field := range ticket.Fields {
			valid := false
			for _, rule := range input.Rules {
				for _, r := range rule.Ranges {
					if field >= r[0] && field <= r[1] {
						valid = true
					}
				}
			}
			if !valid {
				invalidTickets[i] = struct{}{}
				errorRate += field
			}
		}
	}

	return errorRate
}

func part2(input *TicketInput) int {
	allTickets := []*Ticket{}
	allTickets = append(allTickets, input.NearbyTickets...)
	allTickets = append(allTickets, input.YourTicket)

	ruleCandidates := map[int]map[string]struct{}{}
	for i, ticket := range input.NearbyTickets {
		if _, ok := invalidTickets[i]; ok {
			continue
		}
		for i, field := range ticket.Fields {
			validRules := map[string]struct{}{}
			for _, rule := range input.Rules {
				for _, r := range rule.Ranges {
					if field >= r[0] && field <= r[1] {
						validRules[rule.Name] = struct{}{}
					}
				}
			}
			ruleCandidates[i] = checkRules(ruleCandidates[i], validRules)
		}
	}

	assignedRules := reduceRules(ruleCandidates)

	departureFieldsProduct := 1
	for rule, i := range assignedRules {
		if strings.HasPrefix(rule, "departure ") {
			departureFieldsProduct *= input.YourTicket.Fields[i]
		}
	}

	return departureFieldsProduct
}

func main() {
	ticketInput, err := getTicketInput()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error getting starting numbers: %s", err))
	}

	log.Println(fmt.Sprintf("The error rate in part 1 is %d", part1(ticketInput)))
	log.Println(fmt.Sprintf("The product of departure fields in part 2 is %d", part2(ticketInput)))
}
