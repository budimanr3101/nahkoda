package completer

import (
	"os/exec"
	"strings"

	"github.com/c-bata/go-prompt"
)

var actionNames = map[string]bool{
	"liat":   true,
	"cek":    true,
	"atur":   true,
	"baca":   true,
	"tukar":  true,
	"bikin":  true,
	"pantau": true,
	"hapus":  true,
	"masuk":  true,
}

var actions = []prompt.Suggest{
	{Text: "liat", Description: "Lihat daftar resource (get)"},
	{Text: "cek", Description: "Lihat detail resource (describe)"},
	{Text: "atur", Description: "Atur jumlah replika (scale)"},
	{Text: "baca", Description: "Baca jurnal/log (logs)"},
	{Text: "tukar", Description: "Restart armada/kru (rollout restart)"},
	{Text: "bikin", Description: "Buat resource baru (create)"},
	{Text: "pantau", Description: "Pantau penggunaan resource (top)"},
	{Text: "hapus", Description: "Hapus resource (delete)"},
	{Text: "masuk", Description: "Masuk ke dalam kru (exec)"},
	{Text: "keluar", Description: "Keluar dari Nahkoda"},
}

var objectNames = map[string]bool{
	"kru":        true,
	"mesin":      true,
	"armada":     true,
	"pelabuhan":  true,
	"mercusuar":  true,
	"peta":       true,
	"sandi":      true,
	"geladak":    true,
	"penjaga":    true,
	"perbekalan": true,
	"kesehatan":  true,
	"jurnal":     true,
	"berita":     true,
}

var objects = []prompt.Suggest{
	{Text: "kru", Description: "Pod"},
	{Text: "mesin", Description: "Node"},
	{Text: "armada", Description: "Deployment"},
	{Text: "pelabuhan", Description: "Service"},
	{Text: "mercusuar", Description: "Ingress"},
	{Text: "peta", Description: "ConfigMap"},
	{Text: "sandi", Description: "Secret"},
	{Text: "geladak", Description: "Namespace"},
	{Text: "penjaga", Description: "DaemonSet"},
	{Text: "perbekalan", Description: "Resource Limits/Requests"},
	{Text: "kesehatan", Description: "Audit kesehatan klaster"},
	{Text: "jurnal", Description: "Log harian"},
	{Text: "berita", Description: "Events"},
}

var keywords = []prompt.Suggest{
	{Text: "di", Description: "Preposisi lokasi"},
	{Text: "geladak", Description: "Namespace"},
	{Text: "terus", Description: "Follow logs (-f)"},
	{Text: "kabin", Description: "Container spesifik (-c)"},
	{Text: "ke", Description: "Target nominal"},
}

func Completer(d prompt.Document) []prompt.Suggest {
	return GetSuggestions(d.TextBeforeCursor(), d.GetWordBeforeCursor())
}

// GetSuggestions is a testable pure function for command completion
func GetSuggestions(textBefore, wordBefore string) []prompt.Suggest {
	if wordBefore == "" {
		// If at start or after a space, suggest based on context
		args := strings.Fields(textBefore)
		if len(args) == 0 {
			return actions
		}

		lastWord := args[len(args)-1]
		if isAction(lastWord) {
			return objects
		}

		if isObject(lastWord) {
			return keywords
		}

		if lastWord == "di" {
			return []prompt.Suggest{{Text: "geladak", Description: "Namespace (Wajib)"}}
		}

		if lastWord == "geladak" {
			return getDynamicNamespaces()
		}

		if (lastWord == "kru" || lastWord == "jurnal" || lastWord == "cek" || lastWord == "masuk") && !isAction(lastWord) {
			return getDynamicPods()
		}

		return keywords
	}

	return prompt.FilterHasPrefix(getAllSuggestions(), wordBefore, true)
}

func getDynamicNamespaces() []prompt.Suggest {
	out, err := exec.Command("kubectl", "get", "ns", "-o", "jsonpath={.items[*].metadata.name}").Output()
	if err != nil {
		return nil
	}

	nsList := strings.Fields(string(out))
	var suggestions []prompt.Suggest
	for _, ns := range nsList {
		suggestions = append(suggestions, prompt.Suggest{Text: ns, Description: "Geladak (Namespace)"})
	}
	return suggestions
}

func getDynamicPods() []prompt.Suggest {
	// For simplicity, we just get pods from default namespace or all-namespaces
	// In a real scenario, we might want to look at the previous 'geladak' keyword
	out, err := exec.Command("kubectl", "get", "pods", "-A", "-o", "jsonpath={.items[*].metadata.name}").Output()
	if err != nil {
		return nil
	}

	podList := strings.Fields(string(out))
	var suggestions []prompt.Suggest
	for _, pod := range podList {
		suggestions = append(suggestions, prompt.Suggest{Text: pod, Description: "Kru (Pod)"})
	}
	return suggestions
}

func isAction(word string) bool {
	return actionNames[word]
}

func isObject(word string) bool {
	return objectNames[word]
}

func getAllSuggestions() []prompt.Suggest {
	res := append(actions, objects...)
	res = append(res, keywords...)
	return res
}
