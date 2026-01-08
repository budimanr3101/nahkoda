#!/bin/bash

set -e

OUT="nahkoda_test_v0.7.0.txt"
BIN="./nahkoda"

echo "⚓ Nahkoda Test Run v0.7.0 - Comprehensive Integration Test" > $OUT
echo "Timestamp: $(date)" >> $OUT
echo "==============================" >> $OUT
echo "" >> $OUT

# Function to run test
run_test() {
    echo "▶ COMMAND: $1" >> $OUT
    echo "------------------------------" >> $OUT
    # Capture output and exit code
    if ! $BIN $1 >> $OUT 2>&1; then
        echo "❌ command failed" >> $OUT
        echo "❌ TEST FAILED: $1"
        exit 1
    fi
    echo "" >> $OUT
}

# Function to run test expected to fail
run_fail() {
    echo "▶ COMMAND (EXPECT FAIL): $1" >> $OUT
    echo "------------------------------" >> $OUT
    # Capture output and exit code
    if $BIN $1 >> $OUT 2>&1; then
        echo "❌ command succeeded (but expected to fail)" >> $OUT
        echo "❌ TEST FAILED (EXPECTED ERROR): $1"
        exit 1
    fi
    echo "✅ Success (failed as expected)" >> $OUT
    echo "" >> $OUT
}

# Build binary
echo "🔨 Building binary..." | tee -a $OUT
go build -o $BIN
echo "✅ Build success" | tee -a $OUT
echo "" >> $OUT

# ==========================================
# SECTION 1: HELP & BASIC COMMANDS
# ==========================================
echo "📖 SECTION 1: Help & Basic Commands" | tee -a $OUT
echo "==============================" >> $OUT
run_test "--help"

# ==========================================
# SECTION 2: LIAT KRU (LIST PODS)
# ==========================================
echo "📦 SECTION 2: Liat Kru (List Pods)" | tee -a $OUT
echo "==============================" >> $OUT
run_test "liat kru"
run_test "liat kru sehat"
run_test "liat kru rusak"
run_test "liat kru terdampar"

# ==========================================
# SECTION 3: LIAT KRU WITH NAMESPACE
# ==========================================
echo "🏢 SECTION 3: Liat Kru with Namespace" | tee -a $OUT
echo "==============================" >> $OUT
run_test "liat kru di geladak default"
run_test "liat kru di geladak auth"
run_test "liat kru di geladak payment"
run_test "liat kru di geladak kube-system"

# ==========================================
# SECTION 4: LIAT KRU WITH CONDITION + NAMESPACE
# ==========================================
echo "🔍 SECTION 4: Liat Kru with Condition + Namespace" | tee -a $OUT
echo "==============================" >> $OUT
run_test "liat kru sehat di geladak default"
run_test "liat kru rusak di geladak default"
run_test "liat kru sehat di geladak auth"
run_test "liat kru sehat di geladak payment"

# ==========================================
# SECTION 5: CEK KRU (DESCRIBE POD)
# ==========================================
echo "🔎 SECTION 5: Cek Kru (Describe Pod)" | tee -a $OUT
echo "==============================" >> $OUT
run_test "cek kru healthy-pod-1"
run_test "cek kru crashloop-pod"
run_test "cek kru imagepull-pod"
run_test "cek kru pending-pod"
run_test "cek kru pod-tidak-ada"

# ==========================================
# SECTION 6: LIAT MESIN (LIST NODES)
# ==========================================
echo "🖥️  SECTION 6: Liat Mesin (List Nodes)" | tee -a $OUT
echo "==============================" >> $OUT
run_test "liat mesin"
run_test "liat mesin siap"
run_test "liat mesin mogok"

# ==========================================
# SECTION 7: CEK MESIN (DESCRIBE NODE)
# ==========================================
echo "🔧 SECTION 7: Cek Mesin (Describe Node)" | tee -a $OUT
echo "==============================" >> $OUT
run_test "cek mesin docker-desktop"
run_test "cek mesin node-tidak-ada"

# ==========================================
# SECTION 8: ERROR HANDLING
# ==========================================
echo "⚠️  SECTION 8: Error Handling" | tee -a $OUT
echo "==============================" >> $OUT
run_fail "liat kru xyz"
run_fail "liat kru bocor"
run_fail "terbangkan kapal"
run_fail "cek mesin"
run_fail "liat"
run_fail "kru rusak"

# ==========================================
# SECTION 9: EDGE CASES
# ==========================================
echo "🎯 SECTION 9: Edge Cases" | tee -a $OUT
echo "==============================" >> $OUT
run_fail "liat kru di auth"
run_test "LIAT KRU"
run_test "Liat Kru Sehat"

# ==========================================
# SECTION 10: KAPAL (CONTEXT)
# ==========================================
echo "⚓ SECTION 10: Kapal (Context)" | tee -a $OUT
echo "==============================" >> $OUT
run_test "liat kapal"

# ==========================================
# SUMMARY
# ==========================================
echo "" | tee -a $OUT
echo "==============================" | tee -a $OUT
echo "✅ Test selesai. Output tersimpan di $OUT" | tee -a $OUT
echo "" | tee -a $OUT

# Count tests
TOTAL=$(grep -c "▶ COMMAND:" $OUT)
echo "📊 Total tests run: $TOTAL" | tee -a $OUT
