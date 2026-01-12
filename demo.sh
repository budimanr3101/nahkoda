#!/bin/bash

# Demo script for Nahkoda terminal recording
# This script demonstrates the key features of Nahkoda

echo "🎬 Nahkoda Demo - Starting..."
echo ""
sleep 2

echo "⚓ Building Nahkoda..."
go build -o nahkoda main.go
sleep 1
echo ""

echo "📋 Running Unit Tests..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
sleep 1
go test ./... -v
sleep 2
echo ""

echo "✅ All tests passed!"
echo ""
sleep 2

echo "🚀 Demonstrating Nahkoda CLI..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
sleep 1

echo ""
echo "$ ./nahkoda -h"
sleep 1
./nahkoda -h
sleep 3

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎬 Demo Complete! Ready for documentation."
echo ""
