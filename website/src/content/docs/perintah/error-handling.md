---
title: Sistem Error
description: Penjelasan teknis sistem error handling di Nahkoda.
---

Nahkoda menggunakan sistem error yang terstruktur untuk memberikan informasi yang jelas kepada pengguna (dan pengembang).

### Tipe Error

Berikut adalah daftar tipe error yang dikelola oleh Nahkoda:

| Tipe Error | Kategori | Deskripsi |
| :--- | :--- | :--- |
| `ErrUnknownWord` | Parser | Kata tidak dikenali dalam input. |
| `ErrUnknownAction` | Semantic | Aksi yang diminta tidak valid. |
| `ErrUnknownObject` | Semantic | Objek (kru, mesin, dsb) tidak valid. |
| `ErrMissingTarget` | Semantic | Target perintah (nama pod/node) belum diisi. |
| `ErrKubectlFailed` | Executor | Perintah `kubectl` yang asli gagal dieksekusi. |
| `ErrResourceNotFound` | Executor | Resource yang dicari tidak ada di cluster. |

### Mengapa Terstruktur?

1. **Human-Friendly**: Pesan kesalahan disesuaikan dengan konteks Nahkoda.
2. **Contextual**: Error dapat membawa data tambahan (seperti perintah kubectl yang gagal).
3. **Typo Suggestions**: Sistem error terintegrasi dengan mesin fuzzy matching untuk memberikan saran perbaikan.
