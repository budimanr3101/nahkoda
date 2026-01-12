---
title: Kru (Pod)
description: Mengenal Kru, unit terkecil dalam pelayaran Nahkoda.
---

**Kru** adalah metafora untuk **Pod** dalam Kubernetes. 

### ⚓ Filosofi
Kenapa disebut Kru? Karena sebuah kapal tidak akan bisa berjalan tanpa awak kapal yang bekerja di dalamnya. Kru adalah unit terkecil yang melakukan pekerjaan nyata (container). Mereka bisa datang dan pergi, bisa sakit, atau bisa diganti dengan kru baru yang lebih sehat.

### 📋 Ekuivalen Kubernetes
| Nahkoda | Kubernetes |
| :--- | :--- |
| **Kru** | **Pod** |

### 🚀 Perintah Populer
Berikut adalah beberapa perintah yang bisa Anda jalankan terhadap **Kru**:

- `nahkoda liat kru` - Melihat semua kru di seluruh geladak.
- `nahkoda liat kru rusak` - Mencari kru yang sedang bermasalah atau tidak berjalan.
- `nahkoda cek kru [nama]` - Melihat detail kesehatan dan aktivitas seorang kru.
- `nahkoda baca jurnal [nama]` - Membaca catatan aktivitas (logs) dari kru.
- `nahkoda masuk [nama]` - Masuk ke dalam ruangan kerja kru (exec shell).
- `nahkoda hapus kru [nama]` - Meliburkan atau menghapus kru dari tugasnya.

### ⚠️ Yang Belum Bisa Dilakukan
Saat ini, Nahkoda belum mendukung:
- Melihat log dari container spesifik jika satu Kru (Pod) memiliki banyak container.
- Melakukan port-forwarding langsung lewat perintah `masuk`.
- Mengedit manifest Pod secara interaktif.
