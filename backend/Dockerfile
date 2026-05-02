FROM golang:1.25-alpine

WORKDIR /app

# Dependências do sistema
RUN apk add --no-cache git bash build-base

# Instala ferramentas de desenvolvimento
RUN go install github.com/air-verse/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

# Dependências do projeto
COPY go.mod go.sum ./
RUN go mod download

# Código-fonte
COPY . .

# Portas:
# 8080 -> aplicação
# 2345 -> Delve debugger
EXPOSE 8080 2345

# Modo padrão: desenvolvimento com hot reload
CMD ["air", "-c", ".air.toml"]