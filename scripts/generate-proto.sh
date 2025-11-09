#!/bin/bash

# Script to generate Go and Python code from domain-local protobuf definitions
# Generates:
#   - Server code in each service's api/v1/ directory
#   - Client code (all Go) in gateway/api/ directory
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
GATEWAY_API_DIR="$ROOT_DIR/gateway/api"

echo -e "${GREEN}=== Scout Proto Generation (Domain-Local + Gateway Clients) ===${NC}"
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
    
    # 1. Generate server code in service directory
    if [ -f "$SERVICE_DIR/go.mod" ]; then
        # Go service - generate Go server code
        echo "  Generating Go server code for $service..."
        
        protoc \
            --proto_path="$PROTO_DIR" \
            --go_out="$PROTO_DIR" \
            --go_opt=paths=source_relative \
            --go-grpc_out="$PROTO_DIR" \
            --go-grpc_opt=paths=source_relative \
            "$PROTO_FILE"
        
        echo -e "  ${GREEN}✓ Go server code generated${NC}"
        
    elif [ -f "$SERVICE_DIR/pyproject.toml" ]; then
        # Python service - generate Python server code
        echo "  Generating Python server code for $service..."
        
        cd "$SERVICE_DIR"
        
        poetry run python -m grpc_tools.protoc \
            --proto_path="api/v1" \
            --python_out="api/v1" \
            --grpc_python_out="api/v1" \
            "api/v1/${service}.proto"
        
        # Fix Python imports (make them relative)
        GRPC_FILE="$PROTO_DIR/${service}_pb2_grpc.py"
        if [ -f "$GRPC_FILE" ]; then
            if [[ "$OSTYPE" == "darwin"* ]]; then
                sed -i '' "s/import ${service}_pb2/from . import ${service}_pb2/" "$GRPC_FILE"
            else
                sed -i "s/import ${service}_pb2/from . import ${service}_pb2/" "$GRPC_FILE"
            fi
        fi
        
        # Create __init__.py if it doesn't exist
        if [ ! -f "$PROTO_DIR/__init__.py" ]; then
            cat > "$PROTO_DIR/__init__.py" << 'EOF'
"""Generated protocol buffer code."""
EOF
        fi
        
        echo -e "  ${GREEN}✓ Python server code generated${NC}"
        
        cd "$ROOT_DIR"
    fi
    
    # 2. Generate Go client code in gateway/api/ for ALL services
    echo "  Generating Go client code in gateway/api/${service}/v1/..."
    
    GATEWAY_SERVICE_DIR="$GATEWAY_API_DIR/$service/v1"
    mkdir -p "$GATEWAY_SERVICE_DIR"
    
    protoc \
        --proto_path="$PROTO_DIR" \
        --go_out="$GATEWAY_SERVICE_DIR" \
        --go_opt=paths=source_relative \
        --go-grpc_out="$GATEWAY_SERVICE_DIR" \
        --go-grpc_opt=paths=source_relative \
        "$PROTO_FILE"
    
    echo -e "  ${GREEN}✓ Go client code generated in gateway${NC}"
    echo ""
done

# Summary
echo -e "${GREEN}=== Generation Complete ===${NC}"
echo ""
echo "Generated files:"
echo ""
echo "  Service Implementations (Server):"
echo "    services/agents/api/v1/{agents.pb.go,agents_grpc.pb.go}"
echo "    services/incidents/api/v1/{incidents.pb.go,incidents_grpc.pb.go}"
echo "    services/tools/api/v1/{tools_pb2.py,tools_pb2_grpc.py}"
echo "    services/analytics/api/v1/{analytics_pb2.py,analytics_pb2_grpc.py}"
echo ""
echo "  Gateway Clients (All Go):"
echo "    gateway/api/agents/v1/{agents.pb.go,agents_grpc.pb.go}"
echo "    gateway/api/incidents/v1/{incidents.pb.go,incidents_grpc.pb.go}"
echo "    gateway/api/tools/v1/{tools.pb.go,tools_grpc.pb.go}"
echo "    gateway/api/analytics/v1/{analytics.pb.go,analytics_grpc.pb.go}"
echo ""
echo -e "${GREEN}✓ All done!${NC}"