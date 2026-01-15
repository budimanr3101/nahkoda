---
title: Instalasi
description: Cara memasang Nahkoda di sistem Anda.
---

## ⚙️ Konfigurasi (Opsional)

Nahkoda mendukung file konfigurasi untuk menyesuaikan perilaku aplikasi:

```bash
mkdir -p ~/.nahkoda
cat > ~/.nahkoda/config.json << 'EOF'
{
  "kubectl_path": "",
  "default_namespace": "default",
  "cache_ttl": 30000000000,
  "timeout": 30000000000,
  "enable_suggestions": true
}
EOF
```

**Catatan**: File config bersifat opsional. Jika tidak ada, Nahkoda menggunakan default values.

📖 **[Dokumentasi Lengkap Konfigurasi →](/teknis/konfigurasi)**

---

Nahkoda dapat dipasang melalui berbagai cara sesuai dengan sistem operasi yang Anda gunakan.

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
# Tambahkan bucket
scoop bucket add nahkoda https://github.com/budimanr3101/homebrew-nahkoda

# Install
scoop install nahkoda
```

### Manual Download

Anda bisa langsung download file binary-nya (ekstrak `.zip` atau `.tar.gz`) di halaman **[Releases](https://github.com/budimanr3101/nahkoda/releases)**.

### via Go (Manual Build)

```bash
git clone https://github.com/budimanr3101/nahkoda.git
cd nahkoda
go build -o nahkoda main.go
```
