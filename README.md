# 🎞️ Go Video Encoder Microservice

Este projeto é um microserviço de codificação de vídeo desenvolvido em Go, projetado para converter arquivos de vídeo MP4 em formato MPG. Ele utiliza o RabbitMQ para comunicação assíncrona e o Google Cloud Storage para armazenamento dos vídeos processados.

## 🚀 Funcionalidades

- Conversão automática de vídeos MP4 para MPG
- Processamento assíncrono de mensagens via RabbitMQ
- Armazenamento dos vídeos convertidos no Google Cloud Storage
- Gerenciamento de mensagens não processadas através de Dead Letter Exchange (DLX)

## ⚙️ Configuração do Ambiente

1. Duplique o arquivo `.env.example` e renomeie para `.env`
2. Execute o comando:

   ```bash
   docker-compose up -d
   ```

3. Acesse o painel de administração do RabbitMQ e:

- Crie uma exchange do tipo `fanout` para atuar como Dead Letter Exchange
- Crie uma fila de Dead Letter e vincule-a à exchange criada
- No arquivo `.env`, defina o nome da exchange no parâmetro `RABBITMQ_DLX`

4. No Google Cloud Platform (GCP):

- Crie uma conta de serviço com permissões para escrever no Google Cloud Storage
- Baixe o arquivo JSON com as credenciais e salve-o na raiz do projeto com o nome `bucket-credential.json`

## ▶️ Execução

Para iniciar o serviço de codificação, execute o seguinte comando dentro do contêiner:

```bash
docker exec app make server
```

`Nota: app é o nome do contêiner definido no docker-compose.`

## 📩 Formato da Mensagem de Entrada

As mensagens enviadas para o encoder devem estar no seguinte formato JSON:

```json
{
  "resource_id": "meu-id-de-recurso",
  "file_path": "caminho/do/video.mp4"
}
```

- resource_id: Identificador único do vídeo a ser convertido
- file_path: Caminho completo do arquivo MP4 no bucket

## 📤 Formato da Mensagem de Saída

### ✅ Sucesso

Após o processamento bem-sucedido, o encoder enviará uma mensagem com o seguinte formato:

```json
{
  "id": "uuid-do-processamento",
  "output_bucket_path": "bucket",
  "status": "COMPLETED",
  "video": {
    "encoded_video_folder": "pasta-do-video-convertido",
    "resource_id": "id-do-recurso",
    "file_path": "video.mp4"
  },
  "Error": "",
  "created_at": "2020-05-27T19:43:34.850479-04:00",
  "updated_at": "2020-05-27T19:43:38.081754-04:00"
}
```

### ❌ Erro

Em caso de falha no processamento, a mensagem de retorno será:

```json
{
  "message": {
    "resource_id": "id-do-recurso",
    "file_path": "video.mp4"
  },
  "error": "motivo do erro"
}
```

Além disso, a mensagem original será encaminhada para a Dead Letter Exchange definida no parâmetro RABBITMQ_DLX do arquivo .env.

## 📁 Estrutura do Projeto

- application/: Lógica de aplicação
- domain/: Entidades e interfaces do domínio
- framework/: Integrações com frameworks e bibliotecas externas
- Dockerfile e docker-compose.yaml: Configurações para containerização
- Makefile: Comandos para automação de tarefas
- main.go: Ponto de entrada da aplicação
