[![Strict CI](https://github.com/budimanr3101/nahkoda/actions/workflows/ci.yaml/badge.svg)](https://github.com/budimanr3101/nahkoda/actions/workflows/ci.yaml)

# ⚓ Nahkoda v1.0.0 ⚓🚀

> **Kubernetes dalam Bahasa Manusia.**

Nahkoda adalah CLI yang mengubah interaksi Kubernetes (kubectl) menjadi **Bahasa Indonesia** sehari-hari yang natural, ekspresif, dan tetap presisi. Sekarang dengan dukungan fitur **Bikin**, **Pantau**, dan **Toleransi Typo**.

📖 **[Dokumentasi Lengkap (Web)](https://budimanr3101.github.io/nahkoda/)**

---

## 📥 Instalasi

### via Homebrew (Recommended)

Cara termudah untuk pengguna macOS dan Linux.

```bash
# Tambahkan tap (repository)
brew tap budimanr3101/nahkoda

# Install
brew install nahkoda
```

### via Scoop (Windows)

```powershell
# Tambahkan bucket (diambil dari repo yang sama dengan brew)
scoop bucket add nahkoda https://github.com/budimanr3101/homebrew-nahkoda

# Install
scoop install nahkoda
```

### Manual Download (Windows/Linux/macOS)

Jika tidak ingin pakai Package Manager, Anda bisa langsung download file binary-nya (ekstrak `.zip` atau `.tar.gz`) di halaman **[Releases](https://github.com/budimanr3101/nahkoda/releases)**.

### via Go (Manual Build)

```bash
git clone https://github.com/budimanr3101/nahkoda.git
cd nahkoda
go build -o nahkoda main.go
mv nahkoda /usr/local/bin/
```

---

## 📖 Kamus Nahkoda

Nahkoda menggunakan metafora pelayaran untuk mempermudah pemahaman resource Kubernetes:

| Komponen | Nama Nahkoda | Filosofi |
| :--- | :--- | :--- |
| **Cluster/Context** | **Kapal** | Kapal besar yang kita kendalikan. |
| **Namespace** | **Geladak** | Lantai/ruangan spesifik di dalam kapal. |
| **Node** | **Mesin** | Sumber tenaga (server) yang menjalankan kapal. |
| **Pod** | **Kru** | Anggota yang melaksanakan pekerjaan. |
| **Deployment** | **Armada** | Kelompok besar kapal/kru yang bergerak bersama. |
| **DaemonSet** | **Penjaga** | Kru yang wajib ada di setiap unit mesin (Node). |
| **Service** | **Pelabuhan** | Titik temu logistik agar kru bisa dihubungi. |
| **Ingress** | **Mercusuar** | Pemandu trafik dari luar menuju pelabuhan. |
| **ConfigMap** | **Peta** | Panduan konfigurasi/instruksi bagi kru. |
| **Secret** | **Sandi** | Data rahasia dan kunci-kunci pengaman. |
| **Logs** | **Jurnal** | Catatan kegiatan harian kru. |
| **Events** | **Berita** | Kabar terbaru tentang kondisi kapal. |

---

## 🚀 Panduan Cepat

**Nahkoda** bukan pengganti `kubectl`, tapi "penerjemah" agar perintah lebih mudah diingat.

### 🗺️ Navigasi (Context)

| Perintah Nahkoda | Ekuivalen Kubectl | Fungsi |
| :--- | :--- | :--- |
| `nahkoda liat kapal` | `kubectl config get-contexts` | List semua cluster/context |
| `nahkoda pindah kapal [nama]` | `kubectl config use-context [nama]` | Pindah cluster aktif |

### 📦 Monitoring (Pods & Nodes)

| Perintah Nahkoda | Ekuivalen Kubectl | Fungsi |
| :--- | :--- | :--- |
| `nahkoda liat kru` | `kubectl get pods -A` | List semua pod di semua namespace |
| `nahkoda liat kru di geladak [ns]` | `kubectl get pods -n [ns]` | List pod di namespace tertentu |
| `nahkoda liat kru rusak` | `kubectl get pods ... \| grep -v Running` | Cari pod yang error/crash |
| `nahkoda liat mesin` | `kubectl get nodes` | List worker nodes |
| `nahkoda liat berita` | `kubectl get events --sort-by=...` | Lihat event cluster terbaru |

### 🔧 Debugging (Logs, Exec, Describe)

| Perintah Nahkoda | Ekuivalen Kubectl | Fungsi |
| :--- | :--- | :--- |
| `nahkoda baca jurnal [pod]` | `kubectl logs [pod]` | Baca logs dari pod |
| `nahkoda masuk [pod]` | `kubectl exec -it [pod] -- /bin/sh` | Masuk ke container (shell) |
| `nahkoda cek kru [pod]` | `kubectl describe pod [pod]` | Lihat detail/status pod |
| `nahkoda cek mesin [node]` | `kubectl describe node [node]` | Lihat detail node |

### 🛠️ Operation & Metrics

| Perintah Nahkoda | Ekuivalen Kubectl | Fungsi |
| :--- | :--- | :--- |
| `nahkoda liat armada` | `kubectl get deployment -A` | List semua armada (deployment) |
| `nahkoda liat pelabuhan` | `kubectl get service -A` | List semua pelabuhan (service) |
| `nahkoda liat mercusuar` | `kubectl get ingress -A` | List semua mercusuar (ingress) |
| `nahkoda liat penjaga` | `kubectl get daemonset -A` | List semua penjaga (daemonset) |
| `nahkoda liat peta` | `kubectl get configmap -A` | List semua peta (configmap) |
| `nahkoda liat sandi` | `kubectl get secret -A` | List semua sandi (secret) |
| `nahkoda bikin geladak [nama]` | `kubectl create namespace [nama]` | Buat namespace baru |
| `nahkoda bikin kru [nama]` | `kubectl run [nama] --image=nginx...` | Buat pod baru (nginx) |
| `nahkoda pantau kru` | `kubectl top pod` | Lihat penggunaan resource pod |
| `nahkoda pantau mesin` | `kubectl top node` | Lihat penggunaan resource node |
| `nahkoda hapus kru [nama]` | `kubectl delete pod [nama]` | Hapus pod |

---

## 🧠 Filosofi & Desain

### 1. **Analogi Pelayaran**
Kami menggunakan metafora kapal untuk membuat Kubernetes lebih *approachable*:
- **Kapal** = Cluster / Context
- **Geladak** (Deck) = Namespace
- **Kru** = Pod (pekerja)
- **Mesin** = Node (infrastruktur)
- **Jurnal** = Logs
- **Berita** = Events

### 2. **Strict tapi Ramah (Typo Tolerance)**
Nahkoda didesain untuk **tidak menebak-nebak secara liar**, tapi membantu Kapten yang sedang lelah.
- Jika ada typo seperti `liat kur`, Nahkoda akan bertanya: *"Mungkin maksud Kapten: **kru**? (y/n)"*.
- Jika resource tidak ada di namespace default, Nahkoda akan memberi **Tips** (*"Coba cek di geladak lain..."*).

### 3. **Native Go Implementation**
Versi ini (v0.9.0+) ditulis ulang menjadi **100% Pure Go** (tanpa dependensi `cobra` atau library berat). Ukuran binary sangat kecil (~2.9MB) dan cepat.

---

## 🛠️ Pengembangan

Struktur project ini modular dan mudah dipelajari:

```
nahkoda/
├── internal/
│   ├── parser/       # Tokenizer & AST (Memahami input teks)
│   ├── semantic/     # Resolver & Validasi (Memastikan perintah masuk akal)
│   ├── planner/      # Penerjemah ke perintah kubectl
│   ├── exec/         # Eksekutor perintah
│   └── errors/       # Sistem error handling bahasa manusia
├── test.sh           # Integration Test Suite
└── main.go           # Entrypoint
```

### Menjalankan Test
```bash
./test.sh
```

---

## 📄 Lisensi
MIT License.

> **"Happy Sailing, Captain!"** ⚓
