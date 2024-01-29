package interpreter

import "strings"

type parserManager struct{}

// parse parses a given line of code and returns an Instruction.
func (p *parserManager) parse(line string) Instruction {
	tokens := strings.Split(line, " ")
	command := tokens[0]
	args := make([]string, 0)
	isStringLiteral := false
	for i := 1; i < len(tokens); i++ {
		if isStringLiteral {
			lastArg := args[len(args)-1]
			lastArg += " " + tokens[i]
			if strings.HasSuffix(tokens[i], "\"") {
				isStringLiteral = false
			}
			args[len(args)-1] = lastArg
		} else {
			if len(tokens[i]) == 0 {
				continue
			}
			if strings.HasPrefix(tokens[i], "\"") && !strings.HasSuffix(tokens[i], "\"") {
				isStringLiteral = true
			}
			args = append(args, tokens[i])
		}

	}
	return Instruction{Command: command, Args: args}
}
