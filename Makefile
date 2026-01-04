PROTO_DIR := src/proto
GEN_DIR   := src/app/pb
GO_PB_PKG := github.com/simivar/tibia-sprites-exporter/src/app/pb

.PHONY: proto precheck tools

precheck:
	@command -v protoc >/dev/null 2>&1 || (echo "Missing dependency: protoc (Protocol Buffers compiler)"; exit 1)
	@command -v protoc-gen-go >/dev/null 2>&1 || (echo "Missing dependency: protoc-gen-go (Go protobuf plugin)"; exit 1)
	@test -f "$(PROTO_DIR)/appearances.proto" || (echo "Missing file: $(PROTO_DIR)/appearances.proto"; exit 1)
	@test -f "$(PROTO_DIR)/shared.proto" || (echo "Missing file: $(PROTO_DIR)/shared.proto"; exit 1)
	@echo "Dependencies OK"
	@echo " - protoc:        $$(command -v protoc)"
	@echo " - protoc-gen-go: $$(command -v protoc-gen-go)"
	@echo " - proto file:    $(PROTO_DIR)/appearances.proto"
	@echo " - proto file:    $(PROTO_DIR)/shared.proto"

proto: precheck
	@mkdir -p "$(GEN_DIR)"
	protoc \
	  --proto_path="$(PROTO_DIR)" \
	  --go_out="$(GEN_DIR)" \
	  --go_opt=paths=source_relative \
	  --go_opt=Mshared.proto=$(GO_PB_PKG) \
	  --go_opt=Mappearances.proto=$(GO_PB_PKG) \
	  appearances.proto shared.proto
	@echo "Generated protobuf sources in $(GEN_DIR)"

tools:
	@echo "Installing protoc-gen-go via 'go install'..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@echo "Done. Ensure GOPATH/bin is in your PATH."