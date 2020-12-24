# apiGo
API desenvolvida em GoLang

# Condições para executar a importação do CSV

Ter o arquivo para importação em:
 - Arquivo/q1_catalog.csv
 
Rodar o comando
 - go build
 - ./apiGo

# Adicionando SQLite

Rodar o comando
 - go get github.com/mattn/go-sqlite3

# Adicionando controle de rodas

Rodar o comando
 - go get github.com/gorilla/mux
 
# Utilizando a API

GET http://localhost:8000/companies
 - Params:
  - name
  - zip

POST http://localhost:8000/companies
 - Body (Json)
  - {"name":"MÁRIO MOTOS","zip":"54896","website":"motomario.com.br"}
