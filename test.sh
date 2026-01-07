#!/usr/bin/env bash

set -e

OUT="nahkoda_test_v0.1.0.txt"
BIN="./nahkoda"

echo "⚓ Nahkoda Test Run" > $OUT
echo "Timestamp: $(date)" >> $OUT
echo "==============================" >> $OUT
echo "" >> $OUT

echo "🔨 Building binary..." | tee -a $OUT
go build -o nahkoda
echo "✅ Build success" | tee -a $OUT
echo "" >> $OUT

run_test () {
  CMD="$1"
  echo "▶ COMMAND: $CMD" | tee -a $OUT
  echo "------------------------------" >> $OUT
  $BIN $CMD >> $OUT 2>&1 || echo "❌ command failed" >> $OUT
  echo "" >> $OUT
}

run_test "--help"
run_test "liat kru"
run_test "liat kru rusak"
run_test "liat kru di geladak auth"
run_test "liat kru bocor di geladak payment"
run_test "liat kru sehat"
run_test "liat kru xyz"
run_test "terbangkan kapal"

echo "✅ Test selesai. Output tersimpan di $OUT"
