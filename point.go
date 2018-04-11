package mapSimplify

import (
	"fmt"
)

type Point [2]float64

func (pt Point) String() string {
	return fmt.Sprintf("[%f, %f]", pt[0], pt[1])
}


