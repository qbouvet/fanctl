package curve

import "fmt"



type Curve interface {
	Lookup(int) int
	MapY(int, int, int, int) Curve
}



/*
 *	Linear interpolation between a min and a max. Should work well for 
 *	everything
 */
type ClampedLinearCurve struct {
	x1, y1, x2, y2 int
}

func ClampedLinear(x1, y1, x2, y2 int) Curve {
	if x1>=x2 {
		panic("Invalid parameters: x1>=x2")
	}
	return &ClampedLinearCurve{x1, y1, x2, y2}
}

func (C *ClampedLinearCurve) Lookup(x int) int {
	if x<=C.x1 {
		return C.y1
	}
	if x>=C.x2 {
		return C.y2
	}
	return C.y1 + (C.y2-C.y1) * (x-C.x1) / (C.x2-C.x1)
	//completion := (x-C.x1)/(C.x2-C.x1)
	//return C.y2*completion + C.y1*(1-completion)
}

 //	Actuators transform the existing curve to directly fit to their min/max
func (C *ClampedLinearCurve) MapY(min, max, newmin, newmax int) Curve {
	newy1 := newmin + C.y1 * (newmax-newmin) / (max-min)
	newy2 := newmin + C.y2 * (newmax-newmin) / (max-min)
	fmt.Printf("Remapped y1=%d->%d, y2=%d->%d ([%d-%d]->[%d-%d])\n", 
		C.y1, newy1, C.y2, newy2, min, max, newmin, newmax,
	)
	return &ClampedLinearCurve {C.x1, newy1, C.x2, newy2}
}