package main

import "sort"

func NonConstructibleChange(coins []int) int {
	if len(coins) == 0 {
		return 1
	}
	sort.Ints(coins)
	change := 0
	sum := 0
	for i := 0; i < len(coins); i++ {
		if sum+coins[i] != change+1 {
			return change + 1
		}
		sum += coins[i]
		change++
	}
	return coins[len(coins)-1] + 1
}
