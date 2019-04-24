package main

import (
	"testing"
)

func TestRegistrar(t *testing.T) {
	t.Run(
		"Registrar should return a unique value/id for each process that calls Register",
		func(t *testing.T) {
			r := Registrar{}
			ids := make(map[int]bool, 0)
			idsChan := make(chan int, 100)
			for i := 0; i < 100; i++ {
				go func() {
					id := r.Register()
					idsChan <- id
				}()
			}
			for i := 0; i < 100; i++ {
				id := <-idsChan
				_, ok := ids[id]
				if ok {
					t.Errorf("Expected Regestrar to only return unique id values. Duplicate id value: %d", id)
				}
				ids[id] = true
			}
		},
	)
}
