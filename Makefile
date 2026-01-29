.PHONY: lint format check

# Menjalankan linter secara manual
lint:
	golangci-lint run ./...

# Memperbaiki format kode otomatis
format:
	go fmt ./...

# Menjalankan semua pengecekan sebelum commit
check: format lint
	go test ./...