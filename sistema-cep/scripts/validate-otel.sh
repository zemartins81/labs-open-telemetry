#!/bin/bash

# =============================================================================
# Script de Validação do OTEL Collector
# Sistema de Temperatura por CEP
# =============================================================================

set -e

echo "🔍 Validando configuração do OTEL Collector..."

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para log com cores
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Função para verificar se um serviço está rodando
check_service() {
    local service_name=$1
    local port=$2
    local path=${3:-""}

    log_info "Verificando $service_name na porta $port..."

    if curl -s --connect-timeout 5 "http://localhost:${port}${path}" > /dev/null; then
        log_success "$service_name está respondendo"
        return 0
    else
        log_error "$service_name não está respondendo"
        return 1
    fi
}

# Função para verificar configuração YAML
validate_yaml() {
    local file=$1
    log_info "Validando sintaxe YAML: $file"

    if command -v yamllint &> /dev/null; then
        if yamllint "$file" > /dev/null 2>&1; then
            log_success "Sintaxe YAML válida: $file"
        else
            log_error "Erro de sintaxe YAML: $file"
            yamllint "$file"
            return 1
        fi
    else
        log_warning "yamllint não encontrado, pulando validação de sintaxe"
    fi
}

# Função para testar endpoints do OTEL Collector
test_otel_endpoints() {
    log_info "Testando endpoints do OTEL Collector..."

    # Health check
    if curl -s "http://localhost:13133/health" | grep -q "Server available"; then
        log_success "Health check endpoint OK"
    else
        log_warning "Health check endpoint não disponível"
    fi

    # Métricas internas
    if curl -s "http://localhost:8888/metrics" | grep -q "otelcol_"; then
        log_success "Métricas internas disponíveis"
    else
        log_warning "Métricas internas não disponíveis"
    fi

    # zpages
    if curl -s "http://localhost:55679/debug/tracez" > /dev/null; then
        log_success "zpages endpoint disponível"
    else
        log_warning "zpages endpoint não disponível"
    fi
}

# Função para testar conectividade com Zipkin
test_zipkin_connectivity() {
    log_info "Testando conectividade com Zipkin..."

    if curl -s "http://localhost:9411/health" | grep -q "OK"; then
        log_success "Zipkin está acessível"

        # Testar endpoint de spans
        if curl -s -X POST "http://localhost:9411/api/v2/spans" \
           -H "Content-Type: application/json" \
           -d '[]' > /dev/null; then
            log_success "Endpoint de spans do Zipkin OK"
        else
            log_warning "Endpoint de spans do Zipkin pode ter problemas"
        fi
    else
        log_error "Zipkin não está acessível"
    fi
}

# Função para enviar trace de teste
send_test_trace() {
    log_info "Enviando trace de teste via OTLP HTTP..."

    local test_trace='{
        "resourceSpans": [{
            "resource": {
                "attributes": [{
                    "key": "service.name",
                    "value": {"stringValue": "test-service"}
                }]
            },
            "scopeSpans": [{
                "scope": {
                    "name": "test-scope"
                },
                "spans": [{
                    "traceId": "5B8EFFF798038103D269B633813FC60C",
                    "spanId": "EEE19B7EC3C1B174",
                    "name": "test-span",
                    "kind": 1,
                    "startTimeUnixNano": "'$(date +%s)'000000000",
                    "endTimeUnixNano": "'$(($(date +%s) + 1))'000000000",
                    "attributes": [{
                        "key": "test.type",
                        "value": {"stringValue": "validation"}
                    }]
                }]
            }]
        }]
    }'

    if curl -s -X POST "http://localhost:4318/v1/traces" \
       -H "Content-Type: application/json" \
       -d "$test_trace" > /dev/null; then
        log_success "Trace de teste enviado com sucesso"
        log_info "Verifique o Zipkin UI em http://localhost:9411"
    else
        log_error "Falha ao enviar trace de teste"
    fi
}

# Função principal
main() {
    echo "🚀 Iniciando validação do OTEL Collector..."
    echo "================================================"

    # 1. Validar arquivos de configuração
    if [ -f "configs/otel-collector.yml" ]; then
        validate_yaml "configs/otel-collector.yml"
    else
        log_error "Arquivo configs/otel-collector.yml não encontrado"
        exit 1
    fi

    if [ -f "configs/prometheus.yml" ]; then
        validate_yaml "configs/prometheus.yml"
    fi

    echo ""

    # 2. Verificar se os serviços estão rodando
    log_info "Verificando serviços..."

    services_ok=true

    # OTEL Collector
    if ! check_service "OTEL Collector (Health)" 13133 "/health"; then
        services_ok=false
    fi

    if ! check_service "OTEL Collector (Metrics)" 8888 "/metrics"; then
        services_ok=false
    fi

    # Zipkin
    if ! check_service "Zipkin" 9411 "/health"; then
        services_ok=false
    fi

    # Verificar se OTLP endpoints estão listening
    if ss -tuln | grep -q ":4317"; then
        log_success "OTLP gRPC endpoint (4317) está ouvindo"
    else
        log_error "OTLP gRPC endpoint (4317) não está ouvindo"
        services_ok=false
    fi

    if ss -tuln | grep -q ":4318"; then
        log_success "OTLP HTTP endpoint (4318) está ouvindo"
    else
        log_error "OTLP HTTP endpoint (4318) não está ouvindo"
        services_ok=false
    fi

    echo ""

    # 3. Testes específicos do OTEL
    if [ "$services_ok" = true ]; then
        test_otel_endpoints
        echo ""

        test_zipkin_connectivity
        echo ""

        # 4. Enviar trace de teste
        if [ "${1:-}" = "--test-trace" ]; then
            send_test_trace
            echo ""
        fi

        # 5. Mostrar informações úteis
        echo "📊 Informações Úteis:"
        echo "- OTEL Collector Health: http://localhost:13133/health"
        echo "- OTEL Collector Metrics: http://localhost:8888/metrics"
        echo "- OTEL Collector Debug: http://localhost:55679"
        echo "- Zipkin UI: http://localhost:9411"
        echo "- OTLP gRPC: localhost:4317"
        echo "- OTLP HTTP: localhost:4318"
        echo ""

        log_success "Validação concluída com sucesso!"
        echo ""
        echo "💡 Para enviar um trace de teste, execute:"
        echo "   $0 --test-trace"

    else
        log_error "Alguns serviços não estão funcionando corretamente"
        echo ""
        echo "🔧 Sugestões:"
        echo "1. Verifique se o docker-compose está rodando: docker-compose ps"
        echo "2. Verifique os logs: docker-compose logs otel-collector"
        echo "3. Reinicie os serviços: docker-compose restart"
        exit 1
    fi
}

# Executar função principal
main "$@"
