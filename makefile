runGo:
	go mod tidy
	go run ./cmd/server/main.go ./cmd/server/top.go

runUi:
	$(SHELL) cd ./pkg/ui
	npm run dev