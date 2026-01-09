package semantic

// Keywords adalah kumpulan kata kunci yang dimengerti Nahkoda
var Keywords = []string{
	// Actions
	"liat", "hapus", "cek", "pindah", "baca", "masuk", "bikin", "pantau",
	// Objects
	"kru", "mesin", "kapal", "jurnal", "berita", "geladak",
	// Conditions
	"rusak", "bocor", "sehat", "terdampar", "siap", "mogok",
}

// FindSuggestion mencari kata kunci yang paling mirip dengan input
func FindSuggestion(input string) string {
	if len(input) < 2 {
		return ""
	}

	bestMatch := ""
	bestDistance := 99

	for _, kw := range Keywords {
		dist := LevenshteinDistance(input, kw)
		if dist < bestDistance {
			bestDistance = dist
			bestMatch = kw
		}
	}

	// Threshold: hanya sarankan jika jaraknya dekat (max 2 perubahan)
	// dan tidak lebih dari setengah panjang katanya
	if bestDistance <= 2 && bestDistance < len(bestMatch) {
		return bestMatch
	}

	return ""
}

// LevenshteinDistance menghitung jarak antara dua string
func LevenshteinDistance(s, t string) int {
	m := len(s)
	n := len(t)
	d := make([][]int, m+1)
	for i := range d {
		d[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		d[i][0] = i
	}
	for j := 1; j <= n; j++ {
		d[0][j] = j
	}

	for j := 1; j <= n; j++ {
		for i := 1; i <= m; i++ {
			cost := 1
			if s[i-1] == t[j-1] {
				cost = 0
			}
			d[i][j] = min(
				d[i-1][j]+1, // deletion
				min(d[i][j-1]+1, // insertion
					d[i-1][j-1]+cost)) // substitution
		}
	}

	return d[m][n]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
