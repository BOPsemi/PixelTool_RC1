package util

import (
	"math"
)

/*
DeltaLabCalculator :interface of delta lab calculator
*/
type DeltaLabCalculator interface {
	DeltaLab(ref []float64, comp []float64, kvalues []float64) float64
}

// definition of lab for calculation
type lab struct {
	l float64
	a float64
	b float64
}

// definition of k-values for calculation
type consts struct {
	kL float64
	kC float64
	kH float64
}

// deltaLabCalculater str definition
type deltaLabCalculator struct {
	refLab  *lab
	compLab *lab
	kconsts *consts
}

/*
NewDeltaLabCalculator : initializer
*/
func NewDeltaLabCalculator() DeltaLabCalculator {
	obj := new(deltaLabCalculator)

	// initialize properties
	obj.refLab = new(lab)
	obj.compLab = new(lab)
	obj.kconsts = new(consts)

	return obj
}

/*
DeltaLab :
	in	;ref []float64, comp []float64, kvalues []float64
	out	;float64
*/
func (dc *deltaLabCalculator) DeltaLab(ref []float64, comp []float64, kvalues []float64) float64 {
	distance := 0.0

	if (len(ref) * len(comp) * len(kvalues)) != 0 {
		// setup properties
		// -- Reference
		dc.refLab.l = ref[0]
		dc.refLab.a = ref[1]
		dc.refLab.b = ref[2]

		// -- compare target
		dc.compLab.l = comp[0]
		dc.compLab.a = comp[1]
		dc.compLab.b = comp[2]

		// -- k values
		dc.kconsts.kL = kvalues[0]
		dc.kconsts.kC = kvalues[1]
		dc.kconsts.kH = kvalues[2]

		// start calculation
		distance = dc.deltaCalculator(dc.refLab, dc.compLab, dc.kconsts)
	}

	return distance
}

// delta calculator
// calculation sequence is based on CIE 2000
func (dc *deltaLabCalculator) deltaCalculator(ref *lab, comp *lab, kconsts *consts) float64 {
	// Radian to Degree converter
	radianToDegreee := func(radian float64) float64 {
		return (radian * 180.0 / math.Pi)
	}

	// Degree to Radian converter
	degreeToRadian := func(degree float64) float64 {
		return (degree * math.Pi / 180.0)
	}

	deltaLp := comp.l - ref.l
	lAve := (comp.l + ref.l) / 2.0

	// calculate cAve
	c1 := math.Sqrt(math.Pow(ref.a, 2.0) + math.Pow(ref.b, 2.0))
	c2 := math.Sqrt(math.Pow(comp.a, 2.0) + math.Pow(comp.b, 2.0))
	cAve := (c1 + c2) / 2.0

	ap1 := ref.a + (ref.a/2.0)*(1.0-math.Sqrt(math.Pow(cAve, 7.0)/(math.Pow(cAve, 7.0)+math.Pow(25, 7))))
	ap2 := comp.a + (comp.a/2.0)*(1.0-math.Sqrt(math.Pow(cAve, 7.0)/(math.Pow(cAve, 7.0)+math.Pow(25, 7))))

	cp1 := math.Sqrt(math.Pow(ap1, 2.0) + math.Pow(ref.b, 2.0))
	cp2 := math.Sqrt(math.Pow(ap2, 2.0) + math.Pow(comp.b, 2.0))

	cpAve := (cp1 + cp2) / 2.0
	deltaCp := cp2 - cp1

	var hp1 float64
	if ref.b == 0 && ap1 == 0 {
		hp1 = 0.0
	} else {
		hp1 = radianToDegreee(math.Atan2(ref.b, ap1))
		if hp1 < 0 {
			hp1 += 360.0
		}
	}

	var hp2 float64
	if ref.b == 0 && ap1 == 0 {
		hp2 = 0.0
	} else {
		hp2 = radianToDegreee(math.Atan2(comp.b, ap2))
		if hp2 < 0 {
			hp2 += 360.0
		}
	}

	var deltahp float64
	if c1 == 0.0 || c2 == 0.0 {
		deltahp = 0.0
	} else if math.Abs(hp1-hp2) <= 180.0 {
		deltahp = hp2 - hp1
	} else if hp2 <= hp1 {
		deltahp = hp2 - hp1
	} else {
		deltahp = hp2 - hp1 - 360.0
	}
	deltaHp := 2.0 * math.Sqrt(cp1*cp2) * math.Sin(degreeToRadian(deltahp)/2.0)

	var HpAve float64
	if math.Abs(hp1-hp2) > 180.0 {
		HpAve = (hp1 + hp2 + 360.0) / 2.0
	} else {
		HpAve = (hp1 + hp2 + 360.0) / 2.0
	}

	t := 1.0 -
		0.17*math.Cos(degreeToRadian(HpAve-30.0)) +
		0.24*math.Cos(degreeToRadian(2.0*HpAve)) +
		0.32*math.Cos(degreeToRadian(3.0*HpAve+6.0)) -
		0.20*math.Cos(degreeToRadian(4.0*HpAve-63.0))

	sl := 1.0 + ((0.015 * math.Pow(lAve-50.0, 2.0)) / math.Sqrt(20.0+math.Pow(lAve-50.0, 2.0)))
	sc := 1.0 + 0.045*cpAve
	sh := 1.0 + 0.015*cpAve*t

	rt := -2.0 * math.Sqrt(math.Pow(cpAve, 7.0)/(math.Pow(cpAve, 7.0)+math.Pow(25.0, 7.0))) *
		math.Sin(degreeToRadian(60.0*math.Exp(-math.Pow((HpAve-275.0)/25.0, 2.0))))

	result := math.Sqrt(
		math.Pow(deltaLp/(kconsts.kL*sl), 2.0) +
			math.Pow(deltaCp/(kconsts.kC*sc), 2.0) +
			math.Pow(deltaHp/(kconsts.kH*sh), 2.0) +
			rt*(deltaCp/(kconsts.kC*sc))*(deltaHp/(kconsts.kH*sh)))

	return result
}
