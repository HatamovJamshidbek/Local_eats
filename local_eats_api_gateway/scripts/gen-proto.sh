#!/bin/sh

# Use the current directory if none is provided
if [ -z "$1" ]; then
  CURRENT_DIR=$(pwd)
else
  CURRENT_DIR=$1
fi

GENPROTO_DIR="${CURRENT_DIR}/genproto"
PROTOS_DIR="${CURRENT_DIR}/protos"

# Check if the user.proto directory exists
if [ ! -d "$PROTOS_DIR" ]; then
  echo "Error: Protos directory $PROTOS_DIR does not exist."
  exit 1
fi

# Remove old user.proto directory and create a new one
rm -rf "$GENPROTO_DIR"
mkdir -p "$GENPROTO_DIR"

# Find all directories under the user.proto directory and run protoc on them
for x in $(find "$PROTOS_DIR" -type d); do
  if ls "$x"/*.proto 1> /dev/null 2>&1; then
    protoc -I="$x" \
      --go_out="$GENPROTO_DIR" --go_opt=paths=source_relative \
      --go-grpc_out="$GENPROTO_DIR" --go-grpc_opt=paths=source_relative "$x"/*.proto
  else
    echo "Warning: No .proto files found in directory $x"
  fi
done
