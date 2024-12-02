# Use a imagem oficial do Go como base
FROM golang:1.22

# Atualizar e instalar FFmpeg e outras dependências necessárias
RUN apt-get update && \
    apt-get install -y ffmpeg && \
    rm -rf /var/lib/apt/lists/*  # Limpar cache do apt-get para reduzir o tamanho da imagem

# Configurar diretório de trabalho
WORKDIR /app

# Copiar todos os arquivos do projeto para o container
COPY . .

# Copiar o arquivo .env para o diretório de trabalho no container
# COPY .env .env


# Cria diretórios necessários e define permissões
# Cria diretórios no /tmp e define permissões
RUN mkdir -p /tmp/videos/hls /tmp/videos/videos && chmod -R 755 /tmp/videos

# Instalar dependências Go
RUN go mod tidy

# Compilar o binário Go
RUN go build -o main cmd/main.go

# Expor a porta do servidor
EXPOSE 8080

# Rodar o servidor Go
CMD ["sh", "-c", "mkdir -p /tmp/videos/hls /tmp/videos/videos && chmod -R 755 /tmp/videos && ./main"]
