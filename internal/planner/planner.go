package planner

import (
	"strings"

	"nahkoda/internal/semantic"
)

func Build(intent semantic.Intent) Plan {
	plan := Plan{
		Filters: make(map[string]string),
	}

	// ===============================
	// 1️⃣ AKSI → OPERATION
	// ===============================
	switch intent.Aksi {
	// ===============================
	// BIKIN → CREATE / RUN
	// ===============================
	case "bikin":
		if intent.Objek == "geladak" {
			plan.Operation = "create"
			plan.Resource = "namespace"
			// args: create namespace [target]
		} else if intent.Objek == "kru" {
			plan.Operation = "run"
			plan.Resource = "" // run doesn't use "pod" resource mapping like get
			// args: run [target] --image=nginx --restart=Never
			plan.Flags = append(plan.Flags, "--image=nginx", "--restart=Never")
		}

	// ===============================
	// PANTAU → TOP
	// ===============================
	case "pantau":
		plan.Operation = "top"
		if intent.Objek == "kru" {
			plan.Resource = "pod"
		} else if intent.Objek == "mesin" {
			plan.Resource = "node"
		}

	case "liat":
		plan.Operation = "get"
	case "cek":
		plan.Operation = "describe"
	case "hapus":
		plan.Operation = "delete"
	case "pindah":
		plan.Operation = "config"
	case "baca":
		plan.Operation = "logs"
	case "masuk":
		plan.Operation = "exec"
	default:
		plan.Operation = "unknown"
		plan.Notes = append(plan.Notes, "aksi tidak dikenali")
	}

	// ===============================
	// 2️⃣ OBJEK → RESOURCE
	// ===============================
	switch intent.Objek {
	case "kru":
		plan.Resource = "pod"

	case "mesin":
		plan.Resource = "node"
	case "kapal":
		// Jika aksi "liat" -> get-contexts
		// Jika aksi "pindah" -> use-context
		if intent.Aksi == "liat" {
			plan.Operation = "config" // Override "get" from step 1
			plan.Resource = "get-contexts"
		} else if intent.Aksi == "pindah" {
			plan.Resource = "use-context"
		} else {
			plan.Resource = "unknown"
			plan.Notes = append(plan.Notes, "aksi tidak valid untuk kapal")
		}
	case "jurnal":
		plan.Resource = "pod" // kubectl logs [pod]
	case "berita":
		plan.Operation = "get"
		plan.Resource = "events"
		plan.Flags = append(plan.Flags, "--sort-by=.metadata.creationTimestamp")
	default:
		plan.Resource = "unknown"
		plan.Notes = append(plan.Notes, "objek tidak dikenali")
	}

	// ===============================
	// 3️⃣ LOKASI → NAMESPACE
	// ===============================
	if intent.Lokasi == "semua geladak" {
		plan.Namespace = "all"
	} else {
		plan.Namespace = normalizeNamespace(intent.Lokasi)
	}

	// ===============================
	// 4️⃣ TARGET (UNTUK CEK)
	// ===============================
	plan.Target = intent.Target

	// ===============================
	// 5️⃣ FILTER (LIST / GET)
	// ===============================
	if intent.Aksi == "liat" && intent.Filter != "" {
		key, value := splitFilter(intent.Filter)

		// MATCHING LOGIC (v0.5.0 & v0.6.0)
		// Jika filter adalah "status", kita gunakan Client-Side Grep
		if key == "status" {
			// value formats: "=Running" or "!=Running" or "=Pending" or "=Ready"
			val := value
			invert := false
			if strings.HasPrefix(val, "!=") {
				val = strings.TrimPrefix(val, "!=")
				invert = true
			} else {
				val = strings.TrimPrefix(val, "=")
			}

			// SPECIAL CASE: Ready vs NotReady (v0.6.0)
			// "Ready" match dengan "NotReady" jika pakai simple grep.
			// Solusi: Pakai Regex \bReady\b
			if val == "Ready" {
				plan.Grep = `\bReady\b`
				plan.GrepRegex = true
			} else {
				plan.Grep = val
				plan.GrepRegex = false
			}

			plan.GrepInvert = invert
		} else {
			// Filter lain (jika ada support di masa depan) tetap server-side
			plan.Filters[key] = value
		}
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
		return parts[0], "=" + parts[1]
	}
	return filter, ""
}
