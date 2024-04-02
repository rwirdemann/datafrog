package main

import (
	"fmt"

	"github.com/agnivade/levenshtein"
)

func main() {
	s1 := "update job set description='World', publish_at='<DATE_STR>', publish_trials=1, published_timestamp='<DATE_STR>', tags='', title='Hello' where id=11"
	s2 := "update job set description='World', publish_at='<DATE_STR>', publish_trials=1, published_timestamp='<DATE_STR>', tags='', title='Hello' where id=12"
	distance(s1, s2)

	s1 = "delete from job where id=37"
	s2 = "delete from job where id=37"

	distance(s1, s2)
}

func distance(s1 string, s2 string) {
	distance := float64(levenshtein.ComputeDistance(s1, s2))
	sum := float64(len(s1)) + float64(len(s2))
	ratio := (sum - distance) / sum
	fmt.Printf("Distance: %f\tRatio: %f\n", distance, ratio)

}
