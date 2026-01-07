package planner

type Plan struct {
	Operation string            // list | describe | delete
	Resource  string            // pod | node | ...
	Namespace string            // auth | default | all
	Target    string            // nama resource (untuk cek / describe)
	Filters   map[string]string // status=Running, dll
	Notes     []string          // warning / catatan
}
