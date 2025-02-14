# API DE CONVERSAS

## Descrição

API responsável pelo gerenciamento de conversas e mensagens entre usuários e o sistema. Esta API utiliza MongoDB para armazenar o histórico completo das conversas, permitindo rastrear todas as interações e manter o contexto das comunicações.

## Como Executar apenas essa API

1. Navegue até o diretório da Conversation API
2. Verifique se o arquivo `docker-compose.yaml` está presente
3. Execute o Docker Compose:

```bash
docker-compose up -d
```

**A API estará disponível na porta 8082.**

Isso irá iniciar:

- A API de conversas na porta 8082
- Um servidor MongoDB na porta padrão 27017
- Uma interface Mongo Express na porta 8085 para gerenciar o banco de dados

## Endpoints

### Conversas

A seguir os endpoints relacionados ao gerenciamento de conversas no sistema.

#### Salvar Mensagem

**POST** `/conversations/`

Salva uma nova mensagem. Se não existir uma conversa ativa para o usuário, uma nova conversa será criada automaticamente.

Corpo da requisição:

```json
{
    "user_id": "id_do_usuario",                // Obrigatório, identificador do usuário
    "sender": "origem_mensagem",               // Obrigatório, quem enviou (usuário ou sistema)
    "text": "conteúdo da mensagem",           // Obrigatório, conteúdo da mensagem
    "timestamp": "2024-02-13T10:00:00Z"       // Obrigatório, momento do envio
}
```

Resposta de sucesso (201 Created):

```json
{
    "success": true,
    "data": {
        "id": "conversation_id",
        "user_id": "id_do_usuario",
        "start_time": "2024-02-13T10:00:00Z",
        "messages": [
            {
                "user_id": "id_do_usuario",
                "sender": "origem_mensagem",
                "text": "conteúdo da mensagem",
                "timestamp": "2024-02-13T10:00:00Z"
            }
        ]
    }
}
```

#### Buscar Conversa

**GET** `/conversations/{id}`

Retorna uma conversa específica com todas as suas mensagens.

Parâmetros de URL:

- `id`: ID da conversa (MongoDB ObjectId)

Resposta de sucesso (200 OK):

```json
{
    "success": true,
    "data": {
        "id": "conversation_id",
        "user_id": "id_do_usuario",
        "start_time": "2024-02-13T10:00:00Z",
        "end_time": null,
        "messages": [
            {
                "user_id": "id_do_usuario",
                "sender": "origem_mensagem",
                "text": "conteúdo da mensagem",
                "timestamp": "2024-02-13T10:00:00Z"
            }
            // ... outras mensagens ...
        ]
    }
}
```

### Códigos de Erro

A API pode retornar os seguintes códigos de erro:

- 400 Bad Request
  - Payload inválido
  - ID de conversa inválido
  - Dados obrigatórios faltando
- 404 Not Found
  - Conversa não encontrada
- 500 Internal Server Error
  - Erro interno do servidor

### Observações Importantes

1. O sistema mantém conversas "ativas" e "inativas"
   - Conversas ativas não possuem `end_time`
   - Apenas uma conversa por usuário pode estar ativa por vez
2. As mensagens são ordenadas cronologicamente pelo campo `timestamp`
3. O MongoDB garante a persistência e escalabilidade do histórico de conversas
4. Use o Mongo Express (porta 8085) para:
   - Visualizar conversas e mensagens
   - Monitorar o uso do banco de dados
   - Realizar queries ad-hoc quando necessário
