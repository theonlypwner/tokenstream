package tokenstream

import (
	"math/rand"
)

// Grammar represents a context-free grammar.
type Grammar struct {
	// Productions is a map from a non-terminal (LHS) to its set of production rules (RHS)
	Productions map[string][]ProductionRule
}

// Symbol represents a symbol in a Grammar, which can be
// either terminal or non-terminal
type Symbol struct {
	Terminal bool
	Value    string
}

// ProductionRule is the right-hand-side of a production
type ProductionRule []Symbol

// TokenStream produces a string of tokens from a grammar
type TokenStream struct {
	g     *Grammar
	stack []*Symbol
}

// MakeTokenStream creates a TokenStream for a Grammar and start symbol.
func MakeTokenStream(g *Grammar, startSymbol string) TokenStream {
	return TokenStream{
		g,
		[]*Symbol{{false, startSymbol}},
	}
}

// Next gets the next token
func (t *TokenStream) Next() string {
	for len(t.stack) != 0 {
		production := t.stack[len(t.stack)-1]

		if !production.Terminal {
			if p, ok := t.g.Productions[production.Value]; ok && len(p) != 0 {
				t.stack = t.stack[:len(t.stack)-1]

				// Add non-terminals to stack
				rule := p[rand.Intn(len(p))]
				for i := len(rule) - 1; i >= 0; i-- {
					t.stack = append(t.stack, &rule[i])
				}
				continue
			}
			// If rule does not exist, treat symbol as terminal
		}

		// Return terminal
		r := production.Value
		t.stack = t.stack[:len(t.stack)-1]
		return r
	}

	// Stack is empty (end of string)
	return ""
}

// Produce generates a string of tokens but stops after a length limit is reached.
// The length of the returned string may exceed stopAfter.
func (t *TokenStream) Produce(stopAfter uint) string {
	var s []byte
	for {
		nextToken := t.Next()
		if nextToken == "" {
			break
		}

		s = append(s, nextToken...)
		if uint(len(nextToken)) > stopAfter {
			s = append(s, "..."...)
			break
		}
		stopAfter -= uint(len(nextToken))
	}

	return string(s)
}
