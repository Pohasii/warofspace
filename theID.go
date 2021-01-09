package main

import "sync"

type ID struct {
	sync.RWMutex
	ID int
}

// Get () int
// safe with Mutex
func (the *ID) Get() int {
	// block for read
	the.RLock()

	// change id
	the.ID++

	// unblock
	defer the.RUnlock()

	// return result
	return the.ID
}
