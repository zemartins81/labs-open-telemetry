.PHONY: build run test clean docker-build docker-run

# Variáveis
APP_NAME=servico-b
DOCKER_IMAGE=servico-b:latest
PORT=8081

# Build da aplicação
build:
	go build -o bin/$(APP_NAME) ./cmd

# Executa a aplicação localmente (requer WEATHER_API_KEY)
run:
	@if [ -z "$(WEATHER_API_KEY)" ]; then \
		echo "Erro: WEATHER_API_KEY é obrigatória"; \
		echo "Obtenha uma chave em https://www.weatherapi.com/"; \
		echo "Use: WEATHER_API_KEY=sua_chave make run"; \
		exit 1; \
	fi
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

# Executa com Docker (requer WEATHER_API_KEY)
docker-run:
	@if [ -z "$(WEATHER_API_KEY)" ]; then \
		echo "Erro: WEATHER_API_KEY é obrigatória"; \
		echo "Use: WEATHER_API_KEY=sua_chave make docker-run"; \
		exit 1; \
	fi
	docker run -p $(PORT):$(PORT) \
		-e PORT=$(PORT) \
		-e WEATHER_API_KEY=$(WEATHER_API_KEY) \
		$(DOCKER_IMAGE)

# Executa em modo desenvolvimento
dev:
	@if [ -z "$(WEATHER_API_KEY)" ]; then \
		echo "Erro: WEATHER_API_KEY é obrigatória"; \
		echo "Use: WEATHER_API_KEY=sua_chave make dev"; \
		exit 1; \
	fi
	go run ./cmd

# Testa com CEP real
test-cep:
	@if [ -z "$(WEATHER_API_KEY)" ]; then \
		echo "Erro: WEATHER_API_KEY é obrigatória para teste"; \
		exit 1; \
	fi
	curl -X POST http://localhost:8081/temperature \
		-H "Content-Type: application/json" \
		-d '{"cep": "01310-100"}'

# Help
help:
	@echo "Comandos disponíveis:"
	@echo "  build          - Compila a aplicação"
	@echo "  run            - Executa a aplicação localmente (requer WEATHER_API_KEY)"
	@echo "  test           - Executa os testes"
	@echo "  test-coverage  - Executa os testes com relatório de cobertura"
	@echo "  clean          - Remove arquivos gerados"
	@echo "  fmt            - Formata o código"
	@echo "  lint           - Executa o linter"
	@echo "  docker-build   - Constrói a imagem Docker"
	@echo "  docker-run     - Executa com Docker (requer WEATHER_API_KEY)"
	@echo "  dev            - Executa em modo desenvolvimento"
	@echo "  test-cep       - Testa com CEP real"
	@echo ""
	@echo "Exemplo de uso:"
	@echo "  WEATHER_API_KEY=sua_chave_aqui make run"
