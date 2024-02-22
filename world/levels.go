package world

import . "playful-patterns.com/bakoko/ints"

func RandomLevel(nRows, nCols, extraMin, extraMax Int) (m Matrix) {
	m.Init(nRows, nCols)

	// Set left and right borders.
	for row := I(0); row.Lt(m.NRows()); row.Inc() {
		m.Set(row, I(0), I(1))
		m.Set(row, m.NCols().Minus(I(1)), I(1))
	}

	// Set top and bottom borders.
	for col := I(0); col.Lt(m.NCols()); col.Inc() {
		m.Set(I(0), col, I(1))
		m.Set(m.NRows().Minus(I(1)), col, I(1))
	}

	extra := RInt(extraMin, extraMax)
	for i := I(0); i.Lt(extra); i.Inc() {
		row := RInt(I(1), m.NRows().Minus(I(1)))
		col := RInt(I(1), m.NCols().Minus(I(1)))
		m.Set(row, col, I(1))
	}

	return
}

func ManualLevel() (m Matrix) {
	var level string
	level = `
xxxxxxxxxxxxx
x           x
x  x     x  x
x           x
x           x
x    x  x   x
x           x
x           x
x           x
x         xxx
x           x
x           x
x           x
x           x
xxxxxxxxxxxxx
	`
	//
	//	level = `
	//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	//x   x                         x      x                 xx       x                x      x
	//x   x     xx              x   x      x                  x       x    xxxxxxx     x  x   x
	//x        xxx     x         x x       x          xx      x       x     xxxxxx     x      x
	//xxx              x          x        x         xx       x       x                x  x   x
	//x       x        x     xxxxxxxxxxxx   xx      xx        x  x    xxxxx      xx           x
	//x       xx       x xxxxxxxxxx           xx  xx             x             xxxxx     x xx x
	//xxx              x     x          x       xx               x   xx         xxxx          x
	//x      x         x     x          x             xxxxxxxxxxxx    x    xx            xxxxxx
	//x     xxx                    xxxx x                        x                  xxx       x
	//x      x         x           x    x          xxx      x    xxxxxxxxxxxxxxxxx            x
	//xxx             xxx     x    x    x    x    xxxx    xxx    x                            x
	//x    xxxx        x                x    x            x x         xxx  x                  x
	//x      xx   x         x           x    xxxx         x             xxxx   xxxxxx         x
	//xxx    xx   x       xxx           x                                       xxxxx         x
	//x           x       x   x   xx    x    x        xxxxxxxxxxxxxxxxxxx                     x
	//x           x           x   xx         x                                                x
	//x           x           x   xx     xxxxxxxxx      x    x     x     x   x   x    x  x    x
	//x           x       x                             x    xxxxxxx     x   x   x    x  x    x
	//x           x       x                   xxxxxx    x         x      x   x   x    x  x    x
	//x           x     xxxxx                 xxxxx     x         x      x   x   x    x  x    x
	//x       x                 xxxxxxxxxxx   xxxxx     x         x                           x
	//x     xxx    xx                                             x      xxxxxxxxxxxxxxxx     x
	//x   xxxxxx       xxx    xxxx                      xxxx      x                           x
	//x     xxx               x       x         x       xxxx      x      xxxxxxxxxxxxxxxx     x
	//x       x        xxxxxxxx       x         x       xxxx      x                           x
	//x                               x         x       xxxx      x      xxxxxxxxxxxxxxxx     x
	//x    xxxxxxxxxxxxxxxxxxxxxxxxxxxx         x                 x                           x
	//x                               x         x          xx     x      xxxxxxxxxxxxxxxx     x
	//xxxxxxxxxxxxxxxxxx              x                   xxx     x                           x
	//x                      xxxxxxxxxxxxxxxxxxx        xxxxx            xx     xxxxxxxx      x
	//x      xxxx                     x                xxxxxx              xxxxx         xx   x
	//x     x    xx                                               xxxxxx                  x   x
	//x     x      xx    xxxx                    xx  xx           xxxxxxx         x  xx   x   x
	//x      x             x          x   xxx                                     x  xx       x
	//x    x  x            x   x  x   x   xxx    xx  xx                           x       x   x
	//x   x    x           x          x   xxx    xx  xx                           x  xx   x   x
	//x  x                 x          x                              x         xxxx  xx   x   x
	//x     xxxxxxxx       x   x  x   x    x     xxxxxxx             x            x       x   x
	//x                    x          x    x     xxxxxxx             x            x    xxxxxxxx
	//x         x                     x    x                         x            x           x
	//x   x      xxx         x x      x     xxxxxx    x     xxxxxxxxxxxxxxxxxxx   x           x
	//x   x                           x               x              x            x    xxxxxxxx
	//x   x      xxxxxxxx    x x      x       xxx     x              x    x x     x           x
	//x   x      xxxxxxxx             x         x     xxxxxxxxxxx   x     x x     x    xxxxxxxx
	//x   x         xxxxx             x         x                    x            x           x
	//x                               x     xxxxxxxxxx     x               x                  x
	//x                                                    x               x                  x
	//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	//`

	row := -1
	col := 0
	maxCol := 0
	for i := 0; i < len(level); i++ {
		c := level[i]
		if c == '`' {
			continue
		} else if c == '\n' {
			maxCol = col
			col = 0
			row++
			continue
		}
		col++
	}
	//m.Init(I(int64(row+1)), I(int64(maxCol/3+1)))
	m.Init(I(row), I(maxCol))

	row = -1
	col = 0
	for i := 0; i < len(level); i++ {
		c := level[i]
		if c == '`' {
			continue
		} else if c == '\n' {
			col = 0
			row++
			continue
		} else if c == 'x' {
			m.Set(I(row), I(col), I(1))
		}
		col++
	}
	//row = -1
	//col = 0
	//for i := 0; i < len(level); i++ {
	//	c := level[i]
	//	if c == '`' {
	//		continue
	//	} else if c == '\n' {
	//		col = 0
	//		row++
	//	} else if c == '<' {
	//		m.Set(I(int64(row)), I(int64(col)/3), I(1))
	//	}
	//	col++
	//}

	return
}

func LevelFromString(level string) (m Matrix, balls1 []Pt, balls2 []Pt) {
	// This is the kind of string that can get turned into a level.
	//	level = `
	//xxxxxxxxxxxxx
	//x           x
	//x  x     x  x
	//x           x
	//x           x
	//x    x  x   x
	//x           x
	//x           x
	//x           x
	//x         xxx
	//x           x
	//x           x
	//x           x
	//x           x
	//xxxxxxxxxxxxx
	//	`

	row := 0
	col := 0
	maxCol := 0
	for i := 0; i < len(level); i++ {
		c := level[i]
		if c == '`' {
			continue
		} else if c == '\n' {
			maxCol = col
			col = 0
			row++
			continue
		}
		col++
	}
	// If the string does not end with an empty line, count the last row.
	if col > 0 {
		row++
	}
	m.Init(I(row), I(maxCol))

	row = 0
	col = 0
	for i := 0; i < len(level); i++ {
		c := level[i]
		if c == '`' {
			continue
		} else if c == '\n' {
			col = 0
			row++
			continue
		} else if c == 'x' {
			m.Set(I(row), I(col), I(1))
		} else if c == '1' {
			balls1 = append(balls1, IPt(col, row))
		} else if c == '2' {
			balls2 = append(balls2, IPt(col, row))
		}
		col++
	}
	return
}
