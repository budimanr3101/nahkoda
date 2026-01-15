package semantic

import (
	"nahkoda/internal/parser"
)

// IntentResolver is the interface for resolving a specific action's intent.
type IntentResolver interface {
	Resolve(ast parser.AST, intent *Intent) error
}
