BUILD_DIR 	:= build
SCRIPTS_DIR := scripts
MSG ?= "Default commit message"

build:
	go build -o $(BUILD_DIR)/server main/*

run:
	cd $(BUILD_DIR) && ./server

runserver: build run

format:
	./$(SCRIPTS_DIR)/pre-commit

commit: format
	echo "\033[32mstagging changes...\033[0m"
	@git add .
	echo "\033[32mcommiting changes...\033[0m"
	@git commit -m "$(MSG)"

push: commit
	echo "\033[32mpushing to remote repo...\033[0m"
	@git push git@github.com:rag-nar1/github.com/rag-nar1/TCP-Duckdb.git

clean:
	rm -rf $(BUILD_DIR)/*

.PHONY: build run clean