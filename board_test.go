package sudoku

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initializeBasicBoard() [9][9]cell {
	var c [9][9]cell
	k := 1
	for i :=0; i<9;i++ {
		for j :=0; j<9;j++ {
			c[i][j] = cell{value:k}
			k++
		}
	}

	return c
}

func Test_board_String(t *testing.T) {
	type fields struct {
		cells   [9][9]cell
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "base",
			fields: fields{initializeBasicBoard()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := board{
				cells:   tt.fields.cells,
			}
			fmt.Printf(b.String())
		})
	}
}

func TestNewBoard(t *testing.T) {
	var tests = []struct {
		name string
		values [9][9]int
		expectError bool
		error error
	}{
		{
			name: "clean board",
			//Source: https://www.websudoku.com/?level=1
			values: [9][9]int{
				{0,6,0,  0,0,8,  0,9,0},
				{8,0,7,  1,0,2,  4,6,0},
				{0,4,0,  3,6,0,  0,0,5},

				{2,0,0,  6,8,1,  0,0,0},
				{0,0,0,  0,0,0,  0,0,0},
				{0,0,0,  9,4,3,  0,0,1},

				{5,0,0,  0,2,9,  0,1,0},
				{0,9,4,  8,0,6,  7,0,2},
				{0,2,0,  5,0,0,  0,3,0},
			},
		},
		{
			name: "duplicated value 5 at row of 4,7",
			values: [9][9]int{
				{0,6,0,  0,0,8,  0,9,0},
				{8,0,7,  1,0,2,  4,6,0},
				{0,4,0,  3,6,0,  0,0,5},

				{2,0,0,  6,8,1,  0,0,0},
				{0,5,0,  0,0,0,  0,5,0},
				{0,0,0,  9,4,3,  0,0,1},

				{5,0,0,  0,2,9,  0,1,0},
				{0,9,4,  8,0,6,  7,0,2},
				{0,2,0,  5,0,0,  0,3,0},
			},
			expectError: true,
			error: fmt.Errorf("duplicate value '%d' at row of [%d,%d]", 5,4,7),
		},
		{
			name: "duplicated value 9 at column of 7,7",
			values: [9][9]int{
				{0,6,0,  0,0,8,  0,9,0},
				{8,0,7,  1,0,2,  4,6,0},
				{0,4,0,  3,6,0,  0,0,5},

				{2,0,0,  6,8,1,  0,0,0},
				{0,0,0,  0,0,0,  0,0,0},
				{0,0,0,  9,4,3,  0,0,1},

				{5,0,0,  0,2,9,  0,1,0},
				{0,0,4,  8,0,6,  7,9,2},
				{0,2,0,  5,0,0,  0,3,0},
			},
			expectError: true,
			error: fmt.Errorf("duplicate value '%d' at column of [%d,%d]", 9,7,7),
		},
		{
			name: "duplicated value 1 at region of 5,8",
			values: [9][9]int{
				{0,6,0,  0,0,8,  0,9,0},
				{8,0,7,  1,0,2,  4,6,0},
				{0,4,0,  3,6,0,  0,0,5},

				{2,0,0,  6,8,0,  1,0,0},
				{0,0,0,  0,0,0,  0,0,0},
				{0,0,0,  9,4,3,  0,0,1},

				{5,0,0,  0,2,9,  0,1,0},
				{0,0,4,  8,0,6,  7,9,2},
				{0,2,0,  5,0,0,  0,3,0},
			},
			expectError: true,
			error: fmt.Errorf("duplicate value '%d' at region of [%d,%d]", 1,5,8),
		},
	}

	for _, test := range tests {

		b, err := NewBoard(test.values)
		if test.expectError {
			assert.EqualError(t, err, test.error.Error())
			continue
		}

		assert.NoError(t, err)

		//check values
		for i := 0; i<9; i++ {
			for j := 0; j<9; j++ {
				value := test.values[i][j]
				assert.Equal(t, value, (b.cells[i][j]).value, fmt.Sprintf("%s: position [%d,%d]", test.name, i, j))
				if value > 0 {
					_, rowOk := (*b.cells[i][j].row)[value]
					_, columnOk := (*b.cells[i][j].column)[value]
					_, regionOk := (*b.cells[i][j].region)[value]
					assert.True(t, rowOk, fmt.Sprintf("%s: rowOK for value %d at possition [%d,%d]", test.name, value, i, j))
					assert.True(t, columnOk, fmt.Sprintf("%s: columnOK for value %d at possition [%d,%d]", test.name, value, i, j))
					assert.True(t, regionOk, fmt.Sprintf("%s: regionOk for value %d at possition [%d,%d]", test.name, value, i, j))
				}
			}
		}
	}
}

func TestUpdatePossibilities(t *testing.T) {
	tests := []struct{
		name string
		cell cell
		expectedPossibilities map[int]bool
	}{
		{
			name: "1 in row, 2 in column, 3 in region",
			cell: cell{
				row: &map[int]bool{1: true},
				column: &map[int]bool{2: true},
				region: &map[int]bool{3: true},
				possibilities: make(map[int]bool),
			},
			expectedPossibilities: map[int]bool{
				4: true,
				5: true,
				6: true,
				7: true,
				8: true,
				9: true,
			},
		},
		{
			name: "row, column and region empty",
			cell: cell{
				row: &map[int]bool{},
				column: &map[int]bool{},
				region: &map[int]bool{},
				possibilities: make(map[int]bool),
			},
			expectedPossibilities: map[int]bool{
				1: true,
				2: true,
				3: true,
				4: true,
				5: true,
				6: true,
				7: true,
				8: true,
				9: true,
			},
		},
		{
			name: "all values set",
			cell: cell{
				row: &map[int]bool{1: true, 2: true, 3: true},
				column: &map[int]bool{4: true, 5: true, 6: true},
				region: &map[int]bool{7: true, 8: true, 9: true},
				possibilities: make(map[int]bool),
			},
			expectedPossibilities: make(map[int]bool),
		},
	}

	for _, test := range tests {
		test.cell.updatePossibilities()
		assert.Equal(t, test.expectedPossibilities, test.cell.possibilities, test.name)
	}
}
