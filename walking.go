package bakoko

import (
	. "playful-patterns.com/bakoko/ints"
)

func ObstacleFree(m Matrix, upperLeft, lowerRight Pt) bool {
	for y := upperLeft.Y; y.Leq(lowerRight.Y); y.Inc() {
		for x := upperLeft.X; x.Leq(lowerRight.X); x.Inc() {
			if !m.InBounds(Pt{x, y}) || m.Get(y, x).Neq(I(0)) {
				return false
			}
		}
	}
	return true
}

func GetWalkableMatrix(m Matrix, mSquareSize Int, charSize Int) (mw Matrix, sizeW Int, offset Pt) {
	mw.Init(m.NRows().Times(I(2)).Plus(I(1)), m.NCols().Times(I(2)).Plus(I(1)))
	sizeW = mSquareSize.DivBy(I(2))
	//offset.X = mSquareSize.DivBy(I(2)).Negative()
	//offset.Y = mSquareSize.DivBy(I(2)).Negative()

	for y := I(0); y.Lt(mw.NRows()); y.Inc() {
		for x := I(0); x.Lt(mw.NCols()); x.Inc() {
			// Build rectangle in world coordinates.
			var worldUpperLeft Pt
			worldUpperLeft.X = x.Times(sizeW).Minus(charSize.DivBy(I(2)))
			worldUpperLeft.Y = y.Times(sizeW).Minus(charSize.DivBy(I(2)))
			var worldLowerRight Pt
			worldLowerRight.X = worldUpperLeft.X.Plus(charSize)
			worldLowerRight.Y = worldUpperLeft.Y.Plus(charSize)

			// Translate rectangle to original matrix coordinates.
			var matrixUpperLeft Pt
			matrixUpperLeft.X = worldUpperLeft.X.DivBy(mSquareSize)
			matrixUpperLeft.Y = worldUpperLeft.Y.DivBy(mSquareSize)

			var matrixLowerRight Pt
			matrixLowerRight.X = worldLowerRight.X.DivBy(mSquareSize)
			if worldLowerRight.X.Mod(mSquareSize).Eq(I(0)) {
				matrixLowerRight.X.Dec()
			}
			matrixLowerRight.Y = worldLowerRight.Y.DivBy(mSquareSize)
			if worldLowerRight.Y.Mod(mSquareSize).Eq(I(0)) {
				matrixLowerRight.Y.Dec()
			}

			if !ObstacleFree(m, matrixUpperLeft, matrixLowerRight) {
				mw.Set(y, x, I(1))
			} else {
				mw.Set(y, x, I(0))
			}
		}
	}
	return
}
