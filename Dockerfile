FROM golang:1.24.1-alpine

# Instalar dependencias (SQLite, compiladores y el driver)
RUN apk add --no-cache gcc musl-dev sqlite sqlite-dev

# Configuración específica para el driver SQLite
ENV CGO_ENABLED=1

WORKDIR /site

# Copiar módulos primero para cachear dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto
COPY . .

# Compilar (¡con CGO habilitado!)
RUN go build -ldflags="-s -w" -o main .

EXPOSE 8080
CMD ["./main"]