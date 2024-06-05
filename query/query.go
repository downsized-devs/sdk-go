package query

import (
	"context"
	"fmt"
	"strings"

	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
)

// FormatQueryForRows I hate this, find a better way for insert many rows
func FormatQueryForRows(ctx context.Context, q string, inputs [][]interface{}) (string, []interface{}, error) {
	// Add () based on rows
	// Add ? based on cols
	lRow := len(inputs)
	if lRow < 1 {
		return "", nil, errors.NewWithCode(codes.CodeSQLPrepareStmt, "no inputs rows supplied")
	}
	lCol := len(inputs[0])
	if lCol < 1 {
		return "", nil, errors.NewWithCode(codes.CodeSQLPrepareStmt, "no inputs cols supplied")
	}

	ins := []string{}
	for i := 0; i < lCol; i++ {
		ins = append(ins, "?")
	}
	inputTemplate := fmt.Sprintf("(%s)", strings.Join(ins, ", "))

	insQuery := []string{}
	for i := 0; i < lRow; i++ {
		insQuery = append(insQuery, inputTemplate)
	}

	result := fmt.Sprintf("%s %s", q, strings.Join(insQuery, ", "))

	args := []interface{}{}
	for _, r := range inputs {
		args = append(args, r...)
	}

	return result, args, nil
}
