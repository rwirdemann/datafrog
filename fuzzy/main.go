package main

import (
	"fmt"

	"github.com/agnivade/levenshtein"
)

func main() {

	s1 := "update job set description='World', publish_at='<DATE_STR>', publish_trials=1, published_timestamp='<DATE_STR>', tags='', title='Hello' where id=11"
	for i := 1000; i < 1200; i++ {
		s2 := fmt.Sprintf("update job set description='World', publish_at='<DATE_STR>', publish_trials=1, published_timestamp='<DATE_STR>', tags='', title='Hello' where id=%d", i)
		distance := float64(levenshtein.ComputeDistance(s1, s2))
		sum := float64(len(s1)) + float64(len(s2))
		ratio := (sum - distance) / sum

		fmt.Printf("Distance: %f\tRatio: %f\n", distance, ratio)
	}
}
