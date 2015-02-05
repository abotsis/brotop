VERSION=$(shell cat brotop.go | grep -oP "Version\s+?\=\s?\"\K.*?(?=\"$|$\)")
CWD=$(shell pwd)

NAME="brotop"
DESCRIPTION="Top for bro log files."

NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
UNAME := $(shell uname -s)

ifeq ($(UNAME),Darwin)
	ECHO=echo
else
	ECHO=/bin/echo -e
endif

all: deps
	@mkdir -p bin/
	@$(ECHO) "$(OK_COLOR)==> Building $(NAME) - $(VERSION) $(NO_COLOR)"
	@godep go build -o bin/$(NAME)
	@chmod +x bin/$(NAME)
	@$(ECHO) "$(OK_COLOR)==> Done$(NO_COLOR)"


deps:
	@$(ECHO) "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	@godep get

test: deps
	@$(ECHO) "$(OK_COLOR)==> Testing $(NAME)...$(NO_COLOR)"
	go test ./...

clean:
	rm -rf bin/
	rm -rf pkg/

install: clean all

uninstall: clean

tar: 

.PHONY: all clean deps
