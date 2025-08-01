# Build stage
FROM golang:1.24-alpine AS builder

# Instala dependências necessárias
RUN apk add --no-cache git ca-certificates

# Define diretório de trabalho
WORKDIR /app

# Copia os arquivos go.mod e go.sum (se existir)
COPY go.mod go.sum* ./

# Baixa as dependências
RUN go mod download

# Copia o código fonte
COPY . .

# Compila a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Production stage
FROM alpine:latest

# Instala ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copia o binário da stage de build
COPY --from=builder /app/main .

# Expõe a porta 8080
EXPOSE 8080

# Define variáveis de ambiente padrão
ENV PORT=8080
ENV SERVICE_B_URL=http://localhost:8081

# Comando para executar a aplicação
CMD ["./main"]
