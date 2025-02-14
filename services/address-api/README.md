# API DE ENDEREÇOS

## Descrição

API responsável pelo gerenciamento de endereços, UBS (Unidades Básicas de Saúde) e equipes de saúde. Esta API é fundamental para o mapeamento territorial da área de atendimento, permitindo vincular segmentos de ruas específicos a equipes de saúde e suas respectivas UBS.

## Como Executar apenas essa API

1. Navegue até o diretório da Address API
2. Verifique se o arquivo `docker-compose.yaml` está presente
3. Execute o Docker Compose:

```bash
docker-compose up -d
```

**A API estará disponível na porta 8083.**

Isso irá iniciar:

- A API de endereços na porta 8083
- Um banco PostgreSQL na porta 5432
- Uma instância do Adminer na porta 8084 para gerenciar o banco de dados

## Endpoints

### UBS (Unidades Básicas de Saúde)

A seguir os endpoints relacionados ao gerenciamento de UBS no sistema.

#### Criar UBS

**POST** `/ubs/`

Cria uma nova Unidade Básica de Saúde no sistema.

Corpo da requisição:

```json
{
    "name": "UBS Nome",               // Obrigatório, máximo 200 caracteres
    "address": "Endereço Completo",   // Obrigatório, máximo 200 caracteres
    "city": "Cidade",                 // Obrigatório, máximo 100 caracteres
    "state": "SP",                    // Obrigatório, exatamente 2 caracteres
    "cep": "12345678"                 // Obrigatório, exatamente 8 dígitos
}
```

Resposta de sucesso (201 Created):

```json
{
    "success": true,
    "data": {
        "id": 1,
        "name": "UBS Nome",
        "address": "Endereço Completo",
        "city": "Cidade",
        "state": "SP",
        "cep": "12345678",
        "created_at": "2024-02-13T10:00:00Z",
        "updated_at": "2024-02-13T10:00:00Z"
    }
}
```

#### Listar todas UBS

**GET** `/ubs/`

Retorna uma lista de todas as UBS cadastradas no sistema.

Resposta de sucesso (200 OK):

```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "name": "UBS Nome",
            "address": "Endereço Completo",
            "city": "Cidade",
            "state": "SP",
            "cep": "12345678",
            "teams": [
                // Lista de equipes vinculadas (se houver)
            ]
        }
        // ... outras UBS ...
    ]
}
```

#### Buscar UBS específica

**GET** `/ubs/{id}`

Retorna os dados de uma UBS específica, incluindo suas equipes.

Parâmetros de URL:

- `id`: ID numérico da UBS (obrigatório)

Resposta de sucesso (200 OK):

```json
{
    "success": true,
    "data": {
        "id": 1,
        "name": "UBS Nome",
        "address": "Endereço Completo",
        "city": "Cidade",
        "state": "SP",
        "cep": "12345678",
        "teams": [
            {
                "id": 1,
                "name": "Nome da Equipe",
                "ubs_id": 1
            }
        ]
    }
}
```

#### Atualizar UBS

**PUT** `/ubs/{id}`

Atualiza os dados de uma UBS existente.

Parâmetros de URL:

- `id`: ID numérico da UBS (obrigatório)

Corpo da requisição:

```json
{
    "name": "Novo Nome da UBS",
    "address": "Novo Endereço",
    "city": "Nova Cidade",
    "state": "RJ",
    "cep": "87654321"
}
```

#### Deletar UBS

**DELETE** `/ubs/{id}`

Remove uma UBS do sistema. Só é possível deletar uma UBS que não possui equipes vinculadas.

Parâmetros de URL:

- `id`: ID numérico da UBS (obrigatório)

### Equipes

A seguir os endpoints relacionados ao gerenciamento de equipes de saúde.

#### Criar Equipe

**POST** `/teams/`

Cria uma nova equipe de saúde vinculada a uma UBS.

Corpo da requisição:

```json
{
    "name": "Nome da Equipe",    // Obrigatório, máximo 100 caracteres
    "ubs_id": 1                  // Obrigatório, ID da UBS existente
}
```

#### Listar Equipes

**GET** `/teams/`

Retorna todas as equipes cadastradas com suas respectivas UBS.

#### Buscar Equipe

**GET** `/teams/{id}`

Retorna os dados de uma equipe específica.

#### Atualizar Equipe

**PUT** `/teams/{id}`

Atualiza os dados de uma equipe existente.

#### Deletar Equipe

**DELETE** `/teams/{id}`

Remove uma equipe. Só é possível deletar equipes que não possuem segmentos de rua vinculados.

### Segmentos de Rua

Endpoints utilizados para gerenciar os segmentos de ruas e suas associações com equipes de saúde.

#### Criar Segmento

**POST** `/streets/`

Cria um novo segmento de rua e o vincula a uma equipe.

Corpo da requisição:

```json
{
    "street_name": "Nome da Rua",     // Obrigatório, máximo 200 caracteres
    "street_type": "Rua",             // Obrigatório (Rua, Avenida, etc)
    "neighborhood": "Bairro",         // Obrigatório, máximo 100 caracteres
    "city": "Cidade",                 // Obrigatório, máximo 100 caracteres
    "state": "SP",                    // Obrigatório, exatamente 2 caracteres
    "start_number": 1,                // Número inicial do segmento
    "end_number": 100,               // Número final do segmento
    "cep_prefix": "12345",           // 5 primeiros dígitos do CEP
    "even_odd": "all",               // "even", "odd" ou "all"
    "team_id": 1                     // ID da equipe responsável
}
```

#### Buscar Equipe por Endereço

**GET** `/streets/search`

Busca a equipe responsável por um endereço específico.

Parâmetros de query:

- `street`: Nome da rua
- `number`: Número
- `city`: Cidade
- `state`: Estado

Exemplo:

```
GET /streets/search?street=Rua%20Exemplo&number=123&city=São%20Paulo&state=SP
```

Resposta de sucesso (200 OK):

```json
{
    "success": true,
    "data": {
        "street_segment": {
            // dados do segmento encontrado
        },
        "team": {
            // dados da equipe responsável
        },
        "ubs": {
            // dados da UBS
        }
    }
}
```

### Códigos de Erro

A API pode retornar os seguintes códigos de erro:

- 400 Bad Request
  - Payload inválido
  - ID inválido
  - Dados obrigatórios faltando
  - Violação de regras de negócio
- 404 Not Found
  - Recurso não encontrado
- 409 Conflict
  - Conflito com recursos existentes
- 500 Internal Server Error
  - Erro interno do servidor

### Observações Importantes

1. O sistema utiliza a extensão pg_trgm do PostgreSQL para busca fuzzy de endereços
2. Todos os nomes de ruas são normalizados (removendo acentos e padronizando maiúsculas/minúsculas)
3. Os segmentos de rua podem ser configurados para números pares, ímpares ou ambos
4. O sistema valida a sobreposição de segmentos de rua para evitar conflitos de territórioParâmetros: street, number, city, state
