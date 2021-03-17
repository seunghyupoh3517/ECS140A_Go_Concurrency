package nfa

import (
	"sync"
	// "fmt"
)

// A nondeterministic Finite Automaton (NFA) consists of states,
// symbols in an alphabet, and a transition function.

// A state in the NFA is represented as an unsigned integer.
type state uint

// Given the current state and a symbol, the transition function
// of an NFA returns the set of next states the NFA can transition to
// on reading the given symbol.
// This set of next states could be empty.
type TransitionFunction func(st state, sym rune) []state
var root state
var routine_counter int		// use to count how many threads we create
var mu sync.Mutex		// lock for the counter

// Reachable returns true if there exists a sequence of transitions
// from `transitions` such that if the NFA starts at the start state
// `start` it would reach the final state `final` after reading the
// entire sequence of symbols `input`; Reachable returns false otherwise.
func Reachable(
	// `transitions` tells us what our NFA looks like
	transitions TransitionFunction,
	// `start` and `final` tell us where to start, and where we want to end up
	start, final state,
	// `input` is a (possible empty) list of symbols to apply.
	input []rune,
) bool {
	result := make(chan bool)

	root = start
 	go goReachable(transitions, start, final, input, result)
	
	// if val, ok := <- result; ok {
	// 	fmt.Println(" $$$ In the main !! !!!")
	// 	close(result)
	// 	return val
	// }

	return <- result
}

func goReachable(transitions TransitionFunction, start, final state, input []rune, result chan bool) {

	var wg sync.WaitGroup
	if len(input) == 0 {
		if start == final {
			// send to channel only when its final step
			// fmt.Println(" *** WRITE true!!!")
			result <- true
			// close(result)
		
			// if _, ok := <- result; ok {
			// 	fmt.Println(" $$$ WRITE true!!!")
			// 	result <- true
			// 	close(result)
			// }	
			return
		} 
	} else {
		next := transitions(start, input[0])

		wg.Add(len(next))
		// fmt.Println(" $$$ Add():", len(next))
		for i, _ := range next {
			// fmt.Println(" $$$ index i =",i)

			var counter = 0
			mu.Lock()
			counter = routine_counter
			routine_counter++
			mu.Unlock()

			if counter < 10 {
				go func(next_state state) {
					// fmt.Println(" *** Test Concurrency!!!!")
					goReachable(transitions, next_state, final, input[1:], result)
					wg.Done()
				}(next[i])
			} else {
				goReachable(transitions, next[i], final, input[1:], result)
				wg.Done()
			}			
		}	
		wg.Wait() 
	}

	// only one place can WRITE false, which is the root
	if start == root {
		// if _, ok := <- result; ok {
			// fmt.Println(" $$$ WRITE false!!!")
			result <- false
		// }
	}

	return
}
