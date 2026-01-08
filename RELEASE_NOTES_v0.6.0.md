# Release Notes - Nahkoda v0.6.0

## 🎉 What's New

### Node Status Support (v0.6.0)
- **New Commands**: `siap` (Ready) and `mogok` (NotReady) for checking node health
- **Regex Filtering**: Precise matching using word boundaries (`\bReady\b`) to distinguish "Ready" from "NotReady"
- **Example**:
  ```bash
  nahkoda liat mesin siap    # Show only Ready nodes
  nahkoda liat mesin mogok   # Show NotReady nodes
  ```

## 🔧 Improvements

### Graceful Error Handling (v0.5.1)
- **No More Panic**: Missing resources (`cek kru hantu`) no longer exit with error code 1
- **Clean Output**: kubectl "NotFound" errors are treated as informative feedback, not system failures
- **Example**:
  ```bash
  nahkoda cek kru hantu
  # Output: Error from server (NotFound): pods "hantu" not found
  # (No "❌ command failed")
  ```

### Graceful Input Validation (v0.5.2)
- **Friendly Errors**: Invalid input (typos, unknown words) returns clean feedback without panic
- **Exit Code 0**: Validation errors no longer treated as system crashes
- **Removed Prefix**: No more "❌ " prefix for cleaner, more natural output
- **Examples**:
  ```bash
  nahkoda liat kru xyz
  # Output: kata tidak dikenali: "xyz"
  
  nahkoda cek mesin
  # Output: cek mesin butuh nama mesin
  ```

### Mature Filtering (v0.5.0)
- **Client-Side Text Filtering**: More accurate status detection
- **CrashLoopBackOff Detection**: Now correctly identified as "rusak" (broken)
- **Human-Aligned Results**: Filters based on what humans see in STATUS column, not just API phase
- **Example**:
  ```bash
  nahkoda liat kru rusak
  # Now shows: CrashLoopBackOff, ImagePullBackOff, etc.
  ```

### kubectl Integration (v0.4.0)
- **Real Execution**: Replaced mock executor with actual `kubectl` commands
- **Dynamic Arguments**: Constructs proper kubectl commands with namespaces and field selectors
- **Example**:
  ```bash
  nahkoda liat kru sehat
  # Runs: kubectl get pod -A | grep 'Running' (invert=false)
  ```

### Mesin (Node) Support (v0.3.1)
- **New Object**: `mesin` as alias for Kubernetes nodes
- **Commands**:
  ```bash
  nahkoda liat mesin          # List all nodes
  nahkoda cek mesin node-1    # Describe specific node
  ```

## 📚 Complete Vocabulary

| Nahkoda   | Kubernetes          |
|-----------|---------------------|
| kru       | pod                 |
| mesin     | node                |
| geladak   | namespace           |
| rusak     | status != Running   |
| sehat     | status = Running    |
| terdampar | Pending             |
| siap      | status = Ready      |
| mogok     | status = NotReady   |

## 🔄 Breaking Changes

None. All changes are backward compatible.

## 🐛 Bug Fixes

- Fixed default filter logic: `liat mesin` no longer applies `status.phase=Running` filter
- Fixed `cek` command to use correct namespace (default instead of `-A`)
- Fixed semantic gap where CrashLoopBackOff was incorrectly shown as "healthy"
- Fixed "Ready" substring matching issue with "NotReady" using regex

## 🏗️ Technical Details

### Architecture Improvements
- **Planner**: Added `GrepRegex` field for precise text matching
- **Executor**: Implemented regex support using `regexp.MatchString`
- **Semantic**: Abstracted status filters from `status.phase` to `status` for flexibility
- **Parser**: Extended condition tokens to include `siap`, `mogok`, `terdampar`

### Code Quality
- Removed duplicate `resolveCondition` function
- Centralized condition mapping in `conditions.go`
- Improved error handling across all layers
- Better separation of concerns (server-side vs client-side filtering)

## 📊 Test Coverage

All features verified via `test.sh`:
- ✅ Pod listing with various filters
- ✅ Node listing and status filtering
- ✅ Resource description (cek)
- ✅ Namespace scoping
- ✅ Error handling (missing resources, invalid input)
- ✅ Graceful validation

## 🚀 Upgrade Guide

No special steps required. Simply pull the latest version:

```bash
git pull origin main
git checkout v0.6.0
go build -o nahkoda
```

## 🙏 Acknowledgments

This release represents a significant maturation of Nahkoda, transforming it from a proof-of-concept to a production-ready CLI tool with:
- Real kubectl integration
- Robust error handling
- Human-friendly output
- Precise filtering capabilities

---

**Full Changelog**: v0.3.0...v0.6.0
