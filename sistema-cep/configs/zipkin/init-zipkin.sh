#!/bin/bash

# =============================================================================
# Script de Inicialização do Zipkin
# Sistema de Temperatura por CEP
# =============================================================================

set -e

# Variáveis de ambiente
ZIPKIN_PORT=${ZIPKIN_PORT:-9411}
STORAGE_TYPE=${STORAGE_TYPE:-mem}
JAVA_OPTS=${JAVA_OPTS:-"-Xms512m -Xmx1024m"}

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Função para verificar dependências
check_dependencies() {
    log_info "Verificando dependências..."

    # Verificar Java
    if ! command -v java &> /dev/null; then
        log_error "Java não encontrado"
        exit 1
    fi

    local java_version=$(java -version 2>&1 | grep version | awk '{print $3}' | sed 's/"//g')
    log_success "Java encontrado: $java_version"

    # Verificar curl
    if ! command -v curl &> /dev/null; then
        log_warning "curl não encontrado - health checks podem falhar"
    fi

    # Verificar nc (netcat)
    if ! command -v nc &> /dev/null; then
        log_warning "netcat não encontrado - verificações de porta podem falhar"
    fi
}

# Função para configurar JVM
setup_jvm() {
    log_info "Configurando JVM..."

    # Configurações básicas de JVM
    export JAVA_OPTS="$JAVA_OPTS \
        -server \
        -Djava.awt.headless=true \
        -XX:+UseG1GC \
        -XX:MaxGCPauseMillis=200 \
        -XX:+UseStringDeduplication \
        -XX:+OptimizeStringConcat \
        -Djava.security.egd=file:/dev/./urandom"

    # Configurações específicas do ambiente
    if [ "$STORAGE_TYPE" = "mem" ]; then
        export JAVA_OPTS="$JAVA_OPTS -XX:+UseCompressedOops"
        log_info "Configurações otimizadas para storage em memória"
    fi

    # Configurações de debug (apenas em desenvolvimento)
    if [ "${ENVIRONMENT:-}" = "development" ]; then
        export JAVA_OPTS="$JAVA_OPTS \
            -XX:+PrintGCDetails \
            -XX:+PrintGCTimeStamps \
            -Xloggc:/tmp/zipkin-gc.log"
        log_info "Configurações de debug habilitadas"
    fi

    log_success "JVM configurada: $JAVA_OPTS"
}

# Função para configurar storage
setup_storage() {
    log_info "Configurando storage: $STORAGE_TYPE"

    case "$STORAGE_TYPE" in
        "mem")
            export STORAGE_TYPE=mem
            export MEM_MAX_SPANS=${MEM_MAX_SPANS:-500000}
            log_info "Storage em memória configurado (max spans: $MEM_MAX_SPANS)"
            ;;
        "elasticsearch")
            export STORAGE_TYPE=elasticsearch
            export ES_HOSTS=${ES_HOSTS:-http://elasticsearch:9200}
            export ES_INDEX=${ES_INDEX:-zipkin}
            log_info "Storage Elasticsearch configurado (hosts: $ES_HOSTS)"
            ;;
        "mysql")
            export STORAGE_TYPE=mysql
            export MYSQL_HOST=${MYSQL_HOST:-mysql}
            export MYSQL_TCP_PORT=${MYSQL_TCP_PORT:-3306}
            export MYSQL_DB=${MYSQL_DB:-zipkin}
            export MYSQL_USER=${MYSQL_USER:-zipkin}
            export MYSQL_PASS=${MYSQL_PASS:-zipkin}
            log_info "Storage MySQL configurado (host: $MYSQL_HOST:$MYSQL_TCP_PORT)"
            ;;
        *)
            log_error "Tipo de storage não suportado: $STORAGE_TYPE"
            exit 1
            ;;
    esac
}

# Função para configurar coletor
setup_collector() {
    log_info "Configurando coletor..."

    # Configurações do coletor HTTP
    export COLLECTOR_HTTP_ENABLED=${COLLECTOR_HTTP_ENABLED:-true}
    export COLLECTOR_SAMPLE_RATE=${COLLECTOR_SAMPLE_RATE:-1.0}

    # Configurações Kafka (se habilitado)
    if [ "${KAFKA_ENABLED:-false}" = "true" ]; then
        export KAFKA_BOOTSTRAP_SERVERS=${KAFKA_BOOTSTRAP_SERVERS:-kafka:9092}
        export KAFKA_TOPIC=${KAFKA_TOPIC:-zipkin}
        log_info "Coletor Kafka habilitado (bootstrap: $KAFKA_BOOTSTRAP_SERVERS)"
    fi

    # Configurações RabbitMQ (se habilitado)
    if [ "${RABBITMQ_ENABLED:-false}" = "true" ]; then
        export RABBIT_ADDRESSES=${RABBIT_ADDRESSES:-rabbitmq:5672}
        export RABBIT_QUEUE=${RABBIT_QUEUE:-zipkin}
        log_info "Coletor RabbitMQ habilitado (addresses: $RABBIT_ADDRESSES)"
    fi

    log_success "Coletor configurado (HTTP: $COLLECTOR_HTTP_ENABLED, Sample Rate: $COLLECTOR_SAMPLE_RATE)"
}

# Função para aguardar dependências
wait_for_dependencies() {
    log_info "Aguardando dependências..."

    case "$STORAGE_TYPE" in
        "elasticsearch")
            wait_for_service "$ES_HOSTS" "Elasticsearch"
            ;;
        "mysql")
            wait_for_port "$MYSQL_HOST" "$MYSQL_TCP_PORT" "MySQL"
            ;;
    esac
}

# Função para aguardar um serviço HTTP
wait_for_service() {
    local url=$1
    local name=$2
    local timeout=${3:-60}
    local count=0

    log_info "Aguardando $name estar disponível ($url)..."

    while [ $count -lt $timeout ]; do
        if curl -s "$url" > /dev/null 2>&1; then
            log_success "$name está disponível"
            return 0
        fi
        count=$((count + 1))
        sleep 1
    done

    log_error "$name não ficou disponível em ${timeout}s"
    return 1
}

# Função para aguardar uma porta TCP
wait_for_port() {
    local host=$1
    local port=$2
    local name=$3
    local timeout=${4:-60}
    local count=0

    log_info "Aguardando $name estar disponível ($host:$port)..."

    while [ $count -lt $timeout ]; do
        if nc -z "$host" "$port" 2>/dev/null; then
            log_success "$name está disponível"
            return 0
        fi
        count=$((count + 1))
        sleep 1
    done

    log_error "$name não ficou disponível em ${timeout}s"
    return 1
}

# Função para verificar saúde do Zipkin
health_check() {
    local timeout=${1:-30}
    local count=0

    log_info "Verificando saúde do Zipkin..."

    while [ $count -lt $timeout ]; do
        if curl -s "http://localhost:$ZIPKIN_PORT/health" | grep -q "OK"; then
            log_success "Zipkin está saudável"
            return 0
        fi
        count=$((count + 1))
        sleep 1
    done

    log_error "Zipkin não passou no health check em ${timeout}s"
    return 1
}

# Função para criar dados de exemplo (desenvolvimento)
create_sample_data() {
    if [ "${CREATE_SAMPLE_DATA:-false}" = "true" ]; then
        log_info "Criando dados de exemplo..."

        # Trace de exemplo
        local sample_trace='[{
            "traceId": "463ac35c9f6413ad48485a3953bb6124",
            "id": "a2fb4a1d1a96d312",
            "name": "get-traces",
            "timestamp": '$(date +%s)'000000,
            "duration": 386000,
            "localEndpoint": {
                "serviceName": "sistema-cep-example",
                "ipv4": "127.0.0.1",
                "port": 8080
            },
            "tags": {
                "http.method": "GET",
                "http.path": "/api/traces"
            }
        }]'

        if curl -s -X POST "http://localhost:$ZIPKIN_PORT/api/v2/spans" \
           -H "Content-Type: application/json" \
           -d "$sample_trace" > /dev/null; then
            log_success "Dados de exemplo criados"
        else
            log_warning "Falha ao criar dados de exemplo"
        fi
    fi
}

# Função principal
main() {
    echo "🚀 Inicializando Zipkin Server..."
    echo "=================================="

    # 1. Verificar dependências
    check_dependencies

    # 2. Configurar JVM
    setup_jvm

    # 3. Configurar storage
    setup_storage

    # 4. Configurar coletor
    setup_collector

    # 5. Aguardar dependências externas
    wait_for_dependencies

    echo ""
    log_info "Configuração concluída:"
    log_info "- Porta: $ZIPKIN_PORT"
    log_info "- Storage: $STORAGE_TYPE"
    log_info "- Java Opts: $JAVA_OPTS"
    echo ""

    # 6. Iniciar Zipkin (se não estiver sendo executado via docker)
    if [ "${START_ZIPKIN:-true}" = "true" ]; then
        log_info "Iniciando Zipkin Server..."

        # Download do Zipkin se não existir
        if [ ! -f "/app/zipkin.jar" ]; then
            log_info "Baixando Zipkin Server..."
            curl -sSL https://search.maven.org/remote_content?g=io.zipkin&a=zipkin-server&v=LATEST&c=exec -o /app/zipkin.jar
        fi

        # Executar Zipkin
        exec java $JAVA_OPTS -jar /app/zipkin.jar &
        ZIPKIN_PID=$!

        # Aguardar inicialização
        sleep 10

        # Verificar saúde
        if health_check; then
            create_sample_data
            log_success "Zipkin iniciado com sucesso (PID: $ZIPKIN_PID)"
            log_info "Acesse a UI em: http://localhost:$ZIPKIN_PORT"

            # Manter processo vivo
            wait $ZIPKIN_PID
        else
            log_error "Falha na inicialização do Zipkin"
            exit 1
        fi
    else
        log_info "Configuração concluída. Zipkin deve ser iniciado externamente."
    fi
}

# Tratar sinais para shutdown graceful
trap 'log_info "Recebido sinal de shutdown..."; kill -TERM $ZIPKIN_PID 2>/dev/null; exit 0' TERM INT

# Executar função principal
main "$@"