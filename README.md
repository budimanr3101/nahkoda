# ⚓ Nahkoda

> **Human‑friendly command layer di atas kubectl, pakai bahasa Indonesia.**

Nahkoda adalah **CLI tool + Domain Specific Language (DSL)** yang bikin interaksi dengan Kubernetes jadi lebih santai, lebih kebaca, tapi tetap akurat.

Kalau biasanya kamu ngetik:

```bash
kubectl get pods -n auth --field-selector status.phase!=Running
```

Dengan Nahkoda, cukup:

```bash
nahkoda liat kru rusak di geladak auth
```

Output-nya tetap presisi. Bedanya, **otak kamu tidak ikut capek**.

---

## 🤔 Masalah yang Diselesaikan

Kubernetes itu powerful, tapi:

* Command panjang dan ribet
* Banyak flag sulit diingat
* Tidak ramah dibaca manusia

Nahkoda hadir sebagai **human‑friendly command layer di atas kubectl**.

Bukan untuk mengganti kubectl, tapi **menerjemahkan bahasa manusia ke bahasa mesin**.

---

## 🧠 Cara Berpikir Nahkoda

Nahkoda dibangun dengan prinsip sederhana:

* User bicara pakai bahasa natural
* Sistem yang mikir bagaimana eksekusinya
* Default boleh ada, tapi **tidak boleh bohong**

Secara internal, Nahkoda bekerja seperti compiler mini:

```
Teks Perintah → Parser → Semantic Resolver → Intent → (kubectl)
```

---

## 🗣️ Vocabulary Nahkoda

Nahkoda pakai istilah yang dekat dengan analogi kapal:

| Nahkoda   | Kubernetes        |
| --------- | ----------------- |
| kru       | pod               |
| geladak   | namespace         |
| rusak     | pod tidak Running |
| sehat     | pod Running       |
| bocor     | OOMKilled         |
| terdampar | Pending           |

Tujuannya satu: **perintah gampang dibaca tanpa kehilangan makna teknis**.

---

## ✨ Contoh Perintah

### Lihat kru sehat (default)

```bash
nahkoda liat kru
```

Output:

```
Filter : status=Running (aturan default: kru sehat)
```

---

### Lihat kru bermasalah

```bash
nahkoda liat kru rusak
```

```
Kondisi: rusak
Filter : status!=Running
```

---

### Pakai lokasi

```bash
nahkoda liat kru bocor di geladak auth
```

```
Lokasi : geladak auth
Kondisi: bocor
Filter : reason=OOMKilled
```

---

## 🧩 Arsitektur Internal

Struktur Nahkoda dibuat rapi dan bisa tumbuh:

* **Parser**
  Ngubah teks jadi AST (apa yang user tulis)

* **Semantic Resolver**
  Menentukan maksud, default, dan filter

* **Executor (planned)**
  Menerjemahkan intent ke kubectl

Tidak ada logika kubectl di parser. Tidak ada asumsi di output.

---

## 📂 Struktur Repository

```
nahkoda/
├── cmd/              # CLI entrypoint
├── internal/
│   ├── parser/       # Tokenizer & AST
│   ├── semantic/     # Resolver & condition mapping
│   └── executor/     # (planned) kubectl execution
├── TODO.md
├── README.md
└── main.go
```

---

## 🧪 Testing

Fokus utama ada di semantic layer:

```bash
go test ./internal/semantic -v
```

Kalau semantic-nya benar, output CLI pasti konsisten.

---

## 🛣️ Roadmap

Rencana selanjutnya bisa dilihat di [`TODO.md`](./TODO.md):

* Executor kubectl
* Explain mode (`nahkoda explain`)
* Shortcut syntax
* Pipeline support

---

## 🚧 Status Proyek

* ✅ DSL v1 stabil
* ✅ Semantic resolver solid
* ⏳ Executor sedang disiapkan

Nahkoda masih berkembang, tapi fondasinya sudah siap dipakai serius.

---

## 📜 Lisensi

MIT License

---

> Nahkoda tidak menggantikan kubectl.
> Ia membuat kubectl lebih enak dipakai.

⚓ Happy sailing dengan cluster kamu.
