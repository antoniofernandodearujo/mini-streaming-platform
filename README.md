## 🎥 Mini Streaming Platform
Este é um projeto de estudo para explorar o uso de processos paralelos, transcodificação de vídeos e integração com Amazon S3. A plataforma permite enviar vídeos para a nuvem de forma eficiente, transcodificá-los utilizando Go e FFmpeg, e disponibilizar os vídeos para streaming com HLS (HTTP Live Streaming).

### 🚀 Objetivos do Projeto
Este projeto foi desenvolvido para entender de forma prática e aprofundada como funcionam os seguintes conceitos:

- Processamento paralelo com goroutines em Go para transcodificação de vídeos.
- Envio de vídeos em chunks para o Amazon S3, garantindo eficiência no upload de arquivos grandes.
- Utilização do HLS (HTTP Live Streaming) para a entrega de vídeos em diferentes qualidades.
- Integração com o FFmpeg para transcodificação de vídeos em múltiplos formatos e resoluções.
  
### 🛠️ Tecnologias Utilizadas
- Backend: Go, Goroutines, FFmpeg, S3 (Amazon Web Services), Docker
- Frontend: React
- Transcodificação de Vídeo: FFmpeg, HLS
- Armazenamento de Vídeos: Amazon S3

### ⚡ Principais Funcionalidades
- Envio de Vídeos em Chunks: Os vídeos são divididos em pequenas partes (chunks) e enviados para o Amazon S3, o que garante um upload mais eficiente e confiável, especialmente para arquivos grandes.
- Transcodificação Paralela com Go: Utilizando as goroutines do Go, a plataforma é capaz de processar vídeos de maneira paralela e eficiente, aproveitando múltiplos núcleos de processamento.
- HLS para Streaming: Os vídeos são transcodificados em várias resoluções e entregues ao usuário por meio de HLS, garantindo uma experiência fluida e de alta qualidade, independentemente da conexão.
- Frontend de Streaming com React: Interface simples e responsiva para assistir aos vídeos, com suporte a múltiplos formatos e resoluções.

### 🧑‍💻 Como Executar o Projeto
1. Clone o Repositório
```bash
git clone https://github.com/antoniofernandodearujo/mini-streaming-platform.git
cd mini-streaming-platform
```
2. Backend
Para rodar o backend, você precisará ter o Go e o FFmpeg instalados.

Instale as dependências do Go:
```bash
cd backend
go mod tidy
```
Inicie o servidor backend:
```bash
docker compose up -d --build
```

3. Frontend
Para rodar o frontend, você precisará do Node.js instalado.

Instale as dependências do frontend:
```bash
cd frontend/streaming-video
npm install
```
Inicie o servidor frontend:
```bash
npm start
```

4. Configuração do AWS S3
Certifique-se de ter as credenciais do AWS S3 configuradas corretamente.
Crie um bucket no S3 para armazenar os vídeos e configure as permissões necessárias.

### 🔧 Desafios e Aprendizados
- Transcodificação de Vídeos com FFmpeg: Durante o desenvolvimento, foi necessário entender como o FFmpeg pode ser usado para transcodificar vídeos em diferentes resoluções e formatos.
- Processamento Paralelo com Go: A utilização de goroutines no Go foi um aprendizado valioso sobre como otimizar o uso de múltiplos núcleos de processamento e realizar tarefas de forma paralela.
- Integração com AWS S3: O envio de vídeos em chunks para o S3 foi um desafio técnico que me permitiu aprender mais sobre o gerenciamento de grandes arquivos na nuvem.

### 🌱 O Que Aprendi
- Como trabalhar com goroutines em Go para otimização de processos paralelos.
- Como integrar FFmpeg e HLS para transcodificação e entrega de vídeos em tempo real.
- Como otimizar o envio de vídeos para Amazon S3 utilizando upload de chunks.

### 📄 Licença
Este projeto está licenciado sob a MIT License.

