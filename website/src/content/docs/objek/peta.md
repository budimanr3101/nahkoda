---
title: Peta (ConfigMap)
description: Memberikan petunjuk konfigurasi bagi Kru melalui Peta.
---

**Peta** adalah metafora untuk **ConfigMap** dalam Kubernetes.

### ⚓ Filosofi
Kenapa disebut Peta? Karena Peta berisi informasi statis dan petunjuk jalan yang dibutuhkan oleh Kru untuk menjalankan tugasnya. Dengan Peta, Anda bisa mengubah perilaku Kru tanpa harus mengganti Kru itu sendiri. Cukup ubah isinya, dan Kru akan mengikuti petunjuk baru.

### 📋 Ekuivalen Kubernetes
| Nahkoda | Kubernetes |
| :--- | :--- |
| **Peta** | **ConfigMap** |

### 🚀 Perintah Populer
Berikut adalah beberapa perintah yang bisa Anda jalankan terhadap **Peta**:

- `nahkoda liat peta` - Melihat daftar konfigurasi yang tersedia di kapal.
- `nahkoda cek peta [nama]` - Membaca isi instruksi yang ada di dalam peta tersebut.

### ⚠️ Yang Belum Bisa Dilakukan
Saat ini, Nahkoda belum mendukung:
- Membuat Peta baru dari file lokal secara otomatis.
- Mengedit isi Peta secara langsung (inline editing).
- Menghubungkan Peta ke Kru secara otomatis.
