# Etapa 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar archivos de modulos
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar codigo fuente
COPY . .

# Compilar la aplicacion
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Etapa 2: Runtime
FROM alpine:latest

WORKDIR /root/

# Copiar el binario compilado desde la etapa anterior
COPY --from=builder /app/main .

# Exponer el puerto
EXPOSE 8080

# Ejecutar la aplicacion
CMD ["./main"]
