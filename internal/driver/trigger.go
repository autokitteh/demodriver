package driver

import (
	"github.com/google/cel-go/cel"

	"go.autokitteh.dev/demodriver/internal/types"
)

type trigger struct {
	trigger types.Trigger
	filter  cel.Program
}
