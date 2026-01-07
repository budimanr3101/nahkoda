package planner

import (
	"strings"

	"nahkoda/internal/semantic"
)

func Build(intent semantic.Intent) Plan {
	plan := Plan{
		Filters: make(map[string]string),
	}

	// ===== AKSI → OPERATION =====
	switch intent.Aksi {
	case "liat":
		plan.Operation = "list"
	case "hapus":
		plan.Operation = "delete"
	default:
		plan.Operation = "unknown"
		plan.Notes = append(plan.Notes, "aksi tidak dikenali")
	}

	// ===== OBJEK → RESOURCE =====
	switch intent.Objek {
	case "kru":
		plan.Resource = "pod"
	case "node":
		plan.Resource = "node"
	default:
		plan.Resource = "unknown"
		plan.Notes = append(plan.Notes, "objek tidak dikenali")
	}

	// ===== LOKASI → NAMESPACE =====
	if intent.Lokasi == "semua geladak" {
		plan.Namespace = "all"
	} else {
		plan.Namespace = normalizeNamespace(intent.Lokasi)
	}

	// ===== FILTER =====
	if intent.Filter != "" {
		key, value := splitFilter(intent.Filter)
		plan.Filters[key] = value
	}

	return plan
}

func normalizeNamespace(lokasi string) string {
	// "geladak auth" → "auth"
	parts := strings.Split(lokasi, " ")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return "default"
}

func splitFilter(filter string) (string, string) {
	// "status!=Running" atau "status=Running"
	if strings.Contains(filter, "!=") {
		parts := strings.Split(filter, "!=")
		return parts[0], "!=" + parts[1]
	}
	if strings.Contains(filter, "=") {
		parts := strings.Split(filter, "=")
		return parts[0], parts[1]
	}
	return filter, ""
}
