package completer

import (
	"context"
	"os/exec"
	"strings"
	"sync"
	"time"

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
	args := strings.Fields(textBefore)

	// Case 1: Start of command
	if len(args) == 0 {
		return actions
	}

	// Case 2: Handing specific word progress (wordBefore is not empty)
	if wordBefore != "" {
		hasAction := false
		for _, a := range args {
			if isAction(a) && a != wordBefore {
				hasAction = true
				break
			}
		}

		var all []prompt.Suggest
		if !hasAction {
			all = append(all, actions...)
		}
		all = append(all, objects...)
		all = append(all, keywords...)

		// If typing something that might be a resource name (pod/ns), add dynamic ones to filter
		if len(args) > 1 {
			lastArg := args[len(args)-2]
			if lastArg == "geladak" {
				all = append(all, getDynamicNamespaces()...)
			} else if lastArg == "kru" || lastArg == "jurnal" {
				all = append(all, getDynamicPods()...)
			} else if lastArg == "armada" {
				all = append(all, getDynamicDeployments()...)
			}
		}
		return prompt.FilterHasPrefix(all, wordBefore, true)
	}

	// Case 3: Word is finished (wordBefore is empty), deciding NEXT word based on context
	lastWord := args[len(args)-1]

	// 3.1 Adaptive Action -> Object Mapping
	switch lastWord {
	case "baca":
		return []prompt.Suggest{{Text: "jurnal", Description: "Log harian (logs)"}}
	case "masuk":
		return []prompt.Suggest{{Text: "kru", Description: "Pod (exec)"}}
	case "atur":
		return []prompt.Suggest{{Text: "armada", Description: "Deployment (scale)"}}
	case "tukar":
		return []prompt.Suggest{
			{Text: "kru", Description: "Pod (rollout)"},
			{Text: "armada", Description: "Deployment (rollout)"},
		}
	}

	// 3.2 Action -> General Objects
	if isAction(lastWord) {
		return objects
	}

	// 3.3 Grammar Sequence: di -> geladak
	if lastWord == "di" {
		return []prompt.Suggest{{Text: "geladak", Description: "Namespace (Wajib)"}}
	}

	// 3.4 Resource Discovery (Dynamic) + Keywords
	var res []prompt.Suggest
	if lastWord == "geladak" {
		res = getDynamicNamespaces()
	} else if lastWord == "kru" || lastWord == "jurnal" {
		res = getDynamicPods()
	} else if lastWord == "armada" {
		res = getDynamicDeployments()
	} else if lastWord == "pelabuhan" {
		res = getDynamicServices()
	}

	// Case 3.5: If we have an object or a resource name, suggest keywords as well
	if isObject(lastWord) || len(res) > 0 || (!isAction(lastWord) && !isKeyword(lastWord)) {
		hasDi := false
		for _, a := range args {
			if a == "di" {
				hasDi = true
				break
			}
		}

		if !hasDi {
			res = append(res, keywords...)
		} else {
			res = append(res, []prompt.Suggest{
				{Text: "terus", Description: "Follow logs (-f)"},
				{Text: "kabin", Description: "Container spesifik (-c)"},
			}...)
		}
		return res
	}

	return res
}

var (
	cache      = make(map[string][]prompt.Suggest)
	cacheTime  = make(map[string]time.Time)
	cacheMutex sync.RWMutex
	cacheTTL   = 30 * time.Second
)

func getFromCacheOrFetch(key string, fetcher func() []prompt.Suggest) []prompt.Suggest {
	cacheMutex.RLock()
	data, ok := cache[key]
	ts, okTime := cacheTime[key]
	cacheMutex.RUnlock()

	if ok && okTime && time.Since(ts) < cacheTTL {
		return data
	}

	// Double check locking for write
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Re-check in case another goroutine updated it
	if ts, okTime := cacheTime[key]; okTime && time.Since(ts) < cacheTTL {
		return cache[key]
	}

	// Fetch new data
	data = fetcher()
	if data != nil {
		cache[key] = data
		cacheTime[key] = time.Now()
	}
	return data
}

func init() {
	// Start background goroutine for cache cleanup
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			cleanupCache()
		}
	}()
}

func cleanupCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	now := time.Now()
	for key, ts := range cacheTime {
		if now.Sub(ts) > 5*time.Minute {
			delete(cache, key)
			delete(cacheTime, key)
		}
	}
}

var (
	currentContext     string
	currentContextTime time.Time
	contextCacheTTL    = 5 * time.Second
	contextMutex       sync.RWMutex
)

func getCurrentContext() string {
	contextMutex.RLock()
	if time.Since(currentContextTime) < contextCacheTTL && currentContext != "" {
		ctx := currentContext
		contextMutex.RUnlock()
		return ctx
	}
	contextMutex.RUnlock()

	// Fetch new context with timeout
	contextMutex.Lock()
	defer contextMutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "kubectl", "config", "current-context")
	out, err := cmd.Output()
	if err != nil {
		// Return cached on error, or default
		if currentContext != "" {
			return currentContext
		}
		return "default"
	}

	currentContext = strings.TrimSpace(string(out))
	currentContextTime = time.Now()
	return currentContext
}

func getCacheKey(resource string) string {
	return getCurrentContext() + ":" + resource
}

func getDynamicNamespaces() []prompt.Suggest {
	return getFromCacheOrFetch(getCacheKey("namespaces"), func() []prompt.Suggest {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "kubectl", "get", "ns", "-o", "jsonpath={.items[*].metadata.name}")
		out, err := cmd.Output()
		if err != nil {
			return nil
		}
		nsList := strings.Fields(string(out))
		var suggestions []prompt.Suggest
		for _, ns := range nsList {
			suggestions = append(suggestions, prompt.Suggest{Text: ns, Description: "Geladak (Namespace)"})
		}
		return suggestions
	})
}

func getDynamicPods() []prompt.Suggest {
	return getFromCacheOrFetch(getCacheKey("pods"), func() []prompt.Suggest {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "kubectl", "get", "pods", "-A", "-o", "jsonpath={.items[*].metadata.name}")
		out, err := cmd.Output()
		if err != nil {
			return nil
		}
		podList := strings.Fields(string(out))
		var suggestions []prompt.Suggest
		for _, pod := range podList {
			suggestions = append(suggestions, prompt.Suggest{Text: pod, Description: "Kru (Pod)"})
		}
		return suggestions
	})
}

func getDynamicDeployments() []prompt.Suggest {
	return getFromCacheOrFetch(getCacheKey("deployments"), func() []prompt.Suggest {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "kubectl", "get", "deployments", "-A", "-o", "jsonpath={.items[*].metadata.name}")
		out, err := cmd.Output()
		if err != nil {
			return nil
		}
		deployList := strings.Fields(string(out))
		var suggestions []prompt.Suggest
		for _, d := range deployList {
			suggestions = append(suggestions, prompt.Suggest{Text: d, Description: "Armada (Deployment)"})
		}
		return suggestions
	})
}

func getDynamicServices() []prompt.Suggest {
	return getFromCacheOrFetch(getCacheKey("services"), func() []prompt.Suggest {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "kubectl", "get", "svc", "-A", "-o", "jsonpath={.items[*].metadata.name}")
		out, err := cmd.Output()
		if err != nil {
			return nil
		}
		svcList := strings.Fields(string(out))
		var suggestions []prompt.Suggest
		for _, s := range svcList {
			suggestions = append(suggestions, prompt.Suggest{Text: s, Description: "Pelabuhan (Service)"})
		}
		return suggestions
	})
}

func isAction(word string) bool {
	return actionNames[word]
}

func isObject(word string) bool {
	return objectNames[word]
}

func isKeyword(word string) bool {
	for _, k := range keywords {
		if k.Text == word {
			return true
		}
	}
	return false
}
