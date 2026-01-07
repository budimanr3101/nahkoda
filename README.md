# ⚓ Nahkoda v0.4.0

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

## 🧠 Filosofi v0.4.0 — *Koneksi Kubectl*

Mulai **v0.4.0**, Nahkoda terintegrasi langsung dengan `kubectl`.

Artinya:

* ❌ Kata tidak dikenal → **ERROR**
* ❌ Aksi tidak valid → **ERROR**
* ❌ Ambigu → **ERROR**
* ✅ Default boleh ada, tapi **harus jujur**
* ✅ Tidak ada silent fallback

Contoh:

```bash
nahkoda liat kru xyz
❌ kata tidak dikenali: "xyz"
```

```bash
nahkoda terbangkan kapal
❌ aksi tidak dikenali
```

Ini desain **sengaja ketat**, supaya:

* CLI bisa diprediksi
* Aman dipakai automation
* Engineer percaya output-nya

---

## 🧠 Cara Kerja Nahkoda

Nahkoda bekerja seperti **mini compiler**:

```
Teks Perintah
   ↓
Parser (AST)
   ↓
Semantic Resolver (strict)
   ↓
Intent
   ↓
Planner
   ↓
Executor (kubectl – planned)
```

Tidak ada:

* logika kubectl di parser
* asumsi tersembunyi
* output ambigu

---

## 🗣️ Vocabulary Nahkoda

Nahkoda menggunakan analogi kapal agar lebih natural:

| Nahkoda   | Kubernetes          |
| --------- | ------------------- |
| kru       | pod                 |
| geladak   | namespace           |
| rusak     | status != Running   |
| sehat     | status = Running    |
| bocor     | OOMKilled           |
| terdampar | Pending *(planned)* |

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
Filter : status=Running
```

(Default ini **eksplisit & jujur**)

---

### Lihat kru bermasalah

```bash
nahkoda liat kru rusak
```

```
Filter : status!=Running
```

---

### Dengan lokasi (namespace)

```bash
nahkoda liat kru bocor di geladak auth
```

```
Namespace : auth
Filter    : reason=OOMKilled
```

---

### Cek detail kru (Describe)

```bash
nahkoda cek kru payments-pod-1
```

```
Operation : describe
Resource  : pod
Target    : payments-pod-1
```

---

## 🚫 Contoh Error (STRICT)

```bash
nahkoda liat kru xyz
❌ kata tidak dikenali: "xyz"
```

```bash
nahkoda terbangkan kapal
❌ aksi tidak dikenali
```

Tidak ada guessing.
Tidak ada asumsi.

---

## 🧩 Arsitektur Internal

Struktur dibuat modular & scalable:

* **parser/**
  Tokenizer & AST

* **semantic/**
  Resolver intent, default, dan validasi

* **planner/**
  Translasi intent → execution plan

* **exec/**
  Executor (simulasi, kubectl planned)

---

## 📂 Struktur Repository

```
nahkoda/
├── cmd/              # CLI entrypoint (cobra)
├── internal/
│   ├── parser/       # Tokenizer & AST
│   ├── semantic/     # Strict semantic resolver
│   ├── planner/      # Execution planning
│   └── exec/         # Executor (simulasi)
├── TODO.md
├── README.md
└── main.go
```

---

## 🧪 Testing

Fokus utama ada di **semantic layer**:

```bash
go test ./internal/semantic -v
```

Jika semantic benar, seluruh CLI behavior konsisten.

---

## 🛣️ Roadmap (Next)

Lihat detail di [`TODO.md`](./TODO.md):

* kubectl executor
* typo handling (`krue` → `kru`)
* ambiguity resolver
* explain mode (`nahkoda explain`)
* shell completion

---

## 🚧 Status Proyek

* ✅ DSL v1 stabil
* ✅ Strict semantic resolver (v0.3.0)
* ✅ Deterministic output
* ✅ kubectl executor (v0.4.0)

Nahkoda **sudah aman dipamerkan**,
dan **fondasinya siap untuk production-grade CLI**.

---

## 📜 License

MIT License

---

> Nahkoda tidak menggantikan kubectl.
> Ia membuat kubectl **lebih manusiawi**.

⚓ **Happy sailing, Captain.**
