#!/bin/bash
set -e

# Always run from the project root, regardless of where the script is called from
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"

# Step 1: Build contracts to generate ABIs
COLLECTION_VAULT_DIR="/Users/andrey/projects/lend.fam MVP/collection-vault"
cd "$COLLECTION_VAULT_DIR"
forge build src
cd "$PROJECT_ROOT"

# Step 2: Generate Go bindings for each interface
INTERFACE_DIR="$COLLECTION_VAULT_DIR/src/interfaces"
ABI_DIR="$COLLECTION_VAULT_DIR/out"
GO_OUT_DIR="pkg/contracts"

mkdir -p "$GO_OUT_DIR"

for sol in "$INTERFACE_DIR"/*.sol; do
    name=$(basename "$sol" .sol)
    artifact_file="$ABI_DIR/$name.sol/$name.json"
    abi_file="$ABI_DIR/$name.sol/$name.abi.json"
    if [ -f "$artifact_file" ]; then
        # Extract only the abi array using jq
        jq .abi "$artifact_file" > "$abi_file"
        abigen --v2 --abi "$abi_file" --pkg contracts --type "$name" --out "$GO_OUT_DIR/$name.go"
        echo "Generated Go binding for $name"
    else
        echo "ABI not found for $name, skipping."
    fi
done

# Generate Go binding for IERC20
INTERFACE_DIR_ERC20="$COLLECTION_VAULT_DIR/lib/forge-std/src/interfaces"
ERC20_SOL="$INTERFACE_DIR_ERC20/IERC20.sol"
ERC20_NAME="IERC20"
ERC20_ARTIFACT_FILE="$ABI_DIR/$ERC20_NAME.sol/$ERC20_NAME.json"
ERC20_ABI_FILE="$ABI_DIR/$ERC20_NAME.sol/$ERC20_NAME.abi.json"

if [ -f "$ERC20_ARTIFACT_FILE" ]; then
    jq .abi "$ERC20_ARTIFACT_FILE" > "$ERC20_ABI_FILE"
    abigen --v2 --abi "$ERC20_ABI_FILE" --pkg contracts --type "$ERC20_NAME" --out "$GO_OUT_DIR/$ERC20_NAME.go"
    echo "Generated Go binding for $ERC20_NAME"
else
    echo "ABI not found for $ERC20_NAME, skipping."
fi
