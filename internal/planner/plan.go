package planner

type Plan struct {
	Operation string
	Resource  string
	Namespace string

	Filters map[string]string
	Notes   []string
}
