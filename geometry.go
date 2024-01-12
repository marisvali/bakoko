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

func LineCircleIntersection4Factors(l Line, circle Circle) (bool, Pt) {
	//d = L - E ( Direction vector of ray, from start to end )
	//f = E - C ( Vector from center sphere to ray start )

	r := circle.Diameter.DivBy(I(2)).Plus(circle.Diameter.Mod(I(2)))
	d := l.Start.To(l.End)
	f := circle.Center.To(l.Start)

	a := d.Dot(d)                   // two multiplications
	b := f.Dot(d).Times(I(2))       // two multiplications
	c := f.Dot(f).Minus(r.Times(r)) // two multiplications

	// The computation below involves 4 multiplications. Which means for a unit
	// of 1000, if we have a range of coordinates between 1 and 10.000 we have
	// (10.000*1000)^4 = (10^4*10^3)^4 = 10^28
	// 2^63 = 9,223,372,036,854,775,808 ~= 10^18
	// 8,000,000,000,000,000 (2000 * 100)^3
	// 1,000,000,000,000,000,000 (10000 * 100)^3
	// 9,223,372,036,854,775,808
	// With 10.000.000 (unit of 1000) we get only two multiplications within int64.
	// With 1.000.000 (unit of 100) I would get three multiplications within int64.
	// But I would still need 4 multiplications for this calculation here.
	// Should I reduce my unit or adjust this computation?
	// I know I start from the space of 1 to 10k units and I know I end up there.
	// The only problem I will ever have is with intermediary calculations that
	// need to go above to another space.
	// By my math notes, I can probably go down to 3 multiplications.
	// I definitely need to modify the unit from 1000 to 100 so that I can have
	// 3 multiplications.
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

// Returns the first point of intersection between 'line' and 'circle', where
// first means the intersection point closest to line.Start.
// If line and circle don't intersect, returns false.
func LineCircleIntersection3Factors(line Line, circle Circle) (bool, Pt) {
	// https://stackoverflow.com/questions/1073336/circle-line-segment-collision-detection-algorithm
	// The idea of the algorithm goes like this:
	// - Write the equation of the line in terms of two points on the line.
	// 		x = startX + t * (endX - startX)
	// 		y = startY + t * (endY - endY)
	// 		So a point on the line must have (x, y) coordinates which obey the
	//		equations above.
	// - Write the equation of the circle in terms of center and radius.
	//		(x - center.x)^2 + (y - center.y)^2 = radius^2
	// 		So a point on the circle must have (x, y) coordinates which obey the
	//		equations above.
	// - Intersection points mean those values of x and y for which the 3
	// equations above are true. We have 3 equations with 3 unknowns: x, y and t.
	// - So now we can replace the x and y from the circle equation with the
	// x and y definitions from the line. This way we get a big equation with
	// only one unknown: t. The t factor will tell us where on the line the
	// point of intersection is.
	// - The equation we need to solve is quadratic (a*t^2 + b*t + c = 0).
	// This means we have 3 solutions for the equation:
	// 		- no solution (corresponding to 'line and circle don't intersect')
	//		- a single solution, t (corresponding to 'line and circle touch')
	//		- two solutions, t1 and t2 (corresponding to 'line goes through circle')
	// - After I get values for t, I can compute the points on the line. These
	// points may or may not fall between (x1, y1) and (x2, y2) on the infinite line.
	// If 0 <= t <= 1 then it falls between these two points on the infinite line.

	// There are some intermediate steps between the original equations and the
	// final form that we need to compute. These won't be covered here.
	// Below are the canonical final formulas:
	// d = line.End - line.start (direction vector of ray, from start to end)
	// f = line.Start - circle.Center (vector from center sphere to ray start)
	// t^2 * ( d · d ) + 2*t*( f · d ) + ( f · f - r^2 ) = 0
	// where f · d is the dot product between the f and d vectors.
	// The solution to this quadratic equation:
	// a = (d · d)
	// b = 2 * f · d
	// c = f · f - r^2
	// discriminant = b*b - 4*a*c
	// t1 = (-b - sqrt(discriminant)) / (2*a)
	// t2 = (-b + sqrt(discriminant)) / (2*a)
	// If discriminant < 0 then we have no solution.
	// If discriminant == 0 then we have one solution (t1 == t2).
	// Final intersection points:
	// x1 = startX + t1 * (endX - startX)
	// y1 = startY + t1 * (endY - endY)
	// x2 = startX + t2 * (endX - startX)
	// y2 = startY + t2 * (endY - endY)

	// The issue is, when the line segment intersects the circle, I will have
	// 0 <= t1, t2 <= 1. I want to do this whole algorithm using integers only.
	// This means I can't compute t1 and t2 directly, as they are sub-unitary.
	// I have to manipulate the equations so that I multiply things before and
	// do the division at the end.
	// For example if t1 is a fraction of the form p/q I can compute:
	// x1 = startX + t1 * (endX - startX)
	//	  = startX + t1 * endX - t1 * startX
	//	  = startX + p/q * endX - p/q * startX
	//	  = startX + p * endX / q - p * startX / q

	// Another problem is that discriminant involves b*b and a*c. a, b and c
	// are all two factor numbers. They are the result of multiplying two ints
	// which are in the range of coordinates I use in my world.
	// The discriminant then becomes a four factor number. I know that my int64
	// limit is 10^18. My world coordinates must be at least between 0 and 1920
	// because I have a resolution of 1920x1080. I let's say at least 0 to
	// 10.000, or 10^4. If I multiply 4 numbers that have a maximum of 10^4
	// I get 10^16. And I want to have a unit, to have some precision. I wanted
	// 1000 initially, but that might be too much. But to have 4-factor numbers
	// my unit can't even be 10, because that would multiply my result by 10^4
	// and I would get 10^20. But if I have only 3-factor numbers, suddenly my
	// options improve. Without unit I have 10^12. I can afford a unit of 100,
	// because that multiplies my 10^12 by (10^2)^3 = 10^6, so 10^18.
	// So I have to somehow arrange my equations so that I only ever rise up to
	// 3-factor numbers and never above. I can do this by moving things around
	// so that if I ever have to multiply 4 numbers and divide by another number
	// I multiply 3 of them, divide, then multiply the result with the 4th number.

	// Again, it's not worth noting here all the moving around that I did for
	// the equations, I did them on a piece of paper. The conclusions are:
	// d = line.End - line.start (direction vector of ray, from start to end)
	// f = line.Start - circle.Center (vector from center sphere to ray start)
	// a = (d · d) 					// 2-factor number
	// b = 2 * f · d				// 2-factor number
	// c = f · f - r^2				// 2-factor number
	// distX = |endX - startX|		// 1-factor number
	// kx = b * distX / a			// 3-factor intermediate, then 1-factor
	// lx = (c * distX / a) * distX	// 3-factor intermediate, then 2-factor
	// mx = kx^2 - lx				// 2-factor
	// signX = endX < startX ? 1 : -1
	// x1 = startX + signX * (kx + sqrt(mx))	// 2-factor intermediate
	//											// then 1-factor
	// This way no calculation goes above 3-factor or below 1-factor (sub-unitary).
	// Same equations apply for y1, just replace all x coords with y.

	// We have one more thing to deal with. We have a sqrt, which means we need
	// to check if the input to the sqrt is negative or not (which translates
	// to 'we have intersection points' or 'we don't').

	// Another thing we need to watch out for is that we lose precision when
	// computing kx and lx, because they each include a divide operation. This
	// means kx ad lx get rounded down to the nearest integer. This can actually
	// cause a problem. For example, a test gave me the values:
	// kx = -26.9285263 = -26
	// lx = 710.806481 = 710
	// The integer values gave me
	// mx = 676 - 710 => mx < 0 => no intersection
	// But if I used the floating point values I would get
	// mx = 725.14552868979169 - 710.806481 => mx > 0 => intersection
	// So the loss of precision will tell me that sometimes there is no
	// intersection when there actually is. True, it's a case where the two
	// 'almost don't intersect'. The issue is that I use this intersection
	// algorithm to do collision detection. If I get a false negative in
	// collision detection I will have an object go through another object when
	// it shouldn't and it will ruin the game. But if something collides when
	// it could have almost not collided, that's perfectly fine.
	// So I'd much rather have false positives and zero false negatives.
	// A solution to this is simply to increase kx's magnitude by 1 in order to
	// compensate for the precision lost during the integer division. This will
	// make mx > 0 sometimes when it shouldn't be, but I will never have
	// mx < 0 when it shouldn't be.

	// r = d/2 -> if we get .5 round up instead of down so that the circle is
	// slightly bigger. This way we get false positives but no false negatives.
	r := circle.Diameter.DivBy(I(2)).Plus(circle.Diameter.Mod(I(2)))
	// d = line.End - line.start (direction vector of ray, from start to end)
	d := line.Start.To(line.End)
	if d.X.Eq(I(0)) && d.Y.Eq(I(0)) {
		// no collision since there is no travel
		return false, Pt{}
	}
	// f = line.Start - circle.Center (vector from center sphere to ray start)
	f := circle.Center.To(line.Start)
	// a = (d · d) 				// 2-factor number
	a := d.Dot(d)
	// b = 2 * f · d			// 2-factor number
	b := f.Dot(d).Times(I(2)) // two factors
	// c = f · f - r^2			// 2-factor number
	c := f.Dot(f).Minus(r.Times(r)) // two factors

	// ----- X coordinate -----
	validX, x1 := lineCircleIntersection3FactorsHelper(a, b, c, line.Start.X, line.End.X)
	if !validX {
		return false, Pt{}
	}

	// ----- Y coordinate -----
	validY, y1 := lineCircleIntersection3FactorsHelper(a, b, c, line.Start.Y, line.End.Y)
	if !validY {
		return false, Pt{}
	}

	return true, Pt{x1, y1}
}

// Function to reduce the code in LineCircleIntersection3Factors as the same
// kind of calculations need to be done for coordinates X and Y.
// Read LineCircleIntersection3Factors to understand what's this does.
func lineCircleIntersection3FactorsHelper(a, b, c, start, end Int) (bool, Int) {
	// dist = |end - start|		// 1-factor number
	dist := end.Minus(start).Abs()
	// sign = end < start ? 1 : -1
	sign := I(1)
	if end.Gt(start) {
		sign = I(-1)
	}
	// k = b * dist / a			// 3-factor intermediate, then 1-factor
	// Make k larger by 1 (in absolute terms) to compensate for the rounding
	// error.
	k := b.Times(dist).DivBy(a.Times(I(2))).EnlargedByOne()
	// l = (c * dist / a) * dist	// 3-factor intermediate, then 2-factor
	l := c.Times(dist).DivBy(a).Times(dist)
	// mx = k^2 - l				// 2-factor
	m := k.Times(k).Minus(l)
	if m.IsNegative() {
		return false, I(0)
	}
	// coord = start + sign * (k + sqrt(m))	// 2-factor intermediate
	coord := start.Plus(sign.Times(k.Plus(m.Sqrt())))
	if !coord.Between(start, end) {
		return false, I(0)
	}
	return true, coord
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
	//circles := [1]Circle{
	//	{upperLeftCorner, circleDiameter},
	//}

	for _, c := range circles {
		if intersects, pt := LineCircleIntersection3Factors(travelLine, c); intersects {
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
