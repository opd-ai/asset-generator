#!/bin/bash
# Demonstration of file-based state sharing for generation tracking
# This shows how the status command can now track generations across processes

set -e

echo "════════════════════════════════════════════════════════════════"
echo "File-Based State Sharing Demo"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo "This demo shows how generation sessions are now tracked across"
echo "different CLI process invocations using a state file."
echo ""

# Create a test directory
TEST_DIR=$(mktemp -d)
echo "Test directory: $TEST_DIR"
echo ""

cd "$TEST_DIR"

echo "1. Checking initial state (no state file exists yet)"
echo "────────────────────────────────────────────────────────────────"
if [ -f .asset-generator-state.json ]; then
    echo "State file exists:"
    cat .asset-generator-state.json | jq .
else
    echo "✓ No state file (as expected)"
fi
echo ""

echo "2. Simulating generation start (creating state file)"
echo "────────────────────────────────────────────────────────────────"
# Create a mock state file to simulate an active generation
cat > .asset-generator-state.json << 'EOF'
{
  "sessions": {
    "demo-session-123": {
      "id": "demo-session-123",
      "status": "generating",
      "progress": 0.45,
      "start_time": "2025-10-10T14:30:00Z",
      "updated_at": "2025-10-10T14:32:15Z"
    }
  },
  "updated_at": "2025-10-10T14:32:15Z"
}
EOF

echo "✓ State file created with active session"
echo ""

echo "3. Viewing state file contents"
echo "────────────────────────────────────────────────────────────────"
cat .asset-generator-state.json | jq .
echo ""

echo "4. What the status command would see:"
echo "────────────────────────────────────────────────────────────────"
echo "Session ID:    demo-session-123"
echo "Status:        generating"
echo "Progress:      45.0%"
echo "Start Time:    2025-10-10T14:30:00Z"
echo ""

echo "5. Querying state with jq"
echo "────────────────────────────────────────────────────────────────"
echo "Number of active sessions:"
jq '.sessions | length' .asset-generator-state.json
echo ""
echo "Session IDs:"
jq -r '.sessions | keys[]' .asset-generator-state.json
echo ""
echo "Progress percentage:"
jq -r '.sessions[] | (.progress * 100)' .asset-generator-state.json
echo ""

echo "6. Demonstrating directory isolation"
echo "────────────────────────────────────────────────────────────────"
ANOTHER_DIR=$(mktemp -d)
cd "$ANOTHER_DIR"
echo "Switched to different directory: $ANOTHER_DIR"
if [ -f .asset-generator-state.json ]; then
    echo "Has state file"
else
    echo "✓ No state file in different directory (project isolation)"
fi
echo ""

cd "$TEST_DIR"
echo "Switched back to: $TEST_DIR"
if [ -f .asset-generator-state.json ]; then
    echo "✓ State file still exists here"
fi
echo ""

echo "7. Cleanup simulation (generation completes)"
echo "────────────────────────────────────────────────────────────────"
# Update state to show completed
cat > .asset-generator-state.json << 'EOF'
{
  "sessions": {},
  "updated_at": "2025-10-10T14:35:00Z"
}
EOF
echo "✓ State updated - no active sessions"
echo ""

echo "Updated state file:"
cat .asset-generator-state.json | jq .
echo ""

echo "8. Cleanup"
echo "────────────────────────────────────────────────────────────────"
cd /
rm -rf "$TEST_DIR" "$ANOTHER_DIR"
echo "✓ Test directories removed"
echo ""

echo "════════════════════════════════════════════════════════════════"
echo "Demo Complete!"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo "Key Takeaways:"
echo "  • State file: .asset-generator-state.json in working directory"
echo "  • Enables cross-process generation tracking"
echo "  • Directory-specific (project isolation)"
echo "  • Automatically managed (created/updated/cleaned)"
echo "  • JSON format for easy querying with jq"
echo ""
echo "Real Usage:"
echo "  Terminal 1: asset-generator generate image --prompt '...'"
echo "  Terminal 2: cd same-directory && asset-generator status"
echo "  Result: See full generation details across processes!"
echo ""
