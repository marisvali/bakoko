package bakoko

import (
	"math"
	. "playful-patterns.com/bakoko/ints"
)

type Line struct {
	Start Pt
	End   Pt
}

type Circle struct {
	Center   Pt
	Diameter Int
}

type Square struct {
	Center Pt
	Size   Int
}

func LineVerticalLineIntersection(l, vert Line) (bool, Pt) {
	// Check if the Lines even intersect.

	// Check if l's min X is at the right of vertX.
	minX, maxX := MinMax(l.Start.X, l.End.X)
	vertX := vert.Start.X // we assume vert.Start.X == vert.End.X

	if minX.Gt(vertX) {
		return false, Pt{}
	}

	// Or if l's max X is at the left of vertX.
	if maxX.Lt(vertX) {
		return false, Pt{}
	}

	//// Check if l's minY is under the vertMaxY.
	//minY, maxY := MinMax(l.Start.Y, l.End.Y)
	//vertMinY, vertMaxY := MinMax(vert.Start.Y, vert.End.Y)
	//
	//if minY.Gt(vertMaxY) {
	//	return false, Pt{}
	//}
	//
	//// Or if l's max Y is above vertMinY.
	//if maxY.Lt(vertMinY) {
	//	return false, Pt{}
	//}

	vertMinY, vertMaxY := MinMax(vert.Start.Y, vert.End.Y)

	// We know the intersection point will have the X coordinate equal to vertX.
	// We just need to compute the Y coordinate.
	// We have to move along the Y axis the same proportion that we moved along
	// the X axis in order to get to the intersection point.

	//factor := (vertX - l.Start.X) / (l.End.X - l.Start.X) // will always be positive
	//y := l.Start.Y + factor * (l.End.Y - l.Start.Y) // l.End.Y - l.Start.Y will
	// have the proper sign so that Y gets updated in the right direction
	//y := l.Start.Y + (vertX - l.Start.X) / (l.End.X - l.Start.X) * (l.End.Y - l.Start.Y)
	//y := l.Start.Y + (vertX - l.Start.X) * (l.End.Y - l.Start.Y) / (l.End.X - l.Start.X)
	var y Int
	if l.End.X.Eq(l.Start.X) {
		y = l.Start.Y
	} else {
		y = l.Start.Y.Plus((vertX.Minus(l.Start.X)).Times(l.End.Y.Minus(l.Start.Y)).DivBy(l.End.X.Minus(l.Start.X)))
	}

	if y.Lt(vertMinY) || y.Gt(vertMaxY) {
		return false, Pt{}
	} else {
		return true, Pt{vertX, y}
	}
}

func LineHorizontalLineIntersection(l, horiz Line) (bool, Pt) {
	// Check if the Lines even intersect.

	// Check if l's minY is under the vertY.
	minY, maxY := MinMax(l.Start.Y, l.End.Y)
	vertY := horiz.Start.Y // we assume vert.Start.Y == vert.End.Y

	if minY.Gt(vertY) {
		return false, Pt{}
	}

	// Or if l's max Y is above vertY.
	if maxY.Lt(vertY) {
		return false, Pt{}
	}

	//// Check if l's min X is at the right of vertMaxX.
	//minX, maxX := MinMax(l.Start.X, l.End.X)
	//vertMinX, vertMaxX := MinMax(horiz.Start.X, horiz.End.X)
	//
	//if minX.Gt(vertMaxX) {
	//	return false, Pt{}
	//}
	//
	//// Or if l's max X is at the left of vertMinX.
	//if maxX.Lt(vertMinX) {
	//	return false, Pt{}
	//}

	vertMinX, vertMaxX := MinMax(horiz.Start.X, horiz.End.X)

	// We know the intersection point will have the Y coordinate equal to vertY.
	// We just need to compute the X coordinate.
	// We have to move along the X axis the same proportion that we moved along
	// the Y axis in order to get to the intersection point.

	//factor := (vertY - l.Start.Y) / (l.End.Y - l.Start.Y) // will always be positive
	//x := l.Start.X + factor * (l.End.X - l.Start.X) // l.End.X - l.Start.X will
	// have the proper sign so that Y gets updated in the right direction
	//x := l.Start.X + (vertY - l.Start.Y) / (l.End.Y - l.Start.Y) * (l.End.X - l.Start.X)
	//x := l.Start.X + (vertY - l.Start.Y) * (l.End.X - l.Start.X) / (l.End.Y - l.Start.Y)
	var x Int
	if l.End.Y.Eq(l.Start.Y) {
		x = l.Start.X
	} else {
		x = l.Start.X.Plus((vertY.Minus(l.Start.Y)).Times(l.End.X.Minus(l.Start.X)).DivBy(l.End.Y.Minus(l.Start.Y)))
	}

	if x.Lt(vertMinX) || x.Gt(vertMaxX) {
		return false, Pt{}
	} else {
		return true, Pt{x, vertY}
	}
}

// (ΔABC) = (1/2) |x1(y2 − y3) + x2(y3 − y1) + x3(y1 − y2)|
func LineCircleIntersection(l Line, circle Circle) (bool, Pt) {
	//d = L - E ( Direction vector of ray, from start to end )
	//f = E - C ( Vector from center sphere to ray start )

	r := circle.Diameter.DivBy(I(2)).Plus(circle.Diameter.Mod(I(2)))
	d := l.Start.To(l.End)
	f := circle.Center.To(l.Start)

	a := d.Dot(d)
	b := f.Dot(d).Times(I(2))
	c := f.Dot(f).Minus(r.Times(r))

	discriminant := b.Times(b).Minus(a.Times(c).Times(I(4)))
	if discriminant.Lt(I(0)) {
		// no intersection
		return false, Pt{}
	} else {
		// ray didn't totally miss sphere,
		// so there is a solution to
		// the equation.

		discriminant = discriminant.Sqrt()

		// either solution may be on or off the ray so need to test both
		// t1 is always the smaller value, because BOTH discriminant and
		// a are nonnegative.
		numarator1 := I(0).Minus(b).Minus(discriminant)
		numarator2 := I(0).Minus(b).Plus(discriminant)
		numitor := a.Times(I(2))

		//t1 := numarator1.DivBy(numitor)
		//t2 := numarator2.DivBy(numitor)

		// 3x HIT cases:
		//          -o->             --|-->  |            |  --|->
		// Impale(t1 hit,t2 hit), Poke(t1 hit,t2>1), ExitWound(t1<0, t2 hit),

		// 3x MISS cases:
		//       ->  o                     o ->              | -> |
		// FallShort (t1>1,t2>1), Past (t1<0,t2<0), CompletelyInside(t1<0, t2>1)

		//t1 >= 0 && t1 <= 1
		t1GreaterThanZero := numarator1.Gt(I(0)) == numitor.Gt(I(0))
		t1LessThanOne := numarator1.Abs().Leq(numitor.Abs())
		//if t1.Geq(I(0)) && t1.Leq(I(1)) {
		if t1GreaterThanZero && t1LessThanOne {
			// t1 is the intersection, and it's closer than t2
			// (since t1 uses -b - discriminant)
			// Impale, Poke

			//P = E + t * d
			//P = l.Start + t1 * d
			//Pt{l.Start.X + t1 * d.X, l.Start.Y + t1 * d.Y}

			x := l.Start.X.Plus(numarator1.Times(d.X).DivBy(numitor))
			y := l.Start.Y.Plus(numarator1.Times(d.Y).DivBy(numitor))
			return true, Pt{x, y}
		}

		// here t1 didn't intersect so we are either started
		// inside the sphere or completely past it

		//t2 >= 0 && t2 <= 1
		t2GreaterThanZero := numarator2.Gt(I(0)) == numitor.Gt(I(0))
		t2LessThanOne := numarator2.Abs().Leq(numitor.Abs())
		if t2GreaterThanZero && t2LessThanOne {
			//if t2.Geq(I(0)) && t2.Leq(I(1)) {
			// ExitWound

			x := l.Start.X.Plus(numarator2.Times(d.X).DivBy(numitor))
			y := l.Start.Y.Plus(numarator2.Times(d.Y).DivBy(numitor))
			return true, Pt{x, y}
		}

		// no intn: FallShort, Past, CompletelyInside
		return false, Pt{}
	}
}

func CirclesIntersect(c1, c2 Circle) bool {
	maxDist := c1.Diameter.Plus(c2.Diameter).DivBy(I(2))
	squaredMaxDist := maxDist.Sqr()
	return c1.Center.SquaredDistTo(c2.Center).Leq(squaredMaxDist)
}

type DebugInfo struct {
	Points  []Pt
	Lines   []Line
	Circles []Circle
}

// CircleSquareCollision doesn't return circleOldPos as a collision point.
func CircleSquareCollision(circleOldPos Pt, circleNewPos Pt,
	circleDiameter Int, s Square) (intersects bool,
	circlePositionAtCollision Pt, collisionNormal Pt,
	debugInfo DebugInfo) {
	// Get the line on which the circle is travelling.
	// Consider the circle to be a point and grow the square using the Minkowski sum concept.
	// Get 4 Circles (one in each corner) and 4 Lines.
	// Compute the intersection between the travel line and the 4 Circles and 4 Lines.
	// If multiple intersection Points exist, get the one closest to circleOldPos.
	// That's the point where the circle will be when it starts touching the square.
	travelLine := Line{circleOldPos, circleNewPos}

	// size / 2 + size % 2 to compensate for the potential precision of the division
	halfSize := s.Size.DivBy(I(2)).Plus(s.Size.Mod(I(2)))

	// square corners
	upperLeftCorner := Pt{s.Center.X.Minus(halfSize), s.Center.Y.Minus(halfSize)}
	lowerLeftCorner := Pt{s.Center.X.Minus(halfSize), s.Center.Y.Plus(halfSize)}
	upperRightCorner := Pt{s.Center.X.Plus(halfSize), s.Center.Y.Minus(halfSize)}
	lowerRightCorner := Pt{s.Center.X.Plus(halfSize), s.Center.Y.Plus(halfSize)}

	// square Lines, moved according to the Minkowski sum concept
	circleRadius := circleDiameter.DivBy(I(2)).Plus(s.Size.Mod(I(2)))

	leftLine := Line{
		Pt{lowerLeftCorner.X.Minus(circleRadius), lowerLeftCorner.Y},
		Pt{upperLeftCorner.X.Minus(circleRadius), upperLeftCorner.Y}}

	rightLine := Line{
		Pt{lowerRightCorner.X.Plus(circleRadius), lowerRightCorner.Y},
		Pt{upperRightCorner.X.Plus(circleRadius), upperRightCorner.Y}}

	upLine := Line{
		Pt{upperLeftCorner.X, upperLeftCorner.Y.Minus(circleRadius)},
		Pt{upperRightCorner.X, upperRightCorner.Y.Minus(circleRadius)}}

	downLine := Line{
		Pt{lowerLeftCorner.X, lowerLeftCorner.Y.Plus(circleRadius)},
		Pt{lowerRightCorner.X, lowerRightCorner.Y.Plus(circleRadius)}}

	// get intersections between the travel line and the (expanded) square's Lines
	var intersectionPoints []Pt
	var intersectionNormals []Pt

	if intersects, pt := LineVerticalLineIntersection(travelLine, leftLine); intersects {
		intersectionPoints = append(intersectionPoints, pt)
		intersectionNormals = append(intersectionNormals, IPt(-1, 0))
	}

	if intersects, pt := LineVerticalLineIntersection(travelLine, rightLine); intersects {
		intersectionPoints = append(intersectionPoints, pt)
		intersectionNormals = append(intersectionNormals, IPt(1, 0))
	}

	if intersects, pt := LineHorizontalLineIntersection(travelLine, upLine); intersects {
		intersectionPoints = append(intersectionPoints, pt)
		intersectionNormals = append(intersectionNormals, IPt(0, -1))
	}

	if intersects, pt := LineHorizontalLineIntersection(travelLine, downLine); intersects {
		intersectionPoints = append(intersectionPoints, pt)
		intersectionNormals = append(intersectionNormals, IPt(0, 1))
	}

	// get intersections between the travel line and the (expanded) square's corner Circles
	circles := [4]Circle{
		{upperLeftCorner, circleDiameter},
		{lowerLeftCorner, circleDiameter},
		{upperRightCorner, circleDiameter},
		{lowerRightCorner, circleDiameter},
	}

	for _, c := range circles {
		if intersects, pt := LineCircleIntersection(travelLine, c); intersects {
			intersectionPoints = append(intersectionPoints, pt)
			intersectionNormals = append(intersectionNormals, c.Center.To(pt))
		}
	}

	debugInfo.Points = intersectionPoints[:]
	debugInfo.Lines = []Line{leftLine, rightLine, upLine, downLine}
	debugInfo.Circles = circles[:]

	// Find the intersection point closest to the start point, but discard those
	// which are exactly the start point.
	minDist := I(math.MaxInt64)
	minIdx := -1
	for idx, pt := range intersectionPoints {
		dist := circleOldPos.SquaredDistTo(pt)
		// Get the point which is closest to but not quite the start point.
		if dist.Gt(I(0)) && dist.Lt(minDist) {
			minDist = dist
			minIdx = idx
		}
	}

	if minIdx < 0 {
		return false, Pt{}, Pt{}, debugInfo
	} else {
		return true, intersectionPoints[minIdx], intersectionNormals[minIdx], debugInfo
	}
}

// CircleSquareCollisionMultiple doesn't return circleOldPos as a collision point.
func CircleSquaresCollision(circleOldPos Pt, circleNewPos Pt,
	circleDiameter Int, squares []Square) (intersectsAny bool,
	circlePositionAtCollision Pt, collisionNormal Pt) {

	minDist := I(math.MaxInt64)
	for _, s := range squares {
		intersects, pt, normal, _ :=
			CircleSquareCollision(circleOldPos, circleNewPos, circleDiameter, s)

		dist := circleOldPos.To(pt).Len()
		if intersects && dist.Lt(minDist) {
			minDist = dist
			circlePositionAtCollision = pt
			collisionNormal = normal
			intersectsAny = true
		}
	}

	return intersectsAny, circlePositionAtCollision, collisionNormal
}
