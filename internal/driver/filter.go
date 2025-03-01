package driver

import (
	"fmt"

	"github.com/google/cel-go/cel"
)

var celEnv, _ = cel.NewEnv(
	cel.Variable("src", cel.StringType),
	cel.Variable("data", cel.AnyType),
)

func parseFilter(text string) (cel.Program, error) {
	if text == "" {
		return nil, nil
	}

	ast, issues := celEnv.Compile(text)
	if issues != nil && issues.Err() != nil {
		return nil, fmt.Errorf("compile: %w", issues.Err())
	}

	return celEnv.Program(ast)
}

func evalFilter(p cel.Program, src string, data any) (bool, error) {
	if p == nil {
		return true, nil
	}

	v, _, err := p.Eval(map[string]any{
		"src":  src,
		"data": data,
	})
	if err != nil {
		return false, fmt.Errorf("eval: %w", err)
	}

	result, ok := v.Value().(bool)
	if !ok {
		return false, fmt.Errorf("eval: unexpected type %T", v.Value())
	}

	return result, nil
}
