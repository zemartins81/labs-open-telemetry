# Build stage
FROM golang:1.24-alpine AS builder

# Instala dependências necessárias
RUN apk add --no-cache git ca-certificates

# Define diretório de trabalho
WORKDIR /app

# Copia os arquivos go.mod e go.sum
COPY go.mod go.sum ./

# Copia o código fonte
COPY . .

# Gera o go.sum e baixa as dependências
RUN go mod tidy && go mod download

# Compila a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Production stage
FROM alpine:latest

# Instala ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copia o binário da stage de build
COPY --from=builder /app/main .

# Expõe a porta 8081
EXPOSE 8081

# Define variáveis de ambiente padrão
ENV PORT=8081
ENV WEATHER_API_KEY=""
ENV VIACEP_URL=https://viacep.com.br/ws
ENV WEATHER_API_URL=http://api.weatherapi.com/v1

# Comando para executar a aplicação
CMD ["./main"]
