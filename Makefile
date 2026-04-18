BACKEND_DIR := backend

.PHONY: run seed

run:
	cd $(BACKEND_DIR) && go run .

seed:
	cd $(BACKEND_DIR) && go run ./cmd/seed/main.go
