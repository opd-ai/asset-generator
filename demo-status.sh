#!/bin/bash
# Demo script for the status command
# Shows various ways to use the status command

set -e

echo "════════════════════════════════════════════════════════════════"
echo "Asset Generator CLI - Status Command Demo"
echo "════════════════════════════════════════════════════════════════"
echo ""

# Demo 1: Basic status check
echo "1. Basic Status Check"
echo "────────────────────────────────────────────────────────────────"
echo "$ asset-generator status"
echo ""
asset-generator status 2>&1 || echo "(Server may be offline - this is expected if SwarmUI is not running)"
echo ""
echo ""

# Demo 2: JSON format
echo "2. JSON Output Format"
echo "────────────────────────────────────────────────────────────────"
echo "$ asset-generator status --format json"
echo ""
asset-generator status --format json 2>&1 || echo '{"status": "offline", "note": "Server not running"}'
echo ""
echo ""

# Demo 3: YAML format
echo "3. YAML Output Format"
echo "────────────────────────────────────────────────────────────────"
echo "$ asset-generator status --format yaml"
echo ""
asset-generator status --format yaml 2>&1 || echo "status: offline"
echo ""
echo ""

# Demo 4: Extract specific information with jq
echo "4. Extract Specific Information (JSON + jq)"
echo "────────────────────────────────────────────────────────────────"
echo "$ asset-generator status --format json | jq -r '.status'"
echo ""
asset-generator status --format json 2>/dev/null | jq -r '.status' 2>/dev/null || echo "offline"
echo ""
echo "$ asset-generator status --format json | jq '.models_count'"
echo ""
asset-generator status --format json 2>/dev/null | jq '.models_count' 2>/dev/null || echo "0"
echo ""
echo ""

# Demo 5: Health check example
echo "5. Health Check Script Example"
echo "────────────────────────────────────────────────────────────────"
cat << 'EOF'
#!/bin/bash
# Example health check script
if asset-generator status > /dev/null 2>&1; then
    echo "✓ Server is online"
    exit 0
else
    echo "✗ Server is offline"
    exit 1
fi
EOF
echo ""
echo "Running health check..."
if asset-generator status > /dev/null 2>&1; then
    echo "✓ Server is online"
else
    echo "✗ Server is offline (this is expected if SwarmUI is not running)"
fi
echo ""
echo ""

# Demo 6: Monitoring example
echo "6. Monitoring Script Example"
echo "────────────────────────────────────────────────────────────────"
cat << 'EOF'
#!/bin/bash
# Example monitoring script
while true; do
    timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    status=$(asset-generator status --format json 2>/dev/null | jq -r '.status' 2>/dev/null || echo "offline")
    echo "[$timestamp] Server status: $status"
    sleep 60
done
EOF
echo ""
echo "Would run continuously, checking status every 60 seconds..."
echo ""
echo ""

# Demo 7: Pre-flight check
echo "7. Pre-flight Check Before Generation"
echo "────────────────────────────────────────────────────────────────"
cat << 'EOF'
#!/bin/bash
# Check server before running expensive operations
echo "Checking server status..."
if ! asset-generator status > /dev/null 2>&1; then
    echo "Error: Server is offline. Cannot proceed with generation."
    exit 1
fi

echo "Server is online. Starting generation..."
asset-generator generate image --prompt "beautiful landscape" --batch 10
EOF
echo ""
echo ""

# Demo 8: Help output
echo "8. Status Command Help"
echo "────────────────────────────────────────────────────────────────"
echo "$ asset-generator status --help"
echo ""
asset-generator status --help
echo ""

echo "════════════════════════════════════════════════════════════════"
echo "Demo Complete!"
echo ""
echo "To run the status command yourself:"
echo "  asset-generator status"
echo ""
echo "For more information:"
echo "  - Full docs: docs/STATUS_COMMAND.md"
echo "  - Quick ref: docs/STATUS_QUICKREF.md"
echo "════════════════════════════════════════════════════════════════"
