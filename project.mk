APP_NAME ?= mtproto-client
BIN_DIR ?= ./bin
CMD_DIR ?= ./cmd

.PHONY: help
## Вывести список доступных команд
help:
	@echo ""
	@echo "Использование:"
	@echo "  make <цель>"
	@echo ""
	@echo "Доступные цели:"
	@awk ' \
		BEGIN { comment=""; } \
		/^[ \t]*##/ { \
			comment = substr($$0, index($$0, $$2)); next; \
		} \
		/^[a-zA-Z0-9_-]+:/ { \
			target = $$1; \
			sub(":", "", target); \
			if (comment != "") { \
				printf "  %-17s %s\n", target, comment; \
				comment=""; \
			} \
		}' $(MAKEFILE_LIST)
	@echo ""

.PHONY: build
## Скомпилировать Go приложение
build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)