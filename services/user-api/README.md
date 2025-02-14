# Descrição

API responsável pelo gerenciamento completo de usuários no sistema. Esta API permite realizar todas as operações CRUD (Criar, Ler, Atualizar e Deletar) relacionadas aos usuários e suas informações pessoais. Além disso, ela se integra com a Address API para vincular cada usuário a uma equipe de saúde com base em seu endereço.

## Como Executar apenas essa API

1. Navegue até o diretório da User API
2. Verifique se o arquivo `docker-compose.yaml` está presente
3. Execute o Docker Compose:

```bash
docker-compose up -d
```

**A API estará disponível na porta 8081.**

Isso irá iniciar:

- A API de usuários na porta 8081
- Um banco PostgreSQL na porta 5432
- Uma instância do Adminer na porta 8084 para gerenciar o banco de dados

## Endpoints

### Usuários

A seguir os endpoints disponíveis para gerenciamento de usuários no sistema.

#### Criar Usuário

POST /users/

Cria um novo usuário no sistema. Todos os campos marcados como obrigatórios devem ser fornecidos.

Corpo da requisição:

```json
{
"name": "Nome do Usuário",           // Obrigatório, máximo 200 caracteres
"cpf": "12345678900",                // Obrigatório, exatamente 11 dígitos
"date_of_birth": "1990-01-01",       // Obrigatório, formato YYYY-MM-DD
"phone_number": "11999999999",       // Opcional, máximo 20 caracteres
"street_name": "Nome da Rua",        // Obrigatório, máximo 200 caracteres
"street_number": "123",              // Obrigatório, máximo 20 caracteres
"complement": "Apto 1",              // Opcional, máximo 100 caracteres
"neighborhood": "Bairro",            // Obrigatório, máximo 100 caracteres
"city": "Cidade",                    // Obrigatório, máximo 100 caracteres
"state": "SP",                       // Obrigatório, exatamente 2 caracteres
"cep": "12345678"                    // Obrigatório, exatamente 8 dígitos
}
```

Resposta de sucesso (201 Created):

```json
{
"success": true,
"data": {
"id": 1,
// ... todos os dados do usuário ...
"created_at": "2024-02-13T10:00:00Z",
"updated_at": "2024-02-13T10:00:00Z"
}
}
```

#### Buscar Usuário por ID

GET /users/{id}

Retorna os dados de um usuário específico, incluindo (se houver) informações da equipe de saúde responsável pela sua região.

Parâmetros de URL:

id: ID numérico do usuário (obrigatório)

Resposta de sucesso (200 OK):

```json
{
"success": true,
"data": {
"user": {
// dados do usuário
},
"team": {
"id": 1,
"name": "Nome da Equipe",
"ubs_name": "Nome da UBS"
}
}
}
```

#### Buscar Usuário por CPF

GET /users/cpf/{cpf}

Busca um usuário pelo seu CPF, retornando também as informações da equipe de saúde.

Parâmetros de URL:

cpf: CPF do usuário com 11 dígitos (obrigatório)

Resposta de sucesso (200 OK):

```json
{
"success": true,
"data": {
"user": {
// dados do usuário
},
"team": {
"id": 1,
"name": "Nome da Equipe",
"ubs_name": "Nome da UBS"
}
}
}
```

#### Atualizar Usuário

PUT /users/{id}

Atualiza os dados de um usuário existente. Apenas os campos enviados serão atualizados.

Parâmetros de URL:

id: ID numérico do usuário (obrigatório)
Corpo da requisição:

```json
{
"name": "Novo Nome",                // Opcional
"phone_number": "11999999999",       // Opcional
"street_name": "Nova Rua",           // Opcional
"street_number": "456",              // Opcional
"complement": "Casa",                // Opcional
"neighborhood": "Novo Bairro",       // Opcional
"city": "Nova Cidade",               // Opcional
"state": "RJ",                       // Opcional
"cep": "87654321"                    // Opcional
}
```

Resposta de sucesso (200 OK):

```json
{
  "success": true,
  "data": {
    // dados atualizados do usuário
  }
}
```

#### Deletar Usuário

DELETE /users/{id}

Remove um usuário do sistema.

Parâmetros de URL:

id: ID numérico do usuário (obrigatório)
Resposta de sucesso (200 OK):

```json
{
  "success": true,
  "data": "User successfully deleted"
}
```

## Códigos de Erro

A API pode retornar os seguintes códigos de erro:

400 Bad Request

- Payload inválido
- ID inválido
- Dados obrigatórios faltando

404 Not Found

- Usuário não encontrado

409 Conflict

- CPF já cadastrado

500 Internal Server Error

- Erro interno do servidor
