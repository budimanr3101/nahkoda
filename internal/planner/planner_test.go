package planner

import (
	"testing"

	"nahkoda/internal/semantic"
)

func TestBuild_Operations(t *testing.T) {
	tests := []struct {
		name   string
		intent semantic.Intent
		wantOp string
	}{
		{
			name:   "liat maps to get",
			intent: semantic.Intent{Aksi: "liat", Objek: "kru"},
			wantOp: "get",
		},
		{
			name:   "cek maps to describe",
			intent: semantic.Intent{Aksi: "cek", Objek: "kru", Target: "pod-1"},
			wantOp: "describe",
		},
		{
			name:   "hapus maps to delete",
			intent: semantic.Intent{Aksi: "hapus", Objek: "kru"},
			wantOp: "delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := Build(tt.intent)
			if plan.Operation != tt.wantOp {
				t.Errorf("Operation = %v, want %v", plan.Operation, tt.wantOp)
			}
		})
	}
}

func TestBuild_Resources(t *testing.T) {
	tests := []struct {
		name         string
		objek        string
		wantResource string
	}{
		{"kru maps to pod", "kru", "pod"},
		{"mesin maps to node", "mesin", "node"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intent := semantic.Intent{Aksi: "liat", Objek: tt.objek}
			plan := Build(intent)
			if plan.Resource != tt.wantResource {
				t.Errorf("Resource = %v, want %v", plan.Resource, tt.wantResource)
			}
		})
	}
}

func TestBuild_Namespaces(t *testing.T) {
	tests := []struct {
		name          string
		lokasi        string
		wantNamespace string
	}{
		{
			name:          "semua geladak maps to all",
			lokasi:        "semua geladak",
			wantNamespace: "all",
		},
		{
			name:          "geladak auth maps to auth",
			lokasi:        "geladak auth",
			wantNamespace: "auth",
		},
		{
			name:          "geladak default maps to default",
			lokasi:        "geladak default",
			wantNamespace: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intent := semantic.Intent{
				Aksi:   "liat",
				Objek:  "kru",
				Lokasi: tt.lokasi,
			}
			plan := Build(intent)
			if plan.Namespace != tt.wantNamespace {
				t.Errorf("Namespace = %v, want %v", plan.Namespace, tt.wantNamespace)
			}
		})
	}
}

func TestBuild_FilterToGrep(t *testing.T) {
	tests := []struct {
		name           string
		filter         string
		wantGrep       string
		wantGrepRegex  bool
		wantGrepInvert bool
	}{
		{
			name:           "status=Running",
			filter:         "status=Running",
			wantGrep:       "Running",
			wantGrepRegex:  false,
			wantGrepInvert: false,
		},
		{
			name:           "status!=Running",
			filter:         "status!=Running",
			wantGrep:       "Running",
			wantGrepRegex:  false,
			wantGrepInvert: true,
		},
		{
			name:           "status=Ready uses regex",
			filter:         "status=Ready",
			wantGrep:       `\bReady\b`,
			wantGrepRegex:  true,
			wantGrepInvert: false,
		},
		{
			name:           "status=Pending",
			filter:         "status=Pending",
			wantGrep:       "Pending",
			wantGrepRegex:  false,
			wantGrepInvert: false,
		},
		{
			name:           "status=NotReady",
			filter:         "status=NotReady",
			wantGrep:       "NotReady",
			wantGrepRegex:  false,
			wantGrepInvert: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intent := semantic.Intent{
				Aksi:   "liat",
				Objek:  "kru",
				Filter: tt.filter,
			}
			plan := Build(intent)

			if plan.Grep != tt.wantGrep {
				t.Errorf("Grep = %v, want %v", plan.Grep, tt.wantGrep)
			}

			if plan.GrepRegex != tt.wantGrepRegex {
				t.Errorf("GrepRegex = %v, want %v", plan.GrepRegex, tt.wantGrepRegex)
			}

			if plan.GrepInvert != tt.wantGrepInvert {
				t.Errorf("GrepInvert = %v, want %v", plan.GrepInvert, tt.wantGrepInvert)
			}
		})
	}
}

func TestBuild_Target(t *testing.T) {
	intent := semantic.Intent{
		Aksi:   "cek",
		Objek:  "kru",
		Target: "my-pod-123",
	}

	plan := Build(intent)

	if plan.Target != "my-pod-123" {
		t.Errorf("Target = %v, want %v", plan.Target, "my-pod-123")
	}
}

func TestBuild_NoFilter(t *testing.T) {
	intent := semantic.Intent{
		Aksi:  "liat",
		Objek: "mesin",
	}

	plan := Build(intent)

	if plan.Grep != "" {
		t.Errorf("Grep should be empty for no filter, got %v", plan.Grep)
	}

	if len(plan.Filters) != 0 {
		t.Errorf("Filters should be empty, got %v", plan.Filters)
	}
}

func TestNormalizeNamespace(t *testing.T) {
	tests := []struct {
		lokasi string
		want   string
	}{
		{"geladak auth", "auth"},
		{"geladak production", "production"},
		{"geladak default", "default"},
		{"invalid", "default"},
		{"", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.lokasi, func(t *testing.T) {
			got := normalizeNamespace(tt.lokasi)
			if got != tt.want {
				t.Errorf("normalizeNamespace(%v) = %v, want %v", tt.lokasi, got, tt.want)
			}
		})
	}
}

func TestSplitFilter(t *testing.T) {
	tests := []struct {
		filter    string
		wantKey   string
		wantValue string
	}{
		{"status=Running", "status", "=Running"},
		{"status!=Running", "status", "!=Running"},
		{"status=Pending", "status", "=Pending"},
		{"invalid", "invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.filter, func(t *testing.T) {
			key, value := splitFilter(tt.filter)
			if key != tt.wantKey {
				t.Errorf("key = %v, want %v", key, tt.wantKey)
			}
			if value != tt.wantValue {
				t.Errorf("value = %v, want %v", value, tt.wantValue)
			}
		})
	}
}
