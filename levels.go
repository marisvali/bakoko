package bakoko

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
	//	level = `
	//<|><|><|><|><|><|><|><|><|><|><|><|><|><|><|>
	//<|>             <|>                       <|>
	//<|>             <|>                       <|>
	//<|>       <|><|><|><|><|><|>              <|>
	//<|>                                       <|>
	//<|>                                       <|>
	//<|>                              <|>      <|>
	//<|>               <|>            <|>      <|>
	//<|>                              <|>      <|>
	//<|>                              <|>      <|>
	//<|>            <|><|>                     <|>
	//<|>                                       <|>
	//<|>                                       <|>
	//<|>                                       <|>
	//<|><|><|><|><|><|><|><|><|><|><|><|><|><|><|>
	//`

	level = `
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
x   x                         x      x                 xx       x                x      x
x   x     xx              x   x      x                  x       x    xxxxxxx     x  x   x
x        xxx     x         x x       x          xx      x       x     xxxxxx     x      x
xxx              x          x        x         xx       x       x                x  x   x
x       x        x     xxxxxxxxxxxx   xx      xx        x  x    xxxxx      xx           x
x       xx       x xxxxxxxxxx           xx  xx             x             xxxxx     x xx x
xxx              x     x          x       xx               x   xx         xxxx          x
x      x         x     x          x             xxxxxxxxxxxx    x    xx            xxxxxx
x     xxx                    xxxx x                        x                  xxx       x
x      x         x           x    x          xxx      x    xxxxxxxxxxxxxxxxx            x
xxx             xxx     x    x    x    x    xxxx    xxx    x                            x
x    xxxx        x                x    x            x x         xxx  x                  x
x      xx   x         x           x    xxxx         x             xxxx   xxxxxx         x
xxx    xx   x       xxx           x                                       xxxxx         x
x           x       x   x   xx    x    x        xxxxxxxxxxxxxxxxxxx                     x
x           x           x   xx         x                                                x
x           x           x   xx     xxxxxxxxx      x    x     x     x   x   x    x  x    x
x           x       x                             x    xxxxxxx     x   x   x    x  x    x
x           x       x                   xxxxxx    x         x      x   x   x    x  x    x
x           x     xxxxx                 xxxxx     x         x      x   x   x    x  x    x
x       x                 xxxxxxxxxxx   xxxxx     x         x                           x
x     xxx    xx                                             x      xxxxxxxxxxxxxxxx     x
x   xxxxxx       xxx    xxxx                      xxxx      x                           x
x     xxx               x       x         x       xxxx      x      xxxxxxxxxxxxxxxx     x
x       x        xxxxxxxx       x         x       xxxx      x                           x
x                               x         x       xxxx      x      xxxxxxxxxxxxxxxx     x
x    xxxxxxxxxxxxxxxxxxxxxxxxxxxx         x                 x                           x
x                               x         x          xx     x      xxxxxxxxxxxxxxxx     x
xxxxxxxxxxxxxxxxxx              x                   xxx     x                           x
x                      xxxxxxxxxxxxxxxxxxx        xxxxx            xx     xxxxxxxx      x
x      xxxx                     x                xxxxxx              xxxxx         xx   x
x     x    xx                                               xxxxxx                  x   x
x     x      xx    xxxx                    xx  xx           xxxxxxx         x  xx   x   x
x      x             x          x   xxx                                     x  xx       x
x    x  x            x   x  x   x   xxx    xx  xx                           x       x   x
x   x    x           x          x   xxx    xx  xx                           x  xx   x   x
x  x                 x          x                              x         xxxx  xx   x   x
x     xxxxxxxx       x   x  x   x    x     xxxxxxx             x            x       x   x
x                    x          x    x     xxxxxxx             x            x    xxxxxxxx
x         x                     x    x                         x            x           x
x   x      xxx         x x      x     xxxxxx    x     xxxxxxxxxxxxxxxxxxx   x           x
x   x                           x               x              x            x    xxxxxxxx
x   x      xxxxxxxx    x x      x       xxx     x              x    x x     x           x
x   x      xxxxxxxx             x         x     xxxxxxxxxxx   x     x x     x    xxxxxxxx
x   x         xxxxx             x         x                    x            x           x
x                               x     xxxxxxxxxx     x               x                  x
x                                                    x               x                  x
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
`

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
		}
		col++
	}
	//m.Init(I(int64(row+1)), I(int64(maxCol/3+1)))
	m.Init(I(int64(row+1)), I(int64(maxCol+1)))

	row = -1
	col = 0
	for i := 0; i < len(level); i++ {
		c := level[i]
		if c == '`' {
			continue
		} else if c == '\n' {
			col = 0
			row++
		} else if c == 'x' {
			m.Set(I(int64(row)), I(int64(col)), I(1))
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
