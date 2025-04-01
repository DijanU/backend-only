FROM golang:1.24.1-alpine

# Instalar dependencias (SQLite y compiladores)
RUN apk add --no-cache gcc musl-dev sqlite sqlite-dev

# Directorio de trabajo
WORKDIR /site

# Copiar primero los archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Compilar
RUN go build -o main .

EXPOSE 8080
CMD ["./main"]