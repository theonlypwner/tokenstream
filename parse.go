package tokenstream

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/scanner"
)

// ParseCFG parses a CFG from a reader.
//
// It expects rules in the form of
// rule = <lhs> "->" <rhs> * ("|" <rhs>) ";"
func ParseCFG(grammar io.Reader) (Grammar, error) {
	p := make(map[string][]ProductionRule)

	const (
		SLHS = iota
		SArrow
		SRHS
		SChoiceOrEnd
	)

	state := SLHS
	var lhs string
	var rule ProductionRule

	s := new(scanner.Scanner).Init(grammar)
	for c := s.Scan(); c != scanner.EOF; c = s.Scan() {
		token := s.TokenText()

		switch state {
		case SLHS:
			lhs = token
			state = SArrow
		case SArrow:
			if !(token == "->" || token == "=") {
				return Grammar{}, fmt.Errorf("expected '->' or '=' at %v", s.Position.String())
			}
			state = SRHS
		case SRHS:
			if (token[0] == '\'' || token[0] == '"') && len(token) >= 2 && token[0] == token[len(token)-1] {
				rule = append(rule, Symbol{true, token[1 : len(token)-1]})
			} else {
				rule = append(rule, Symbol{false, token})
			}
			state = SChoiceOrEnd
		case SChoiceOrEnd:
			if token == "|" {
				state = SRHS
			} else if token == ";" {
				p[lhs] = append(p[lhs], rule)

				// lhs = ""
				rule = nil
				state = SLHS
			} else {
				return Grammar{}, fmt.Errorf("expected '|' or ';' at %v", s.Position.String())
			}
		}
	}

	return Grammar{p}, nil
}

// ParseCFGLines parses a CFG from a reader.
//
// It expects one rule per line in the format
// rule = <lhs> "->" <rhs> * ("|" <rhs>) <newline>
func ParseCFGLines(grammar io.Reader) (Grammar, error) {
	p := make(map[string][]ProductionRule)

	scanner := bufio.NewScanner(grammar)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			// skip blank lines
			continue
		}
		if len(tokens) < 2 || !(tokens[1] == "->" || tokens[1] == "=") {
			return Grammar{}, errors.New("production rule must be 'LHS = RHS'")
		}
		lhs := tokens[0]
		var rule ProductionRule
		for _, token := range tokens[2:] {
			if token == "|" {
				p[lhs] = append(p[lhs], rule)
				rule = nil
			} else if (token[0] == '\'' || token[0] == '"') && len(token) >= 2 && token[0] == token[len(token)-1] {
				rule = append(rule, Symbol{true, token[1 : len(token)-1]})
			} else {
				rule = append(rule, Symbol{false, token})
			}
		}
		p[lhs] = append(p[lhs], rule)
	}

	return Grammar{p}, nil
}
