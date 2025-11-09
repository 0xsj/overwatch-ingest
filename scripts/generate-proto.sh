#!/bin/bash

# Script to generate Go and Python code from domain-local protobuf definitions
# Usage: ./scripts/generate-proto.sh

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Directories
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo -e "${GREEN}=== Scout Proto Generation (Domain-Local) ===${NC}"
echo ""

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo -e "${RED}Error: protoc not found${NC}"
    echo "Install with: brew install protobuf"
    exit 1
fi

echo -e "${GREEN}✓ protoc found:${NC} $(protoc --version)"

# Check if Go plugins are installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo -e "${RED}Error: protoc-gen-go not found${NC}"
    echo "Install with: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    exit 1
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo -e "${RED}Error: protoc-gen-go-grpc not found${NC}"
    echo "Install with: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

echo -e "${GREEN}✓ Go plugins found${NC}"

echo ""
echo -e "${GREEN}=== Generating Proto Code ===${NC}"
echo ""

# Define domain services
SERVICES=("agents" "incidents" "tools" "analytics")

# Generate for each service
for service in "${SERVICES[@]}"; do
    SERVICE_DIR="$ROOT_DIR/services/$service"
    PROTO_DIR="$SERVICE_DIR/api/v1"
    PROTO_FILE="$PROTO_DIR/${service}.proto"
    
    if [ ! -f "$PROTO_FILE" ]; then
        echo -e "${YELLOW}Skipping $service: proto file not found${NC}"
        continue
    fi
    
    echo -e "${YELLOW}Processing $service service...${NC}"
    
    # Determine if this is a Go or Python service
    if [ -f "$SERVICE_DIR/go.mod" ]; then
        # Go service
        echo "  Generating Go code for $service..."
        
        protoc \
            --proto_path="$PROTO_DIR" \
            --go_out="$PROTO_DIR" \
            --go_opt=paths=source_relative \
            --go-grpc_out="$PROTO_DIR" \
            --go-grpc_opt=paths=source_relative \
            "$PROTO_FILE"
        
        echo -e "  ${GREEN}✓ Go code generated${NC}"
        
    elif [ -f "$SERVICE_DIR/pyproject.toml" ]; then
        # Python service
        echo "  Generating Python code for $service..."
        
        cd "$SERVICE_DIR"
        
        # Use poetry run to ensure we're in the right environment
        poetry run python -m grpc_tools.protoc \
            --proto_path="api/v1" \
            --python_out="api/v1" \
            --grpc_python_out="api/v1" \
            "api/v1/${service}.proto"
        
        # Fix Python imports (make them relative)
        GRPC_FILE="$PROTO_DIR/${service}_pb2_grpc.py"
        if [ -f "$GRPC_FILE" ]; then
            if [[ "$OSTYPE" == "darwin"* ]]; then
                # macOS
                sed -i '' "s/import ${service}_pb2/from . import ${service}_pb2/" "$GRPC_FILE"
            else
                # Linux
                sed -i "s/import ${service}_pb2/from . import ${service}_pb2/" "$GRPC_FILE"
            fi
        fi
        
        # Create __init__.py if it doesn't exist
        if [ ! -f "$PROTO_DIR/__init__.py" ]; then
            cat > "$PROTO_DIR/__init__.py" << 'EOF'
"""Generated protocol buffer code."""
EOF
        fi
        
        echo -e "  ${GREEN}✓ Python code generated${NC}"
        
        cd "$ROOT_DIR"
    else
        echo -e "  ${YELLOW}Unknown service type (no go.mod or pyproject.toml)${NC}"
    fi
    
    echo ""
done

# Summary
echo -e "${GREEN}=== Generation Complete ===${NC}"
echo ""
echo "Generated files:"
echo "  Go services:"
echo "    services/agents/api/v1/{agents.pb.go,agents_grpc.pb.go}"
echo "    services/incidents/api/v1/{incidents.pb.go,incidents_grpc.pb.go}"
echo ""
echo "  Python services:"
echo "    services/tools/api/v1/{tools_pb2.py,tools_pb2_grpc.py}"
echo "    services/analytics/api/v1/{analytics_pb2.py,analytics_pb2_grpc.py}"
echo ""
echo -e "${GREEN}✓ All done!${NC}"