#!/bin/bash
#
# Cancel Command Demonstration
#
# This script demonstrates the cancel command functionality.
# It shows how to:
# 1. Cancel a current generation
# 2. Cancel all queued generations
# 3. Use cancel in scripts
#
# Prerequisites:
# - asset-generator installed and in PATH
# - SwarmUI server running and configured
# - Valid API configuration

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "=========================================="
echo "Cancel Command Demonstration"
echo "=========================================="
echo

# Check if asset-generator is available
if ! command -v asset-generator &> /dev/null; then
    echo -e "${RED}Error: asset-generator not found in PATH${NC}"
    echo "Please install it first: go install ./..."
    exit 1
fi

# Check if server is configured
if ! asset-generator status &> /dev/null; then
    echo -e "${RED}Error: Cannot connect to SwarmUI server${NC}"
    echo "Please configure the API URL:"
    echo "  asset-generator config set api-url http://your-server:7801"
    exit 1
fi

echo -e "${GREEN}✓ Prerequisites met${NC}"
echo

# Demo 1: Cancel a single generation
echo "=========================================="
echo "Demo 1: Cancel Current Generation"
echo "=========================================="
echo
echo "This demo shows how to cancel a running generation."
echo
read -p "Press Enter to start a generation and then cancel it..."
echo

# Start a generation in the background
echo "Starting generation..."
asset-generator generate image \
    --prompt "detailed fantasy landscape with mountains" \
    --steps 30 \
    --width 512 \
    --height 512 \
    > /tmp/cancel-demo-gen.log 2>&1 &
GEN_PID=$!

echo -e "${YELLOW}Generation started (PID: $GEN_PID)${NC}"
echo "Waiting 2 seconds..."
sleep 2

echo
echo "Cancelling the generation..."
if asset-generator cancel; then
    echo -e "${GREEN}✓ Generation cancelled successfully${NC}"
else
    echo -e "${RED}✗ Cancel failed (generation may have completed)${NC}"
fi

# Clean up background process
kill $GEN_PID 2>/dev/null || true
wait $GEN_PID 2>/dev/null || true

echo
read -p "Press Enter to continue to Demo 2..."
echo

# Demo 2: Cancel all queued generations
echo "=========================================="
echo "Demo 2: Cancel All Queued Generations"
echo "=========================================="
echo
echo "This demo shows how to cancel multiple queued generations."
echo
read -p "Press Enter to queue multiple generations and cancel them all..."
echo

# Queue multiple generations
echo "Queueing 3 generations..."
for i in 1 2 3; do
    asset-generator generate image \
        --prompt "test generation $i" \
        --steps 5 \
        > /tmp/cancel-demo-gen$i.log 2>&1 &
    echo -e "${YELLOW}Queued generation $i${NC}"
    sleep 0.5
done

echo
echo "Waiting 1 second..."
sleep 1

echo
echo "Cancelling all generations..."
if asset-generator cancel --all; then
    echo -e "${GREEN}✓ All generations cancelled successfully${NC}"
else
    echo -e "${RED}✗ Cancel failed${NC}"
fi

# Clean up any remaining background processes
pkill -P $$ || true

echo
read -p "Press Enter to continue to Demo 3..."
echo

# Demo 3: Scripted usage
echo "=========================================="
echo "Demo 3: Cancel in Scripts"
echo "=========================================="
echo
echo "This demo shows how to use cancel in automated scripts."
echo
read -p "Press Enter to see script examples..."
echo

cat << 'EOF'
# Example 1: Emergency stop script
#!/bin/bash
asset-generator cancel --all -q || echo "Cancel failed"

# Example 2: Conditional cancel
#!/bin/bash
if asset-generator status | grep -q "generating"; then
    echo "Found active generation, cancelling..."
    asset-generator cancel
fi

# Example 3: Cleanup in error handler
#!/bin/bash
cleanup() {
    echo "Cleaning up..."
    asset-generator cancel --all -q
}
trap cleanup EXIT ERR

# Your generation script here...
asset-generator generate image --prompt "..."

# Example 4: Timeout-based cancel
#!/bin/bash
asset-generator generate image --prompt "..." &
GEN_PID=$!

# Wait up to 5 minutes
TIMEOUT=300
if ! timeout $TIMEOUT wait $GEN_PID; then
    echo "Generation timed out, cancelling..."
    asset-generator cancel
    kill $GEN_PID
fi

EOF

echo
read -p "Press Enter to continue to Demo 4..."
echo

# Demo 4: Verbose mode
echo "=========================================="
echo "Demo 4: Verbose Cancel"
echo "=========================================="
echo
echo "This demo shows verbose output for debugging."
echo
read -p "Press Enter to see verbose cancel output..."
echo

# Start a quick generation
echo "Starting a quick generation..."
asset-generator generate image \
    --prompt "test" \
    --steps 1 \
    > /tmp/cancel-demo-verbose.log 2>&1 &
GEN_PID=$!

sleep 1

echo
echo "Cancelling with verbose output..."
asset-generator cancel -v || echo "(Generation may have completed)"

kill $GEN_PID 2>/dev/null || true
wait $GEN_PID 2>/dev/null || true

echo
echo "=========================================="
echo "Demo 5: Status Check Before Cancel"
echo "=========================================="
echo
echo "Best practice: Check status before cancelling."
echo
read -p "Press Enter to see status check example..."
echo

echo "Checking server status..."
asset-generator status

echo
echo "Now you could cancel if there are active generations:"
echo "  asset-generator cancel"
echo

echo
echo "=========================================="
echo "Demonstration Complete!"
echo "=========================================="
echo
echo "Summary of cancel command usage:"
echo
echo "  Cancel current generation:"
echo "    asset-generator cancel"
echo
echo "  Cancel all queued generations:"
echo "    asset-generator cancel --all"
echo
echo "  Cancel quietly (for scripts):"
echo "    asset-generator cancel -q"
echo
echo "  Cancel with verbose output:"
echo "    asset-generator cancel -v"
echo
echo "For more information:"
echo "  asset-generator cancel --help"
echo "  See docs/CANCEL_COMMAND.md"
echo

# Clean up temp files
rm -f /tmp/cancel-demo-*.log

echo -e "${GREEN}✓ Demo complete${NC}"
