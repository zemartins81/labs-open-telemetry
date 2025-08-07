# Sistema de Temperatura por CEP com OpenTelemetry

Sistema distribuído em Go composto por dois microserviços que fornecem informações de temperatura baseadas em CEP brasileiro, integrado com OpenTelemetry e Zipkin para observabilidade.

## 📋 Visão Geral

- **Serviço A** (8080): Recebe e valida CEPs, encaminha para o Serviço B
- **Serviço B** (8081): Busca localização (ViaCEP) e temperatura (WeatherAPI)
- **Zipkin** (9411): Interface de observabilidade para traces distribuídos

## 🏗️ Arquitetura

```
┌─────────────┐    HTTP POST     ┌─────────────┐
│             │   /cep/{cep}     │             │
│  Serviço A  ├─────────────────▶│  Serviço B  │
│   (8080)    │                  │   (8081)    │
│             │                  │             │
└─────────────┘                  └──────┬──────┘
       │                                │
       │                      ┌─────────▼─────────┐
       │                      │                   │
       └──────────────────────┤   APIs Externas   │
              OpenTelemetry    │                   │
              Traces           │ ViaCEP + Weather  │
              ▼                └───────────────────┘
┌─────────────────────────┐
│                         │
│    Zipkin (9411)        │
│   Trace Visualization   │
└─────────────────────────┘
```

## 🚀 Início Rápido

### Pré-requisitos

1. **Docker & Docker Compose** (obrigatório)
2. **Chave WeatherAPI** (gratuita)
3. **Go 1.24+** (opcional, para desenvolvimento local)

### 1. Clone o Repositório

```bash
git clone https://github.com/zemartins81/labs-open-telemetry.git
cd labs-open-telemetry/sistema-cep
```

### 2. Obter Chave da WeatherAPI

1. Acesse [WeatherAPI.com](https://www.weatherapi.com/)
2. Crie uma conta gratuita
3. Vá em "My Account" → "API Keys"
4. Copie sua chave

### 3. Configurar Ambiente

```bash
# Copie o arquivo de exemplo
cp .env.example .env

# Edite o arquivo .env e adicione sua chave
# WEATHER_API_KEY=sua_chave_weatherapi_aqui
```

### 4. Executar com Docker Compose

```bash
# Subir todos os serviços (microserviços + Zipkin)
docker compose up --build

# Em background (recomendado)
docker compose up --build -d
```

### 5. Verificar se os Serviços Estão Funcionando

```bash
# Verificar status dos containers
docker compose ps

# Health checks
curl http://localhost:8080/health  # Serviço A
curl http://localhost:8081/health  # Serviço B
curl http://localhost:9411/health  # Zipkin
```

### 6. Testar o Sistema

```bash
# Teste com CEP válido (Av. Paulista - São Paulo)
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'

# Resposta esperada:
# {
#   "city": "San Paulo",
#   "temp_C": 23.1,
#   "temp_F": 73.6,
#   "temp_K": 296.25
# }
```

### 7. Acessar Interface de Observabilidade

- **Zipkin**: [http://localhost:9411](http://localhost:9411)
  - Visualize traces distribuídos entre os serviços
  - Analise latência e performance das chamadas

## 📁 Estrutura do Projeto

```text
labs-open-telemetry/
├── sistema-cep/
│   ├── servico-a/                  # Microserviço A (validação e roteamento)
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   ├── Dockerfile
│   │   ├── makefile
│   │   └── README.md
│   ├── servico-b/                  # Microserviço B (orquestração de APIs)
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   ├── Dockerfile
│   │   ├── makefile
│   │   └── README.md
│   ├── docker-compose.yml          # Orquestração dos serviços
│   ├── .env.example               # Exemplo de configuração
│   └── .env                       # Configuração local (criar)
└── README.md                      # Esta documentação
```

## 🔧 Desenvolvimento Local

### Executar Serviços Individualmente

```bash
# Terminal 1 - Serviço B (precisa executar primeiro)
cd servico-b
WEATHER_API_KEY=sua_chave make run

# Terminal 2 - Serviço A
cd servico-a
make run

# Terminal 3 - Zipkin (opcional para desenvolvimento)
docker run -d -p 9411:9411 openzipkin/zipkin
```

### Executar Testes

```bash
# Serviço A
cd servico-a && make test

# Serviço B
cd servico-b && make test

# Teste funcional completo (requer serviços rodando)
make test-integration
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
  "city": "San Paulo",
  "temp_C": 23.1,
  "temp_F": 73.6,
  "temp_K": 296.25
}
```

**Response Error (422):**

```json
{
  "message": "invalid zipcode"
}
```

#### `GET /health`

Health check do Serviço A

### Serviço B (Porta 8081)

#### `POST /temperature`

Processa CEP e busca temperatura (uso interno)

#### `GET /health`

Health check do Serviço B

### Zipkin (Porta 9411)

#### Interface Web

- **URL**: [http://localhost:9411](http://localhost:9411)
- **Função**: Visualização de traces distribuídos
- **Uso**: Análise de performance e debugging

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
docker compose down

# Rebuild completo (limpa cache)
docker compose down
docker compose build --no-cache
docker compose up

# Ver logs em tempo real
docker compose logs -f servico-a
docker compose logs -f servico-b
docker compose logs -f zipkin

# Ver logs de todos os serviços
docker compose logs -f

# Escalar serviços (para testes de carga)
docker compose up --scale servico-a=2 --scale servico-b=2
```

## 🔧 Configuração Avançada

### Variáveis de Ambiente

| Variável | Serviço | Descrição | Padrão | Obrigatória |
|----------|---------|-----------|--------|-------------|
| `WEATHER_API_KEY` | B | Chave WeatherAPI | - | ✅ Sim |
| `SERVICE_B_URL` | A | URL do Serviço B | `http://servico-b:8081` | Não |
| `VIACEP_URL` | B | URL da ViaCEP | `https://viacep.com.br/ws` | Não |
| `WEATHER_API_URL` | B | URL da WeatherAPI | `http://api.weatherapi.com/v1` | Não |
| `ZIPKIN_ENDPOINT` | A, B | URL do Zipkin | `http://zipkin:9411/api/v2/spans` | Não |
| `PORT` | A, B | Porta do serviço | `8080`/`8081` | Não |

### Portas Customizadas

Para alterar as portas, crie um arquivo `docker-compose.override.yml`:

```yaml
version: '3.8'
services:
  servico-a:
    ports:
      - "9000:8080"
  servico-b:
    ports:
      - "9001:8081"
  zipkin:
    ports:
      - "9412:9411"
```

## 🐛 Troubleshooting

### ❌ Erro: "WEATHER_API_KEY é obrigatória"

**Solução:**

1. Verifique se a chave está no arquivo `.env`
2. Confirme se a chave é válida no [WeatherAPI.com](https://www.weatherapi.com/)
3. Reinicie os containers: `docker compose restart`

### ❌ Erro: "can not find zipcode"

**Problema:** CEP pode não existir na base ViaCEP

**Solução:**

- Teste com CEPs conhecidos:
  - `01310100` (Av. Paulista, São Paulo)
  - `20040020` (Centro, Rio de Janeiro)
  - `40070110` (Pelourinho, Salvador)

### ❌ Serviços não se comunicam

**Verificações:**

1. Todos os containers estão rodando: `docker compose ps`
2. Containers estão na mesma rede Docker
3. URLs de comunicação estão corretas no docker-compose.yml

### ❌ Zipkin não mostra traces

**Possíveis causas:**

1. Variável `ZIPKIN_ENDPOINT` incorreta
2. Zipkin não está rodando: `curl http://localhost:9411/health`
3. Firewall bloqueando comunicação

## 📊 Observabilidade

### Análise de Traces no Zipkin

1. **Acesse**: [http://localhost:9411](http://localhost:9411)
2. **Faça algumas requisições**:

   ```bash
   curl -X POST http://localhost:8080 -H "Content-Type: application/json" -d '{"cep": "01310100"}'
   ```

3. **No Zipkin**:
   - Clique em "Run Query" para ver traces
   - Analise a latência entre serviços
   - Identifique gargalos de performance

### Métricas Disponíveis

- **Tempo de resposta** por serviço
- **Latência de rede** entre microserviços
- **Tempo de resposta das APIs externas** (ViaCEP, WeatherAPI)
- **Traces de erro** para debugging

## 📚 Próximos Passos

- [ ] **Metrics**: Implementação Prometheus + Grafana
- [ ] **Logging**: Structured logs com ELK Stack
- [ ] **Cache**: Redis para cache de respostas ViaCEP
- [ ] **Rate Limiting**: Controle de taxa de requisições
- [ ] **Circuit Breaker**: Proteção contra falhas em cascata
- [ ] **Health Checks**: Checks mais robustos com dependências
- [ ] **Load Testing**: Testes de carga automatizados

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature: `git checkout -b feature/nova-funcionalidade`
3. Commit suas mudanças: `git commit -m 'Add nova funcionalidade'`
4. Push para a branch: `git push origin feature/nova-funcionalidade`
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

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
