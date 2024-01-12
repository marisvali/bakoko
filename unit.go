package bakoko

import . "playful-patterns.com/bakoko/ints"

/*
Unit is used in order to get some advantages of floating Points when doing
integer-only calculations.

The problem:
I want to use integers for all my calculations. The reason and the solution are
explained in the ints package. However, this means that the smallest unit of
change is 1. If I have an object at position (10, 5) and I want it at position
(30, 5) I can only move with a minimum speed of 1 on the X axis. Which might be
too fast for my needs.

This problem comes up in all sorts of other places. For example whenever I do a
division I lose some precision because I round to the nearest integer. This can
quickly become a problem.

The solution proposed here:
Use a Unit everywhere instead of 1. So instead of an object at position 10 going
to position 30, it is an object at position (10*Unit, 5*Unit) going to position
(30*Unit, 5*Unit). If this Unit is 1000, the speed can be as low as 1
milli-Unit.

For convenience, you can express 10*Unit by U(10). and 10*milli-Units by MU(10).

This gives me two things:
1. It allows me to think in terms of reasonable dimensions, instead of using
positions like (10000, 5000) and (30000, 5000) everywhere. I can think in terms
of tens, hundreds and thousands and know that I can have subunit values when
I need them. For example: a pixel on my screen can represent a unit. This way
I can easily reason about sprites, their sizes and their positions when I'm
drawing and debugging them.
2. It allows me to experiment with how big a unit needs to be. I currently have
no idea what kind of operations I will need and how much leeway I need to give
myself. int64 is big, but it's not infinite. If I multiply 4 numbers together,
they must all be smaller than 55108, or I overflow. Suddenly int64 isn't that
big. So I need to be careful to have a Unit that's big enough to give me the
precision I want in my computations. But it needs to be small enough so that
I don't overflow in my computations. Since I won't know what my computations
will be until I finish the game, I need this flexibility established from the
start.
*/
//const Unit = 1000

const Unit = 100

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
