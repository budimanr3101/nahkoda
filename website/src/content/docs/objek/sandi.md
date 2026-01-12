---
title: Sandi (Secret)
description: Menyimpan data rahasia dan kunci pengaman di dalam Sandi.
---

**Sandi** adalah metafora untuk **Secret** dalam Kubernetes.

### ⚓ Filosofi
Kenapa disebut Sandi? Karena setiap kapal memiliki informasi rahasia yang tidak boleh diketahui sembarang orang, seperti kunci brankas atau kode rahasia navigasi. Sandi menyimpan data sensitif (password, token, key) secara aman agar Kru bisa menggunakannya tanpa membocorkannya ke pihak luar.

### 📋 Ekuivalen Kubernetes
| Nahkoda | Kubernetes |
| :--- | :--- |
| **Sandi** | **Secret** |

### 🚀 Perintah Populer
Berikut adalah beberapa perintah yang bisa Anda jalankan terhadap **Sandi**:

- `nahkoda liat sandi` - Melihat daftar kunci pengaman yang tersedia.
- `nahkoda cek sandi [nama]` - Melakukan inspeksi terhadap metadata sandi (isinya tetap tersembunyi kecuali Anda sengaja membacanya).

### ⚠️ Yang Belum Bisa Dilakukan
Saat ini, Nahkoda belum mendukung:
- Membuat Sandi baru dari literal atau file.
- Melakukan decode base64 secara otomatis dalam tampilan `cek`.
- Mengelola TLS Secret secara khusus.
