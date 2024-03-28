package main

import (
	"fmt"

	"github.com/agnivade/levenshtein"
)

func main() {
	s1 := "update job set description='World', publish_at='2024-03-28 10:57:27', publish_trials=1, published_timestamp='2024-03-28 10:57:36.261562', tags='', title='Hello' where id=11"
	s2 := "update job set description='World', publish_at='2024-03-28 10:57:27', publish_trials=1, published_timestamp='2024-03-28 10:57:36.261562', tags='', title='Hello' where id=11"
	distance := float64(levenshtein.ComputeDistance(s1, s2))
	sum := float64(len(s1)) + float64(len(s2))
	ratio := (sum - distance) / sum

	fmt.Printf("The distance / ratio is %f %f.\n", distance, ratio)
}
