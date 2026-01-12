---
title: Pelabuhan (Service)
description: Menghubungkan Kru melalui Pelabuhan yang stabil.
---

**Pelabuhan** adalah metafora untuk **Service** dalam Kubernetes.

### ⚓ Filosofi
Kenapa disebut Pelabuhan? Karena Kru (Pod) bersifat dinamis dan alamatnya bisa berubah-ubah. Pelabuhan menyediakan alamat yang tetap dan aman bagi kapal lain atau pihak luar untuk mengirimkan kargo (data) menuju Kru yang tepat. Tanpa Pelabuhan, komunikasi antar Kru akan kacau.

### 📋 Ekuivalen Kubernetes
| Nahkoda | Kubernetes |
| :--- | :--- |
| **Pelabuhan** | **Service** |

### 🚀 Perintah Populer
Berikut adalah beberapa perintah yang bisa Anda jalankan terhadap **Pelabuhan**:

- `nahkoda liat pelabuhan` - Melihat semua pintu masuk data internal yang tersedia.
- `nahkoda cek pelabuhan [nama]` - Melihat detail mapping port dan IP dari pelabuhan tersebut.

### ⚠️ Yang Belum Bisa Dilakukan
Saat ini, Nahkoda belum mendukung:
- Membuat Pelabuhan baru secara otomatis.
- Mengubah tipe Pelabuhan (misal dari ClusterIP ke NodePort).
- Mengatur Load Balancer eksternal melalui perintah Nahkoda.
