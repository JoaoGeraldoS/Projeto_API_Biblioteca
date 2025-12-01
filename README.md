# Projeto_API_Biblioteca

API REST para gerenciamento de livros, autores, categorias e usuarios.
Fornece endpoints para cadastro, leitura, atualização, remoção e relacionamento entre entidades.

---

## Tecnologias

- Go 
- Gin Framework
- Mysql
- SQLite (desenvolvimento e tests)
- Swagger (documentação)
- Testes unitários com mocks

---

## Estrutura do projeto

- cmd/
    main.go

internal/
  authors/
  books/
  categories/
  users/

  database/
  routes/
  middleware/
  logger/
  config/

docs/

---

## Como rodar localmente

1. Clone o repositório
    ``` bash
    git clone https://github.com/JoaoGeraldoS/Projeto_API_Biblioteca

2. Instale as dependências
    ~~~~
        go mod tidy
    ~~~~
3. Configure as variáveis de ambiente
    ``` DATABASE_URL: "user:pass@tcp(localhost:3306)/library"
    SECRET_KEY: "secret key"
    LOGGER_APP: "development" || "production"
    