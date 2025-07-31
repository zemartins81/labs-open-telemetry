.PHONY: build run test clean docker-build docker-run

# Variáveis
APP_NAME=servico-a
DOCKER_IMAGE=servico-a:latest
PORT=8080

# Build da aplicação
build:
	go build -o bin/$(APP_NAME) ./cmd

# Executa a aplicação localmente
run:
	go run ./cmd

# Executa os testes
test:
	go test -v ./...

# Executa os testes com coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Limpa arquivos gerados
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Formata o código
fmt:
	go fmt ./...

# Executa o linter
lint:
	golangci-lint run

# Build da imagem Docker
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Executa com Docker
docker-run:
	docker run -p $(PORT):$(PORT) \
		-e PORT=$(PORT) \
		-e SERVICE_B_URL=http://host.docker.internal:8081 \
		$(DOCKER_IMAGE)

# Executa em modo desenvolvimento
dev:
	go run ./cmd

# Instala dependências de desenvolvimento
dev-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Help
help:
	@echo "Comandos disponíveis:"
	@echo "  build          - Compila a aplicação"
	@echo "  run            - Executa a aplicação localmente"
	@echo "  test           - Executa os testes"
	@echo "  test-coverage  - Executa os testes com relatório de cobertura"
	@echo "  clean          - Remove arquivos gerados"
	@echo "  fmt            - Formata o código"
	@echo "  lint           - Executa o linter"
	@echo "  docker-build   - Constrói a imagem Docker"
	@echo "  docker-run     - Executa com Docker"
	@echo "  dev            - Executa em modo desenvolvimento"
	@echo "  dev-deps       - Instala dependências de desenvolvimento"
