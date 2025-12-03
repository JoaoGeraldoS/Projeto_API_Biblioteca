# Projeto_API_Biblioteca

API REST para gerenciamento de livros, autores, categorias e usuários.  
Fornece endpoints para cadastro, leitura, atualização, remoção e relacionamentos entre entidades.


---

## Tecnologias

- Go 
- Gin Framework
- MySQL 
- SQLite (desenvolvimento e tests)
- Swagger (documentação)
- Zap Logger  
- Testes unitários com mocks

---

## Estrutura do projeto

- cmd/
    - main.go

- internal/
    - authors/
    - books/
    - categories/
    - users/

    - database/
    - routes/
    - middleware/
    - logger/
- docs/
- migrations/

---

## Como rodar localmente

### 1. Clone o repositório
~~~    
git clone https://github.com/JoaoGeraldoS/Projeto_API_Biblioteca
cd Projeto_API_Biblioteca
~~~

### 2. Instale as dependências
``` bash
go mod tidy
```

### 3. Configure as variáveis de ambiente
Crie um arquivo 
``` bash
DATABASE_URL: "user:pass@tcp(localhost:3306)/library"
SECRET_KEY: "secret key"
LOGGER_APP: "development" # production
```

### 4. Subir o banco de dados (MySQL via Docker)
``` bash
docker compose up
```

### 5. Rode as migarções
``` bash
migrate -path . -database "mysql://user:pass@tcp(localhost:3306)/library" up
```

### 6. Rode a aplicação
``` bash
go run cmd/main.go
```
### 7. Testes
Test coverage: atualmente apenas o módulo de livros possui testes.
Os módulos restantes seguem o mesmo padrão e terão cobertura adicionada nas próximas versões.
``` bash
go test ./...
```
---

## Endpoints principais
A documentação completa está disponível no Swagger:
```bash 
/swagger/index.html
```
Exemplos básicos
```bash
POST /api/books
GET /public/api/books/:id
POST /api/books/relation
```
---
## Padronização de erros

| Código | Significado |
|-----|-----|
| 400 | BadRequest |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | NotFound |
| 500 | InternalServerError |

---
## Exempls de payload

Criar livro
``` json 

{
    "title": "A menina e o porquinho",
    "description": "Esse livro é sobre conteudo infantil",
    "content": "A menina é o porquinho",
    "author_id": 1
}
```

Ler livro (response)
``` json

{
    "author": {
        "description": "string",
        "id": 0,
        "name": "string"
    },
    "author_id": 0,
    "categories": [
        {
        "createdAT": "string",
        "id": 0,
        "name": "string"
        }
    ],
    "content": "string",
    "created_at": "string",
    "description": "string",
    "id": 0,
    "title": "string",
    "updated_at": "string"
}
```
---
## Middleware de Logs

- Method
- Path
- Status
- client_id
- Latency
- User_agent

---

## RoadMap

- Melhorar corbetura de tests
- Implementar Dockerfile
- Melhorar documentação
