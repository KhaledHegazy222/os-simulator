package parser

import (
	"os"
	"strings"
)

type Instruction struct {
	Command string
	Args    []string
}

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}
func (p *Parser) Parse(file string) ([]Instruction, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return []Instruction{}, err
	}
	lines := strings.Split(string(bytes), "\n")
	instructions := make([]Instruction, 0)
	for _, line := range lines {
		line = strings.Trim(line, " \n\t")
		if len(line) == 0 {
			continue
		}

		tokens := strings.Split(line, " ")
		command := tokens[0]
		args := make([]string, 0)
		for i := 1; i < len(tokens); i++ {
			if len(tokens[i]) == 0 {
				continue
			}
			args = append(args, tokens[i])
		}
		ins := Instruction{Command: command, Args: args}
		instructions = append(instructions, ins)
	}
	return instructions, nil
}
