package util

import "strings"

func ConvertCol(col_square byte) uint8 {
	col_square = strings.ToUpper(string(col_square))[0]
	return uint8(col_square - 'A')
}

func ConvertRow(row_square byte) uint8 {
	return 5 - (row_square - '1')
}

func ConvertColBack(col_index uint8) byte {
	return byte(col_index + 'A')
}

func ConvertSquare(square string) uint8 {
	col, row := square[0], square[1]
	colIndex, rowIndex := ConvertCol(col), ConvertRow(row)
	return colIndex + rowIndex*7
}
