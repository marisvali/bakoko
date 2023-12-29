package bakoko

import . "playful-patterns.com/bakoko/ints"

const Unit = 1000

func Units(numUnits Int) Int {
	return numUnits.Times(I(Unit))
}

func Milliunits(numMilliunits Int) Int {
	return numMilliunits.Times(I(Unit / 1000))
}

func U(numUnits int64) Int {
	return I(numUnits).Times(I(Unit))
}

func UPt(xUnits int64, yUnits int64) Pt {
	return Pt{I(xUnits).Times(I(Unit)), I(yUnits).Times(I(Unit))}
}

func MU(numUnits int64) Int {
	return I(numUnits).Times(I(Unit)).DivBy(I(1000))
}
