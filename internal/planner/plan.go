package planner

type Plan struct {
	Operation  string            // list | describe | delete
	Resource   string            // pod | node | ...
	Namespace  string            // auth | default | all
	Target     string            // nama resource (untuk cek / describe)
	Filters    map[string]string // server-side filters
	Grep       string            // client-side filter (text match)
	GrepRegex  bool              // jika true, gunakan Regex Match
	GrepInvert bool              // jika true, tampilkan yang TIDAK match
	Notes      []string          // warning / catatan
}
