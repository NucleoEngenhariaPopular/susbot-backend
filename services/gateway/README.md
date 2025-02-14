# GATEWAY API

## Descrição

Gateway principal do sistema, responsável por integrar todos os serviços e gerenciar a comunicação entre o Botkit (motor de conversação) e o Twilio (serviço de mensageria). O Gateway atua como ponto central de entrada para mensagens dos usuários, coordenando o fluxo de dados entre os diferentes componentes do sistema.

## Como Executar o Sistema Completo

1. Navegue até o diretório raiz do projeto
2. Crie um arquivo `.env` com as variáveis necessárias
3. Execute o Docker Compose:

```bash
docker-compose up -d
```

**O Gateway estará disponível na porta 8080.**

Isso irá iniciar todos os serviços necessários:

- Gateway (8080)
- User API (8081)
- Conversation API (8082)
- Address API (8083)
- Adminer (8084)
- Mongo Express (8085)
- Botkit/Fluxo (3000)
- MongoDB
- PostgreSQL
- Ngrok

## Variáveis de Ambiente

O arquivo `.env` deve conter:

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=postgres
MONGO_URI=mongodb://root:example@mongo:27017/
MONGODB_NAME=my_database
MONGODB_COLLECTION=conversations
BOTKIT_URL=http://fluxo:3000/api/messages
TWILIO_SID=seu_sid_aqui
```

## Endpoints

### Webhook do Twilio

**POST** `/`

Endpoint principal que recebe as mensagens do Twilio.

Corpo da requisição (form-data):

```json
{
    "MessageSid": "identificador_único",
    "AccountSid": "conta_twilio",
    "From": "número_origem",
    "To":
