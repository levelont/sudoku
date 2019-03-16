package sudoku

import (
	"fmt"
	"strings"
)

type board struct {
	cells [9][9]cell
	//TODO is it really cheaper to map to interface{}?
	rows [9]map[int]bool
	columns [9]map[int]bool
	regions [9]map[int]bool
}

func (b board) region(i,j int) *map[int]bool {
	var region *map[int]bool
	switch{
	case 0 <= i && i <= 2 &&
		0 <= j && j <= 2:
		region = &b.regions[0]
	case 0 <= i && i <= 2 &&
		3 <= j && j <= 5:
		region = &b.regions[1]
	case 0 <= i && i <= 2 &&
		6 <= j && j <= 8:
		region = &b.regions[2]
	case 3 <= i && i <= 5 &&
		0 <= j && j <= 2:
		region = &b.regions[3]
	case 3 <= i && i <= 5 &&
		3 <= j && j <= 5:
		region = &b.regions[4]
	case 3 <= i && i <= 5 &&
		6 <= j && j <= 8:
		region = &b.regions[5]
	case 6 <= i && i <= 8 &&
		0 <= j && j <= 2:
		region = &b.regions[6]
	case 6 <= i && i <= 8 &&
		3 <= j && j <= 5:
		region = &b.regions[7]
	case 6 <= i && i <= 8 &&
		6 <= j && j <= 8:
		region = &b.regions[8]
	}
	return region
}


//TODO can I pass an uninitialised values array?
// - how?
// - how can I warn about it?
func NewBoard(values [9][9]int) (board, error) {
	var b board

	//initialise helper maps
	for i := 0; i < 9; i++ {
		b.rows[i] = make(map[int]bool, 9)
		b.columns[i] = make(map[int]bool, 9)
		b.regions[i] = make(map[int]bool, 9)
	}

	//- initialise cells with provided values
	//- link each cell with helper maps
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			b.cells[i][j] = cell{
				value: values[i][j],
				row: &b.rows[i],
				column: &b.columns[j],
				region: b.region(i,j),
			}
		}
	}

	//fill helper maps
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			//if value is not 0
			//append value to row, column and region helper maps
			value := b.cells[i][j].value
			if value != 0 {
				if _, ok := (*b.cells[i][j].row)[value]; ok {
					return board{}, fmt.Errorf("duplicate value '%d' at row of [%d,%d]", value, i, j)
				}
				if _, ok := (*b.cells[i][j].column)[value]; ok {
					return board{}, fmt.Errorf("duplicate value '%d' at column of [%d,%d]", value, i, j)
				}
				if _, ok := (*b.cells[i][j].region)[value]; ok {
					return board{}, fmt.Errorf("duplicate value '%d' at region of [%d,%d]", value, i, j)
				}
				(*b.cells[i][j].row)[value] = true
				(*b.cells[i][j].column)[value] = true
				(*b.cells[i][j].region)[value] = true
			}
		}
	}

	return b, nil
}

//TODO test that setting a value updates the rows, columns and region values, both from the perspective of a cell and
// from the perspective of the board
func (b *board) set(i,j,value int) {
	b.cells[i][j].value = value
	b.rows[i][value] = true
	b.columns[j][value] = true
	(*b.region(i,j))[value] = true
}

type cell struct {
	value int
	set bool
	possibilities map[int]bool
	row *map[int]bool
	column *map[int]bool
	region *map[int]bool
}

func (b board) String() string {
	var s strings.Builder
	for _, x := range b.cells {
		for _, y := range x {
			if y.value != 0 {
				s.WriteString(fmt.Sprintf("%d ", y.value))
			} else{
				s.WriteString(fmt.Sprintf("%s ", " "))
			}
		}
		s.WriteString(fmt.Sprintln())
	}
	return s.String()
}
