# Serviço A - Sistema de Temperatura por CEP

O Serviço A é responsável por receber requisições com CEP, validar o formato e encaminhar para o Serviço B.

## Estrutura do Projeto

```
servico-a/
├── cmd/
│   └── main.go                 # Ponto de entrada da aplicação
├── internal/
│   ├── config/
│   │   └── config.go           # Configurações da aplicação
│   ├── handlers/
│   │   └── cep_handler.go      # Handlers HTTP
│   ├── models/
│   │   └── models.go           # Estruturas de dados
│   ├── server/
│   │   └── server.go           # Configuração do servidor HTTP
│   ├── services/
│   │   └── serviceb_client.go  # Cliente para comunicação com Serviço B
│   └── validators/
│       ├── cep_validator.go    # Validações de CEP
│       └── cep_validator_test.go # Testes unitários
├── Dockerfile                  # Configuração Docker
├── Makefile                   # Scripts de automação
├── go.mod                     # Dependências Go
├── .gitignore                 # Arquivos ignorados pelo Git
└── README.md                  # Documentação
```

## Funcionalidades

- Recebe CEP via POST no formato JSON
- Valida se o CEP tem exatamente 8 dígitos numéricos
- Encaminha requisições válidas para o Serviço B
- Retorna erros apropriados para CEPs inválidos
- Endpoint de health check

## Endpoints

### POST /
Recebe e processa requisições de CEP

**Request:**
```json
{
  "cep": "29902555"
}
```

**Responses:**

#### Sucesso (200)
```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

#### CEP Inválido (422)
```json
{
  "message": "invalid zipcode"
}
```

### GET /health
Verifica a saúde da aplicação

**Response (200):**
```json
{
  "status": "healthy",
  "service": "servico-a"
}
```

## Variáveis de Ambiente

- `PORT`: Porta do servidor (padrão: 8080)
- `SERVICE_B_URL`: URL do Serviço B (padrão: http://localhost:8081)

## Como Executar

### Pré-requisitos
- Go 1.21 ou superior
- Docker (opcional)

### Desenvolvimento Local

```bash
# Executar diretamente
make run

# Ou usando go run
go run ./cmd

# Executar com variáveis customizadas
PORT=9000 SERVICE_B_URL=http://localhost:8082 make run
```

### Build da Aplicação

```bash
# Compilar
make build

# Executar binário
./bin/servico-a
```

### Docker

```bash
# Construir imagem
make docker-build

# Executar container
make docker-run

# Ou manualmente
docker build -t servico-a .
docker run -p 8080:8080 \
  -e SERVICE_B_URL=http://servico-b:8081 \
  servico-a
```

## Testes

```bash
# Executar testes
make test

# Testes com cobertura
make test-coverage

# Formatar código
make fmt

# Executar linter (requer instalação)
make lint
```

## Exemplos de Uso

```bash
# Teste com CEP válido
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"cep": "29902555"}'

# Teste com CEP inválido
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"cep": "123"}'

# Health check
curl http://localhost:8080/health
```

## Arquitetura

O projeto segue os princípios de Clean Architecture:

- **cmd/**: Ponto de entrada da aplicação
- **internal/config/**: Gerenciamento de configurações
- **internal/handlers/**: Camada de apresentação (HTTP handlers)
- **internal/services/**: Camada de serviços (regras de negócio)
- **internal/models/**: Estruturas de dados
- **internal/validators/**: Validações de entrada
- **internal/server/**: Configuração do servidor HTTP

## Próximos Passos

- Implementação do OpenTelemetry para tracing distribuído
- Integração com Zipkin
- Métricas de observabilidade
- Docker Compose para ambiente completo
