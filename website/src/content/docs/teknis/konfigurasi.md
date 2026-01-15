---
title: Konfigurasi
description: Panduan lengkap mengkustomisasi Nahkoda via config.json
---

Nahkoda mendukung file konfigurasi untuk menyesuaikan perilaku aplikasi sesuai kebutuhan Anda.

## 📁 Lokasi File

File konfigurasi disimpan di:
```
~/.nahkoda/config.json
```

**Catatan**: File ini bersifat **opsional**. Jika tidak ada, Nahkoda akan menggunakan default values.

## ⚙️ Default Values

Jika tidak ada file config, Nahkoda menggunakan:

| Field | Default Value | Keterangan |
|-------|---------------|------------|
| `kubectl_path` | `""` (auto-detect) | Path ke binary kubectl |
| `default_namespace` | `"default"` | Namespace default untuk operasi |
| `cache_ttl` | `30000000000` (30 detik) | Durasi cache autocomplete |
| `timeout` | `30000000000` (30 detik) | Timeout untuk kubectl calls |
| `enable_suggestions` | `true` | Toggle autocomplete/suggestions |

## 📝 Format Config

```json
{
  "kubectl_path": "/usr/local/bin/kubectl",
  "default_namespace": "production",
  "cache_ttl": 60000000000,
  "timeout": 45000000000,
  "enable_suggestions": true
}
```

**Catatan**: `cache_ttl` dan `timeout` dalam **nanoseconds** (Go duration format).

### Konversi Waktu

| Durasi | Nanoseconds |
|--------|-------------|
| 1 detik | 1000000000 |
| 30 detik | 30000000000 |
| 1 menit | 60000000000 |
| 5 menit | 300000000000 |

## 🎯 Use Cases

### 1. Custom Kubectl Path

Berguna jika Anda menggunakan:
- Custom kubectl binary (e.g., `k3s kubectl`, `microk8s kubectl`)
- Kubectl di lokasi non-standard
- Multiple kubectl versions

```json
{
  "kubectl_path": "/snap/bin/microk8s.kubectl"
}
```

### 2. Default Namespace

Atur namespace default agar tidak perlu `di geladak` setiap kali:

```json
{
  "default_namespace": "staging"
}
```

Sebelum:
```bash
nahkoda liat kru di geladak staging
```

Sesudah (dengan config):
```bash
nahkoda liat kru  # otomatis di staging
```

### 3. Performance Tuning

#### Cache Lebih Lama (untuk cluster stabil):
```json
{
  "cache_ttl": 300000000000
}
```
Cache autocomplete selama 5 menit (lebih cepat, tapi kurang fresh).

#### Timeout Lebih Lama (untuk cluster lambat):
```json
{
  "timeout": 60000000000
}
```
Tunggu 1 menit sebelum timeout (cluster lambat tidak akan error terlalu cepat).

### 4. Disable Autocomplete

Untuk debugging atau environment dengan kubectl versi lama:

```json
{
  "enable_suggestions": false
}
```

## 🛠️ Membuat Config File

### Manual

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

### Via Editor

```bash
mkdir -p ~/.nahkoda
nano ~/.nahkoda/config.json
```

Paste config JSON, lalu save (Ctrl+O, Enter, Ctrl+X).

## ✅ Validasi

Nahkoda otomatis memvalidasi config saat load:

### ❌ Error: Invalid kubectl path
```json
{
  "kubectl_path": "/path/yang/tidak/ada"
}
```
**Result**: Nahkoda akan error saat start dengan pesan:
```
⚠️  Gagal memuat konfigurasi: kubectl path tidak valid
```

### ❌ Error: Negative values
```json
{
  "timeout": -5000000000
}
```
**Result**: Error `timeout tidak boleh negatif`

### ✅ Valid Config
Jika config valid, Nahkoda akan load tanpa pesan error.

## 🔍 Debugging

### Cek apakah config ter-load:
```bash
# Test dengan dry-run
nahkoda --dry-run liat kru
```

Jika menggunakan custom `kubectl_path`, akan terlihat di output.

### Fallback ke Default

Jika config error, Nahkoda akan:
1. Print warning
2. Fallback ke default config
3. Tetap berjalan normal

```
⚠️  Gagal memuat konfigurasi: <error message>
# Lanjut dengan default values
```

## 📊 Contoh Lengkap

### Production Environment

```json
{
  "kubectl_path": "/usr/local/bin/kubectl",
  "default_namespace": "production",
  "cache_ttl": 60000000000,
  "timeout": 45000000000,
  "enable_suggestions": true
}
```

- Kubectl standard di `/usr/local/bin`
- Default ke namespace `production`
- Cache 1 menit (cluster stabil)
- Timeout 45 detik (cluster cepat)
- Suggestions enabled

### Development Environment

```json
{
  "kubectl_path": "/snap/bin/microk8s.kubectl",
  "default_namespace": "dev",
  "cache_ttl": 10000000000,
  "timeout": 60000000000,
  "enable_suggestions": true
}
```

- MicroK8s kubectl
- Default ke namespace `dev`
- Cache 10 detik (cluster sering berubah)
- Timeout 1 menit (local cluster kadang lambat)
- Suggestions enabled

### Minimal Config

```json
{
  "default_namespace": "staging"
}
```

Hanya set default namespace, sisanya pakai default values.

## 🚀 Tips

1. **Jangan commit config ke git** - isi config bisa berbeda per developer
2. **Backup config** sebelum edit - copy dulu jika mau experiment
3. **Test dengan --dry-run** - pastikan config valid sebelum pakai
4. **Monitor error.log** - cek `~/.nahkoda/error.log` jika ada masalah

## 📚 Related

- [Instalasi](../../instalasi) - Install Nahkoda
- [Error Handling](../../perintah/error-handling) - Debugging errors
- [Bahasa](../bahasa) - Parser & AST internals
