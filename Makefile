GO_MODULE := github.com/CEBEP9HUH/OTUS_Go_Diploma

# Переменные для генерации grpc
PROTOC := protoc
PROTO_EXT := .proto
GEN_PATH := pkg
PROTOS_PATH := api/grpc/statistic
PROTO_FILES := $(wildcard $(PROTOS_PATH)/*$(echo PROTO_EXT))
PROTO_OPT := --go_opt
API_GEN_PROTO_OPTS := $(foreach file, $(basename $(PROTO_FILES)), $(PROTO_OPT)=M$(file)$(PROTO_EXT)=$(GO_MODULE)/$(GEN_PATH)/$(file))
GRPC_OPT := --go-grpc_opt
API_GEN_GRPC_OPTS := $(foreach file, $(basename $(PROTO_FILES)), $(GRPC_OPT)=M$(file)$(PROTO_EXT)=$(GO_MODULE)/$(GEN_PATH)/$(file))

# Переменные для сборки
BUILD_DST := artifacts
SERVER := server
CLIENT := client


genapi:
	@rm -rf $(GEN_PATH)
	@mkdir $(GEN_PATH)
	$(PROTOC) \
		--go_out=$(GEN_PATH) \
		--go_opt=module=$(GO_MODULE)/$(GEN_PATH) \
		--go-grpc_out=$(GEN_PATH) \
		--go-grpc_opt=module=$(GO_MODULE)/$(GEN_PATH) \
		$(API_GEN_PROTO_OPTS) $(API_GEN_GRPC_OPTS) \
		$(PROTO_FILES)

clean:
	@rm -rf $(GEN_PATH)
	@rm -rf $(BUILD_DST)

test:
	go test ./... -race -count 100

server:
	@if [ ! -d $(BUILD_DST) ]; then mkdir $(BUILD_DST); fi
	@if [ ! -d $(GEN_PATH) ]; then make genapi; fi
	go build -o $(BUILD_DST)/$(SERVER) ./cmd/$(SERVER)/*

client:
	@if [ ! -d $(BUILD_DST) ]; then mkdir $(BUILD_DST); fi
	@if [ ! -d $(GEN_PATH) ]; then make genapi; fi
	go build -o $(BUILD_DST)/$(CLIENT) ./cmd/$(CLIENT)/*

all: clean genapi test server client
	@cp configs/config.json $(BUILD_DST)