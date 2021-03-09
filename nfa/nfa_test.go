package nfa

import (
	"testing"
)

func dagTransitions(st state, sym rune) []state {
	/*
	 * 0 -a-> 1
	 * 0 -a-> 2
	 * 1 -b-> 3
	 * 2 -c-> 3
	 */
	return map[state]map[rune][]state{
		0: map[rune][]state{
			'a': []state{1, 2},
		},
		1: map[rune][]state{
			'b': []state{3},
		},
		2: map[rune][]state{
			'c': []state{3},
		},
	}[st][sym]
}

func expTransitions(st state, sym rune) []state {
	/*
	 * 0 -a-> 1
	 * 0 -a-> 2
	 * 0 -b-> 2
	 * 1 -b->0
	 */
	return map[state]map[rune][]state{
		0: map[rune][]state{
			'a': []state{1, 2},
			'b': []state{2},
		},
		1: map[rune][]state{
			'b': []state{0},
		},
		2: map[rune][]state{},
	}[st][sym]
}

func langTransitions(st state, sym rune) []state {
	/*
	 * 0 -a-> 0
	 * 0 -b-> 1
	 * 1 -a-> 1
	 * 1 -b-> 0
	 */
	return map[state]map[rune][]state{
		0: map[rune][]state{
			'a': []state{0},
			'b': []state{1},
		},
		1: map[rune][]state{
			'a': []state{1},
			'b': []state{0},
		},
	}[st][sym]
}

func TestReachable(t *testing.T) {
	nfas := map[string]TransitionFunction{
		"dagTransitions":  dagTransitions,
		"expTransitions":  expTransitions,
		"langTransitions": langTransitions,
	}

	tests := []struct {
		nfa          string
		start, final state
		input        []rune
		expected     bool
	}{
		{"dagTransitions", 0, 3, []rune{'a', 'b'}, true},
		{"dagTransitions", 0, 3, []rune{'a', 'c'}, true},
		{"dagTransitions", 0, 1, []rune{'a'}, true},
		{"dagTransitions", 0, 0, nil, true},
		{"dagTransitions", 0, 3, []rune{'a', 'a'}, false},
		{"dagTransitions", 0, 3, []rune{'a'}, false},
		{"dagTransitions", 0, 1, []rune{'b'}, false},
		{"dagTransitions", 0, 0, []rune{'b'}, false},

		{"expTransitions", 0, 0, []rune{'a', 'b'}, true},
		{"expTransitions", 0, 2, []rune{'a', 'b', 'a'}, true},
		{"expTransitions", 0, 2, []rune{'a', 'b', 'a', 'b', 'a'}, true},
		{"expTransitions", 0, 0, []rune{'a', 'a'}, false},
		{"expTransitions", 0, 2, []rune{'a', 'b', 'a', 'b'}, false},

		{"langTransitions", 0, 0, []rune{'a', 'b', 'b'}, true},
		{"langTransitions", 0, 1, []rune{'a', 'a', 'b'}, true},
		{"langTransitions", 0, 0, []rune{'a', 'a', 'a', 'a', 'a'}, true},
		{"langTransitions", 0, 0, nil, true},
		{"langTransitions", 0, 1, []rune{'a', 'a'}, false},
		{"langTransitions", 0, 0, []rune{'a', 'b', 'a', 'a'}, false},

		// TODO add more tests for 100% test coverage
	}

	for _, test := range tests {
		func() {
			defer func() {
				if recover() != nil {
					t.Errorf("Reachable panicked on (%s, %d, %d, %v); expected %t",
						test.nfa, test.start, test.final, string(test.input),
						test.expected)
				}
			}()
			actual := Reachable(nfas[test.nfa], test.start, test.final, test.input)
			if test.expected != actual {
				t.Errorf("Reachable failed on (%s, %d, %d, %v); expected %t, actual %t",
					test.nfa, test.start, test.final, string(test.input),
					test.expected, actual)
			}
		}()
	}
}

func BenchmarkReachable(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Reachable(expTransitions, 0, 2, ([]rune)("aaaaaaaaaaaaaa"))
	}
}
