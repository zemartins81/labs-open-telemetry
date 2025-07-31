Serviço B - Sistema de Temperatura por CEP
O Serviço B é responsável por receber CEPs válidos, buscar a localização via ViaCEP e consultar a temperatura via WeatherAPI.

Estrutura do Projeto
servico-b/
├── cmd/
│   └── main.go                      # Ponto de entrada da aplicação
├── internal/
│   ├── config/
│   │   └── config.go               # Configurações da aplicação
│   ├── handlers/
│   │   └── temperature_handler.go  # Handlers HTTP
│   ├── models/
│   │   └── models.go               # Estruturas de dados
│   ├── server/
│   │   └── server.go               # Configuração do servidor HTTP
│   ├── services/
│   │   ├── viacep_service.go       # Cliente ViaCEP
│   │   ├── weather_service.go      # Cliente WeatherAPI
│   │   └── temperature_service.go  # Orquestrador de serviços
│   └── validators/
│       └── cep_validator.go        # Validações de CEP
├── Dockerfile                      # Configuração Docker
├── Makefile                       # Scripts de automação
├── go.mod                         # Dependências Go
├── .gitignore                     # Arquivos ignorados pelo Git
└── README.md                      # Documentação
Funcionalidades
Recebe CEP via POST e valida formato
Busca localização usando API ViaCEP
Consulta temperatura usando API WeatherAPI
Converte temperaturas para Celsius, Fahrenheit e Kelvin
Tratamento completo de erros
Endpoint de health check
APIs Utilizadas
ViaCEP
URL: https://viacep.com.br/ws/{cep}/json/
Função: Buscar localização por CEP
Gratuita: Sim
WeatherAPI
URL: http://api.weatherapi.com/v1/current.json
Função: Consultar temperatura atual
Gratuita: Sim (com limitações)
Chave necessária: Sim
Pré-requisitos
1. Chave da WeatherAPI
Obtenha uma chave gratuita em WeatherAPI.com:

Crie uma conta
Vá para "My Account" → "API Keys"
Copie sua chave
Configure a variável WEATHER_API_KEY
Endpoints
POST /temperature
Processa CEP e retorna temperatura

Request:

json
{
  "cep": "01310100"
}
Responses:

Sucesso (200)
json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
CEP Inválido (422)
json
{
  "message": "invalid zipcode"
}
CEP Não Encontrado (404)
json
{
  "message": "can not find zipcode"
}
Erro Interno (500)
json
{
  "message": "internal server error"
}
GET /health
Verifica a saúde da aplicação

Response (200):

json
{
  "status": "healthy",
  "service": "servico-b"
}
Variáveis de Ambiente
Variável	Descrição	Padrão	Obrigatória
PORT	Porta do servidor	8081	Não
WEATHER_API_KEY	Chave da WeatherAPI	-	Sim
VIACEP_URL	URL base da ViaCEP	https://viacep.com.br/ws	Não
WEATHER_API_URL	URL base da WeatherAPI	http://api.weatherapi.com/v1	Não
Como Executar
Desenvolvimento Local
bash
# Com chave da WeatherAPI
WEATHER_API_KEY=sua_chave_aqui make run

# Ou exportando a variável
export WEATHER_API_KEY=sua_chave_aqui
make dev
Build da Aplicação
bash
# Compilar
make build

# Executar binário
WEATHER_API_KEY=sua_chave_aqui ./bin/servico-b
Docker
bash
# Construir imagem
make docker-build

# Executar container
WEATHER_API_KEY=sua_chave_aqui make docker-run

# Ou manualmente
docker build -t servico-b .
docker run -p 8081:8081 \
  -e WEATHER_API_KEY=sua_chave_aqui \
  servico-b
Testes
bash
# Executar testes unitários
make test

# Testes com cobertura
make test-coverage

# Teste funcional (requer serviço rodando)
WEATHER_API_KEY=sua_chave_aqui make test-cep
Exemplos de Uso
bash
# Teste com CEP válido (Paulista - São Paulo)
curl -X POST http://localhost:8081/temperature \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'

# Teste com CEP inválido
curl -X POST http://localhost:8081/temperature \
  -H "Content-Type: application/json" \
  -d '{"cep": "123"}'

# Teste com CEP não encontrado
curl -X POST http://localhost:8081/temperature \
  -H "Content-Type: application/json" \
  -d '{"cep": "99999999"}'

# Health check
curl http://localhost:8081/health
Fluxo de Processamento
Validação: Verifica se CEP tem 8 dígitos numéricos
ViaCEP: Busca localização (cidade/estado) pelo CEP
WeatherAPI: Consulta temperatura atual pela localização
Conversão: Calcula Kelvin a partir de Celsius (K = C + 273.15)
Resposta: Retorna todas as temperaturas formatadas
Tratamento de Erros
422: CEP com formato inválido
404: CEP não encontrado na ViaCEP
500: Erro na WeatherAPI ou problemas de conectividade
400: JSON malformado
405: Método HTTP não permitido
Logs
O serviço gera logs detalhados sobre:

CEPs processados
Chamadas para APIs externas
Temperaturas encontradas
Erros de processamento
Monitoramento
Endpoint /health para health checks
Logs estruturados para observabilidade
Timeouts configuráveis para APIs externas
Tratamento robusto de erros
Limitações
Depende da disponibilidade das APIs ViaCEP e WeatherAPI
WeatherAPI tem limite de requisições gratuitas
Apenas CEPs brasileiros são suportados
Temperaturas são sempre atuais (não históricas)
Próximos Passos
Implementação do OpenTelemetry para tracing distribuído
Integração com Zipkin
Cache de respostas para otimização
Métricas de performance
Rate limiting
