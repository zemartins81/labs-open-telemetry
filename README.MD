# Sistema de Temperatura por CEP

Sistema distribuído em Go composto por dois microserviços que fornecem informações de temperatura baseadas em CEP brasileiro.

## 📋 Visão Geral

- **Serviço A**: Recebe e valida CEPs, encaminha para o Serviço B
- **Serviço B**: Busca localização (ViaCEP) e temperatura (WeatherAPI)

## 🏗️ Arquitetura

```
┌─────────────┐    HTTP POST     ┌─────────────┐
│             │   /{cep}         │             │
│  Serviço A  ├─────────────────▶│  Serviço B  │
│   (8080)    │                  │   (8081)    │
│             │                  │             │
└─────────────┘                  └──────┬──────┘
                                        │
                              ┌─────────▼─────────┐
                              │                   │
                              │   APIs Externas   │
                              │                   │
                              │ ViaCEP + Weather  │
                              └───────────────────┘
```

## 🚀 Início Rápido

### Pré-requisitos

1. **Go 1.21+** (para desenvolvimento local)
2. **Docker & Docker Compose** (recomendado)
3. **Chave WeatherAPI** (gratuita)

### 1. Obter Chave da WeatherAPI

1. Acesse [WeatherAPI.com](https://www.weatherapi.com/)
2. Crie uma conta gratuita
3. Vá em "My Account" → "API Keys"
4. Copie sua chave

### 2. Configurar Ambiente

```bash
# Clone ou configure os arquivos
cp .env.example .env

# Edite o arquivo .env e adicione sua chave
WEATHER_API_KEY=sua_chave_aqui
```

### 3. Executar com Docker Compose

```bash
# Subir todos os serviços
docker-compose up --build

# Em background
docker-compose up --build -d
```

### 4. Testar o Sistema

```bash
# Teste com CEP válido (Av. Paulista - São Paulo)
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'

# Resposta esperada:
# {
#   "city": "São Paulo",
#   "temp_C": 23.5,
#   "temp_F": 74.3,
#   "temp_K": 296.65
# }
```

## 📁 Estrutura do Projeto

```
sistema-temperatura-cep/
├── servico-a/                  # Microserviço A (validação e roteamento)
│   ├── cmd/main.go
│   ├── internal/
│   ├── Dockerfile
│   ├── Makefile
│   └── README.md
├── servico-b/                  # Microserviço B (orquestração de APIs)
│   ├── cmd/main.go
│   ├── internal/
│   ├── Dockerfile
│   ├── Makefile
│   └── README.md
├── docker-compose.yml          # Orquestração dos serviços
├── .env.example               # Exemplo de configuração
└── README.md                  # Esta documentação
```

## 🔧 Desenvolvimento Local

### Executar Serviços Individualmente

```bash
# Terminal 1 - Serviço B
cd servico-b
WEATHER_API_KEY=sua_chave make run

# Terminal 2 - Serviço A
cd servico-a
make run
```

### Executar Testes

```bash
# Serviço A
cd servico-a && make test

# Serviço B
cd servico-b && make test
```

## 📡 Endpoints da API

### Serviço A (Porta 8080)

#### `POST /`

Recebe CEP e retorna temperatura

**Request:**

```json
{
  "cep": "01310100"
}
```

**Response Success (200):**

```json
{
  "city": "São Paulo",
  "temp_C": 23.5,
  "temp_F": 74.3,
  "temp_K": 296.65
}
```

**Response Error (422):**

```json
{
  "message": "invalid zipcode"
}
```

#### `GET /health`

Health check do serviço

### Serviço B (Porta 8081)

#### `POST /temperature`

Processa CEP e busca temperatura

#### `GET /health`

Health check do serviço

## 🧪 Casos de Teste

```bash
# CEP válido
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'

# CEP inválido (formato)
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"cep": "123"}'

# CEP não encontrado
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"cep": "99999999"}'

# Health checks
curl http://localhost:8080/health
curl http://localhost:8081/health
```

## 🔍 Códigos de Resposta

| Código | Descrição |
|--------|-----------|
| 200    | Sucesso - temperatura encontrada |
| 400    | Bad Request - JSON malformado |
| 404    | CEP não encontrado |
| 405    | Método não permitido |
| 422    | CEP com formato inválido |
| 500    | Erro interno do servidor |

## 🌡️ Conversões de Temperatura

- **Celsius**: Valor original da WeatherAPI
- **Fahrenheit**: F = C × 1.8 + 32
- **Kelvin**: K = C + 273.15

## 🛠️ Comandos Úteis

```bash
# Parar todos os serviços
docker-compose down

# Rebuild completo
docker-compose down
docker-compose build --no-cache
docker-compose up

# Ver logs
docker-compose logs -f servico-a
docker-compose logs -f servico-b

# Escalar serviços
docker-compose up --scale servico-a=2
```

## 🔧 Configuração Avançada

### Variáveis de Ambiente

| Variável | Serviço | Descrição | Padrão |
|----------|---------|-----------|--------|
| `WEATHER_API_KEY` | B | Chave WeatherAPI (obrigatória) | - |
| `SERVICE_B_URL` | A | URL do Serviço B | `http://localhost:8081` |
| `VIACEP_URL` | B | URL da ViaCEP | `https://viacep.com.br/ws` |
| `WEATHER_API_URL` | B | URL da WeatherAPI | `http://api.weatherapi.com/v1` |

### Portas Customizadas

```bash
# docker-compose.override.yml
version: '3.8'
services:
  servico-a:
    ports:
      - "9000:8080"
  servico-b:
    ports:
      - "9001:8081"
```

## 🐛 Troubleshooting

### Erro: "WEATHER_API_KEY é obrigatória"

- Verifique se a chave está no arquivo `.env`
- Confirme se a chave é válida no WeatherAPI

### Erro: "can not find zipcode"

- CEP pode não existir na base ViaCEP
- Teste com CEPs conhecidos (01310100, 20040020)

### Serviços não se comunicam

- Verifique se estão na mesma rede Docker
- Confirme URLs de comunicação entre serviços

## 📚 Próximos Passos

- [ ] Implementação OpenTelemetry
- [ ] Integração Zipkin
- [ ] Métricas Prometheus
- [ ] Cache Redis
- [ ] Rate Limiting
- [ ] Circuit Breaker

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature
3. Implemente testes
4. Submeta um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT.
