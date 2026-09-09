---
title: Operasi Harian
description: Daftar perintah populer Nahkoda untuk mengelola kru, mesin, dan armada.
sidebar:
  order: 1
---

Nahkoda dirancang agar Kapten bisa mengemudikan kapal (klaster) tanpa harus mengingat parameter `kubectl` yang rumit.

:::note[Context dan koneksi cluster]
Saat REPL dibuka, Nahkoda memperingatkan jika belum ada context Kubernetes aktif. Operasi cluster juga gagal cepat dengan petunjuk `liat kapal` dan `pindah kapal <nama>`. Jika context ada tetapi kubectl belum merespons, status menunggu muncul setelah 2 detik dan command finite dihentikan sesuai timeout konfigurasi.
:::

:::caution[Perintah mutatif]
`bikin`, `hapus`, `atur`, dan `tukar` langsung mengubah klaster tanpa prompt konfirmasi tambahan. Gunakan `nahkoda --dry-run ...` terlebih dahulu, terutama pada klaster production.
:::

## ⌨️ Mode Interaktif (REPL)

Jika Kapten menjalankan `nahkoda` tanpa argumen apapun, Kapten akan masuk ke **Mode Interaktif**. Di sini, Kapten bisa mengetik perintah beruntun tanpa perlu mengetik prefix `nahkoda` lagi. Ketik `keluar` untuk kembali ke daratan.

Beberapa perintah operasi dan metrik yang didukung oleh Nahkoda.

| Perintah Nahkoda | Ekuivalen Kubectl | Fungsi |
| :--- | :--- | :--- |
| `nahkoda liat armada` | `kubectl get deployment -A` | List semua armada (deployment) |
| `nahkoda liat pelabuhan` | `kubectl get service -A` | List semua pelabuhan (service) |
| `nahkoda liat mercusuar` | `kubectl get ingress -A` | List semua mercusuar (ingress) |
| `nahkoda liat penjaga` | `kubectl get daemonset -A` | List semua penjaga (daemonset) |
| `nahkoda liat peta` | `kubectl get configmap -A` | List semua peta (configmap) |
| `nahkoda liat sandi` | `kubectl get secret -A` | List semua sandi (secret) |
| `nahkoda liat perbekalan kru` | *(Custom columns view)* | Intip jatah CPU/RAM kru |
| `nahkoda baca jurnal [nama] terus` | `kubectl logs -f [nama]` | Streaming log (tail -f) |
| `nahkoda baca jurnal [nama] kabin [c]` | `kubectl logs [nama] -c [c]` | Baca log dari kontainer spesifik |
| `nahkoda atur armada [nama] ke [jumlah]` | `kubectl scale deployment [nama] --replicas=[jumlah]` | Atur jumlah replika armada |
| `nahkoda tukar kru armada [nama]` | `kubectl rollout restart deployment [nama]` | Restart armada (deployment) |
| `nahkoda cek kesehatan` | *(Multi-command aggregation)* | Audit kesehatan klaster secara menyeluruh |
| `nahkoda bikin geladak [nama]` | `kubectl create namespace [nama]` | Buat namespace baru |
| `nahkoda bikin kru [nama]` | `kubectl run [nama] --image=nginx...` | Buat pod baru (nginx) |
| `nahkoda pantau kru` | `kubectl top pod` | Lihat penggunaan resource pod |
| `nahkoda pantau mesin` | `kubectl top node` | Lihat penggunaan resource node |
| `nahkoda hapus kru [nama]` | `kubectl delete pod [nama]` | Hapus pod |
