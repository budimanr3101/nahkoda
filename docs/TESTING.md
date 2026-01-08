# Nahkoda Integration Testing

## Test Resources

File `test-resources.yaml` berisi Kubernetes resources untuk comprehensive testing:

### Namespaces
- `auth` - Testing namespace filtering
- `payment` - Testing namespace filtering

### Pods

#### Healthy Pods (Running)
- `healthy-pod-1` (default namespace)
- `healthy-pod-2` (auth namespace)
- `my-app-*` (existing deployment)

#### Broken Pods (Rusak)
- `crashloop-pod` - CrashLoopBackOff status
- `imagepull-pod` - ImagePullBackOff status
- `imagepullbackoff-pod` - ImagePullBackOff status (existing)

#### Pending Pods (Terdampar)
- `pending-pod` - Unschedulable (999Gi memory request)

### Deployments
- `test-deployment` (payment namespace) - 2 replicas

## Setup

```bash
# Deploy test resources
kubectl apply -f test-resources.yaml

# Create crashloop pod separately
kubectl run crashloop-pod --image=busybox --restart=Never -- sh -c "exit 1"

# Verify resources
kubectl get pods -A
```

## Running Tests

```bash
# Run comprehensive test suite
./test.sh

# View results
cat nahkoda_test_v0.7.0.txt
```

## Test Coverage (32 Tests)

### Section 1: Help & Basic Commands (1 test)
- `--help`

### Section 2: Liat Kru - List Pods (4 tests)
- `liat kru` (default: sehat)
- `liat kru sehat`
- `liat kru rusak`
- `liat kru terdampar`

### Section 3: Liat Kru with Namespace (4 tests)
- `liat kru di geladak default`
- `liat kru di geladak auth`
- `liat kru di geladak payment`
- `liat kru di geladak kube-system`

### Section 4: Liat Kru with Condition + Namespace (4 tests)
- `liat kru sehat di geladak default`
- `liat kru rusak di geladak default`
- `liat kru sehat di geladak auth`
- `liat kru sehat di geladak payment`

### Section 5: Cek Kru - Describe Pod (5 tests)
- `cek kru healthy-pod-1`
- `cek kru crashloop-pod`
- `cek kru imagepull-pod`
- `cek kru pending-pod`
- `cek kru pod-tidak-ada` (not found)

### Section 6: Liat Mesin - List Nodes (3 tests)
- `liat mesin`
- `liat mesin siap`
- `liat mesin mogok`

### Section 7: Cek Mesin - Describe Node (2 tests)
- `cek mesin docker-desktop`
- `cek mesin node-tidak-ada` (not found)

### Section 8: Error Handling (6 tests)
- `liat kru xyz` (unknown word)
- `liat kru bocor` (unknown condition)
- `terbangkan kapal` (unknown action)
- `cek mesin` (missing target)
- `liat` (missing object)
- `kru rusak` (missing action)

### Section 9: Edge Cases (3 tests)
- `liat kru di auth` (invalid syntax)
- `LIAT KRU` (case insensitive)
- `Liat Kru Sehat` (mixed case)

## Cleanup

```bash
# Remove test resources
kubectl delete -f test-resources.yaml
kubectl delete pod crashloop-pod --force --grace-period=0
```

## Expected Results

All tests should complete successfully:
- ✅ Healthy pods shown in "liat kru" and "liat kru sehat"
- ✅ Broken pods shown in "liat kru rusak"
- ✅ Pending pods shown in "liat kru terdampar"
- ✅ Namespace filtering works correctly
- ✅ Describe commands show detailed info
- ✅ Error messages in Bahasa Indonesia
- ✅ Graceful handling of not found resources
- ✅ Case insensitive parsing
