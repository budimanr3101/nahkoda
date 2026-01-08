[![Strict CI](https://github.com/budimanr3101/nahkoda/actions/workflows/ci.yaml/badge.svg)](https://github.com/budimanr3101/nahkoda/actions/workflows/ci.yaml)
# ⚓ Nahkoda v0.9.0

> **Human-friendly command layer di atas kubectl, pakai bahasa Indonesia.**

Nahkoda adalah **CLI tool + Domain Specific Language (DSL)** yang memungkinkan kamu berinteraksi dengan Kubernetes menggunakan **bahasa manusia**, tanpa kehilangan presisi teknis.

Kalau biasanya kamu ngetik:

```bash
kubectl get pods -n auth --field-selector status.phase!=Running
```

Dengan Nahkoda:

```bash
nahkoda liat kru rusak di geladak auth
```

Hasilnya tetap akurat.
Bedanya: **lebih kebaca, lebih santai, dan lebih manusiawi.**

---

## 🎯 Tujuan Proyek

Kubernetes itu powerful, tapi:

* Command panjang & sulit diingat
* Flag banyak & rawan typo
* Sulit dibaca manusia

Nahkoda hadir sebagai **lapisan semantic di atas kubectl**.

👉 **Bukan pengganti kubectl**,
👉 tapi **penerjemah bahasa manusia ke perintah Kubernetes**.

---

## 🚀 What's New in v0.7.0

### Production-Ready Improvements

#### 1. **Custom Error Types** 
- Structured error handling dengan context
- Type-safe error checking
- Better debugging capabilities
- 71.4% test coverage

#### 2. **Comprehensive Unit Tests**
- Parser: 100% coverage
- Semantic: 97.1% coverage
- Planner: 88.6% coverage
- **Total: 25 test suites, all passing**

#### 3. **Full Bahasa Indonesia**
- Semua error messages dalam Bahasa Indonesia
- Konsisten dari command sampai error output
- `perintah kubectl gagal` instead of `kubectl command failed`

#### 4. **Integration Testing**
- 32 comprehensive tests
- Real Kubernetes cluster testing
- Test resources (pods, deployments, namespaces)
- No tests skipped!

---

## 🧠 Filosofi Nahkoda

### Strict Semantic Mode

Nahkoda menggunakan **Client-Side Text Filtering** untuk hasil yang lebih manusiawi:

- `sehat` = Mengandung kata "Running"
- `rusak` = TIDAK mengandung kata "Running" (termasuk CrashLoopBackOff)
- `siap` = Regex match `\bReady\b` (presisi, menghindari NotReady)
- `mogok` = Status NotReady
- `terdampar` = Status Pending

**Graceful Error Handling:**
- Missing resource tidak error (exit 0)
- Input validation errors memberikan feedback, bukan panic
- Semua error messages dalam Bahasa Indonesia

**Strict Mode:**
- ❌ Kata tidak dikenal → **ERROR**
- ❌ Aksi tidak valid → **ERROR**
- ❌ Ambigu → **ERROR**
- ✅ Default boleh ada, tapi **harus jujur**
- ✅ Tidak ada silent fallback

Contoh:

```bash
nahkoda liat kru xyz
# kata tidak dikenali: "xyz"

nahkoda terbangkan kapal
# aksi tidak dikenali
```

Ini desain **sengaja ketat**, supaya:

* CLI bisa diprediksi
* Aman dipakai automation
* Engineer percaya output-nya

---

## 🗣️ Vocabulary Nahkoda

Nahkoda menggunakan analogi kapal agar lebih natural:

| Nahkoda   | Kubernetes          | Keterangan |
| --------- | ------------------- | ---------- |
| **Objek** |
| kru       | pod                 | Container workload |
| mesin     | node                | Cluster nodes |
| geladak   | namespace           | Logical separation |
| **Kondisi Pod** |
| rusak     | status != Running   | CrashLoop, ImagePull, Error |
| sehat     | status = Running    | Healthy pods |
| terdampar | status = Pending    | Unschedulable |
| **Kondisi Node** |
| siap      | status = Ready      | Node ready |
| mogok     | status = NotReady   | Node not ready |
| `liat mesin` | `kubectl get nodes` |
| `liat mesin siap` | `kubectl get nodes` (filter Ready) |
| `liat mesin mogok` | `kubectl get nodes -l status!=Ready` |
| `cek mesin [nama]` | `kubectl describe node [nama]` |
| `liat kapal` | `kubectl config get-contexts` |
| `pindah kapal [nama]` | `kubectl config use-context [nama]` |
| `baca jurnal [pod]` | `kubectl logs [pod]` |
| `masuk [pod]` | `kubectl exec -it [pod] -- /bin/sh` |
| `liat berita` | `kubectl get events --sort-by=.metadata.creationTimestamp` |
| **Aksi** |
| liat      | get / list          | View resources |
| cek       | describe            | Detailed info |
| hapus     | delete              | Delete resources |
| pindah    | config use-context  | Switch context |
| baca      | logs                | View logs |
| masuk     | exec -it            | Enter container |
| **Objek** |
| kru       | pod                 | Kubernetes Pod |
| mesin     | node                | Kubernetes Node |
| kapal     | context             | Kubernetes Context |
| jurnal    | logs                | Pod Logs |
| berita    | events              | Cluster Events |

Tujuannya:

> **perintah gampang dibaca, makna teknis tetap presisi**

---

## ✨ Contoh Perintah

### Lihat kru sehat (default)

```bash
nahkoda liat kru
```

Output:
```
⚓ Menjalankan: kubectl get pod -A | grep 'Running' (invert=false)
NAMESPACE     NAME                    READY   STATUS    RESTARTS   AGE
default       healthy-pod-1           1/1     Running   0          5m
auth          healthy-pod-2           1/1     Running   0          5m
```

(Default ini **eksplisit & jujur**)

---

### Lihat kru bermasalah

```bash
nahkoda liat kru rusak
```

Output:
```
⚓ Menjalankan: kubectl get pod -A | grep 'Running' (invert=true)
NAMESPACE     NAME                    READY   STATUS             RESTARTS   AGE
default       crashloop-pod           0/1     CrashLoopBackOff   5          2m
default       imagepull-pod           0/1     ImagePullBackOff   0          2m
```

---

### Dengan lokasi (namespace)

```bash
nahkoda liat kru sehat di geladak auth
```

Output:
```
⚓ Menjalankan: kubectl get pod -n auth | grep 'Running' (invert=false)
NAME            READY   STATUS    RESTARTS   AGE
healthy-pod-2   1/1     Running   0          5m
```

---

### Cek detail kru (Describe)

```bash
nahkoda cek kru healthy-pod-1
```

Output:
```
⚓ Menjalankan: kubectl describe pod healthy-pod-1 -n default
Name:         healthy-pod-1
Namespace:    default
...
```

---

### Lihat mesin (nodes)

```bash
nahkoda liat mesin siap
```

Output:
```
⚓ Menjalankan: kubectl get node -A | grep '\bReady\b' (invert=false)
NAME             STATUS   ROLES           AGE   VERSION
docker-desktop   Ready    control-plane   5d    v1.34.1
```

---

## 🚫 Contoh Error (STRICT)

```bash
nahkoda liat kru xyz
# kata tidak dikenali: "xyz"

nahkoda terbangkan kapal
# aksi tidak dikenali

nahkoda cek mesin
# cek mesin butuh nama mesin

nahkoda cek kru pod-tidak-ada
# ⚓ Menjalankan: kubectl describe pod pod-tidak-ada -n default
# Error from server (NotFound): pods "pod-tidak-ada" not found
# (No panic, graceful handling)
```

Tidak ada guessing.
Tidak ada asumsi.

---

## 🧩 Arsitektur Internal

Struktur dibuat modular & scalable:

```
Input (Bahasa Indonesia)
    ↓
Parser (AST) - 100% test coverage
    ↓
Semantic Resolver (strict) - 97.1% test coverage
    ↓
Intent
    ↓
Planner - 88.6% test coverage
    ↓
Executor (kubectl)
```

### Packages

* **parser/** - Tokenizer & AST
* **semantic/** - Resolver intent, default, dan validasi
* **planner/** - Translasi intent → execution plan
* **exec/** - Executor (kubectl dengan client-side filtering)
* **errors/** - Structured error types dengan context

---

## 📂 Struktur Repository

```
nahkoda/
├── cmd/              # CLI entrypoint (cobra)
├── internal/
│   ├── parser/       # Tokenizer & AST + tests
│   ├── semantic/     # Strict semantic resolver + tests
│   ├── planner/      # Execution planning + tests
│   ├── exec/         # Kubectl executor
│   └── errors/       # Custom error types + tests
├── docs/             # Documentation
│   ├── TESTING.md
│   ├── CUSTOM_ERRORS.md
│   ├── INDONESIANIZATION.md
│   └── PRODUCTION_IMPROVEMENTS.md
├── test-resources.yaml  # K8s test manifests
├── test.sh              # Comprehensive test suite
├── README.md
└── main.go
```

---

## 🧪 Testing

### Unit Tests

Fokus utama ada di **semantic layer**:

```bash
# Run all tests
go test ./... -v -cover

# Specific package
go test ./internal/parser -v
go test ./internal/semantic -v
go test ./internal/planner -v
```

**Coverage:**
- Parser: 100%
- Semantic: 97.1%
- Planner: 88.6%
- Errors: 71.4%

### Integration Tests

```bash
# Deploy test resources
kubectl apply -f test-resources.yaml

# Run comprehensive tests (32 tests)
./test.sh

# View results
cat nahkoda_test_v0.7.0.txt
```

Jika semantic benar, seluruh CLI behavior konsisten.

---

## 📦 Installation

```bash
# Clone repository
git clone https://github.com/yourusername/nahkoda.git
cd nahkoda

# Build
go build -o nahkoda

# Run
./nahkoda liat kru
```

---

## 🛣️ Roadmap

Lihat detail di [`TODO.md`](./TODO.md):

* ✅ kubectl executor (v0.4.0)
* ✅ Mature filtering (v0.5.0)
* ✅ Graceful error handling (v0.5.1-v0.5.2)
* ✅ Node status support (v0.6.0)
* ✅ Custom error types (v0.7.0)
* ✅ Unit tests (v0.7.0)
* ✅ Bahasa Indonesia errors (v0.7.0)
* ✅ Integration testing (v0.7.0)
* 🔄 Logging/debugging mode
* 🔄 Config file support
* 🔄 Typo handling (`krue` → `kru`)
* 🔄 Explain mode (`nahkoda explain`)
* 🔄 Shell completion

---

## 🚧 Status Proyek

* ✅ DSL v1 stabil
* ✅ Strict semantic resolver (v0.3.0)
* ✅ Deterministic output
* ✅ kubectl executor (v0.4.0)
* ✅ Mature filtering (v0.5.0)
* ✅ Graceful Error & Validation (v0.5.2)
* ✅ Node Status Support (v0.6.0)
* ✅ Custom Error Types (v0.7.0)
* ✅ Comprehensive Unit Tests (v0.7.0)
* ✅ Full Bahasa Indonesia (v0.7.0)
* ✅ Integration Testing (v0.7.0)

Nahkoda **production-ready**,
dan **fondasinya solid untuk enterprise CLI**.

---

## 📊 Test Statistics

```
Total Unit Tests:     25 suites
Total Integration:    32 tests
Code Coverage:        ~90% average
All Tests:            ✅ PASSING
Error Messages:       🇮🇩 100% Bahasa Indonesia
```

---

## 📜 License

MIT License

---

## 🙏 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Write tests for new features
4. Ensure all tests pass
5. Submit a pull request

---

## 📚 Documentation

- [Testing Guide](docs/TESTING.md)
- [Custom Errors](docs/CUSTOM_ERRORS.md)
- [Indonesianization](docs/INDONESIANIZATION.md)
- [Production Improvements](docs/PRODUCTION_IMPROVEMENTS.md)
- [Release Notes v0.6.0](RELEASE_NOTES_v0.6.0.md)

---

> Nahkoda tidak menggantikan kubectl.
> Ia membuat kubectl **lebih manusiawi**.

⚓ **Happy sailing, Captain.**
