package num

import (
	"context"
	"strings"

	"github.com/xuri/excelize/v2"
)

// i.e. D10:D13 -> [D10, D11, D12, D13]
func ExcelGenerateCoords(ctx context.Context, coordinateRange string) ([]string, error) {
	result := []string{}
	coords := strings.Split(coordinateRange, ":")
	startCol, startRow, err := excelize.CellNameToCoordinates(coords[0])
	if err != nil {
		return result, err
	}
	endCol, endRow, err := excelize.CellNameToCoordinates(coords[1])
	if err != nil {
		return result, err
	}

	for c := startCol; c <= endCol; c++ {
		for r := startRow; r <= endRow; r++ {
			namedCoor, err := excelize.CoordinatesToCellName(c, r)
			if err != nil {
				return result, err
			}
			result = append(result, namedCoor)
		}
	}

	return result, nil
}
