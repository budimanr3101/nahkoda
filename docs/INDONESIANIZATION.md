# Indonesianization of Error Messages

## ✅ Completed

All error messages in Nahkoda have been updated to use **Bahasa Indonesia** for consistency with the project's philosophy.

## Changes Made

### Error Package (`internal/errors/errors.go`)

**Before**:
```go
NewKubectlFailed(err) // "kubectl command failed: ..."
NewResourceNotFound(r) // "resource \"x\" not found"
```

**After**:
```go
NewKubectlFailed(err) // "perintah kubectl gagal: ..."
NewResourceNotFound(r) // "resource \"x\" tidak ditemukan"
```

### Error Messages Mapping

| Error Type | English (Before) | Indonesian (After) |
|------------|------------------|-------------------|
| Kubectl Failed | kubectl command failed | perintah kubectl gagal |
| Resource Not Found | resource "x" not found | resource "x" tidak ditemukan |
| Unknown Word | (already ID) | kata tidak dikenali |
| Unknown Action | (already ID) | aksi tidak dikenali |
| Unknown Object | (already ID) | objek tidak dikenali |
| Unknown Condition | (already ID) | kondisi tidak dikenali |
| Missing Target | (already ID) | cek X butuh nama X |

## Kosa Kata v1.0.0

| Kubernetes | Nahkoda | Filosofi |
| :--- | :--- | :--- |
| Deployment | Armada | Grup besar unit kerja. |
| DaemonSet | Penjaga | Unit yang wajib ada di tiap titik. |
| Service | Pelabuhan | Titik bongkar muat logistik/data. |
| Ingress | Mercusuar | Pemandu arah dari luar. |
| ConfigMap | Peta | Instruksi konfigurasi statis. |
| Secret | Sandi | Data rahasia dan kunci keamanan. |

## Consistency

All user-facing messages now use Bahasa Indonesia:
- ✅ Command vocabulary (liat, cek, kru, mesin, etc.)
- ✅ Error messages
- ✅ Help text
- ✅ Status messages

## Testing

- ✅ All unit tests updated and passing
- ✅ Error messages verified with real commands
- ✅ Build successful
- ✅ No breaking changes

## Example Output

```bash
$ nahkoda liat kru xyz
kata tidak dikenali: "xyz"

$ nahkoda cek mesin
cek mesin butuh nama mesin

$ nahkoda liat kru bocor
kondisi tidak dikenali: bocor
```

## Philosophy Alignment

This change aligns perfectly with Nahkoda's core philosophy:
> **Bahasa manusia untuk Kubernetes**

Now, not just the commands but also the error messages speak in the user's language, making the entire experience more natural and accessible for Indonesian speakers.
