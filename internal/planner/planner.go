package planner

import (
	"strings"

	"nahkoda/internal/semantic"
)

func Build(intent semantic.Intent) Plan {
	plan := Plan{
		Filters: make(map[string]string),
	}

	// 1. AKSI → OPERATION
	switch intent.Aksi {
	case "bikin":
		if intent.Objek == "namespace" {
			plan.Operation = "create"
			plan.Resource = "namespace"
		} else if intent.Objek == "pod" {
			plan.Operation = "run"
			plan.Resource = ""
			plan.Flags = append(plan.Flags, "--image=nginx", "--restart=Never")
		} else {
			// standard bikin: create [resource] [target]
			plan.Operation = "create"
			plan.Resource = intent.Objek
		}

	case "pantau":
		plan.Operation = "top"
		plan.Resource = intent.Objek

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
		if intent.Follow {
			plan.Flags = append(plan.Flags, "-f")
		}
		if intent.SubTarget != "" {
			plan.Flags = append(plan.Flags, "-c", intent.SubTarget)
		}
	case "atur":
		plan.Operation = "scale"
	case "masuk":
		plan.Operation = "exec"
	case "tukar":
		plan.Operation = "rollout restart"
	default:
		plan.Operation = "unknown"
		plan.Notes = append(plan.Notes, "aksi tidak dikenali")
	}

	// 2. OBJEK → RESOURCE
	switch intent.Objek {
	case "kapal":
		if intent.Aksi == "liat" {
			plan.Operation = "config"
			plan.Resource = "get-contexts"
		} else if intent.Aksi == "pindah" {
			plan.Resource = "use-context"
		} else {
			plan.Resource = "unknown"
			plan.Notes = append(plan.Notes, "aksi tidak valid untuk kapal")
		}
	case "jurnal":
		plan.Resource = "pod"
	case "berita":
		plan.Operation = "get"
		plan.Resource = "events"
		plan.Flags = append(plan.Flags, "--sort-by=.metadata.creationTimestamp")
	default:
		// mapping sisa objek (pod, node, deployment, service, ingress, configmap, secret, daemonset, namespace, kesehatan)
		plan.Resource = intent.Objek
		if intent.Aksi == "cek" && intent.Objek == "kesehatan" {
			plan.Operation = "audit"
		}
		if intent.Aksi == "liat" && intent.Objek == "perbekalan" {
			// Perbekalan (Resource Management)
			plan.Operation = "get"
			// Mapping Nilai (resource type) -> Kubernetes Standard
			resMap := map[string]string{
				"kru":     "pod",
				"armada":  "deployment",
				"penjaga": "daemonset",
			}
			res := "pod" // default
			if r, ok := resMap[intent.Nilai]; ok {
				res = r
			}
			plan.Resource = res

			jsonPath := ""
			if res == "pod" {
				jsonPath = "custom-columns=NAME:.metadata.name,CONTAINERS:.spec.containers[*].name,CPU_REQ:.spec.containers[*].resources.requests.cpu,MEM_REQ:.spec.containers[*].resources.requests.memory,CPU_LIM:.spec.containers[*].resources.limits.cpu,MEM_LIM:.spec.containers[*].resources.limits.memory"
			} else {
				// deployment/daemonset have resources in template
				jsonPath = "custom-columns=NAME:.metadata.name,CONTAINERS:.spec.template.spec.containers[*].name,CPU_REQ:.spec.template.spec.containers[*].resources.requests.cpu,MEM_REQ:.spec.template.spec.containers[*].resources.requests.memory,CPU_LIM:.spec.template.spec.containers[*].resources.limits.cpu,MEM_LIM:.spec.template.spec.containers[*].resources.limits.memory"
			}
			plan.Flags = append(plan.Flags, "-o", jsonPath)
		}
	}

	// 3. LOKASI → NAMESPACE
	if intent.Lokasi == "semua geladak" {
		plan.Namespace = "all"
	} else {
		plan.Namespace = normalizeNamespace(intent.Lokasi)
	}

	// 4. TARGET
	plan.Target = intent.Target

	// 5. FILTER (LIST / GET)
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

	// 6. REPLICAS (UNTUK ATUR)
	if intent.Aksi == "atur" && intent.Nilai != "" {
		plan.Flags = append(plan.Flags, "--replicas="+intent.Nilai)
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
