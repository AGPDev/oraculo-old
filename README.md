# Oráculo

Projeto responsável por sincronizar os dumps e repositórios da nuvem para o servidor local, reduzindo o consumo de internet e acelerando o desenvolvimento.

## Desenvolvimento

#### Requisitos

[Git](https://git-scm.com/)
[Go](https://golang.org/)

#### Passos

Oráculo sync

`cp exampleconfig.toml config.toml`

`go mod tidy && go mod download`

`go run oraculo.go`

Oráculo Http

`./gohttpserver --conf confighttp.yml`

Compilação

`go build oraculo.go`

## Utilização

De permissão de execução para os programas

`chmod +x gohttpserver oraculo`

Configure o programa

`cp exampleconfig.toml config.toml`

Execute os programas

`./gohttpserver`

`./oraculo`

Configurações personalizadas

`config.toml`  - Configurações do Oráculo

`confighttp.yml` - Configurações do servidor Http

## Changelog

### [2.0.2] - 2019-02-20

- Fix evitar download de dumps do dia anterior

### [2.0.1] - 2019-02-19

- Ignorando config e adicionando config de exemplo
- Fix jobs faltava configurar a data e hora atual
- Adicionada sincronização para criar a pasta de downloads
- Fix limpar pasta do dia
- Fix demora para cancelar download lento
- Add bloqueio por processos que devem executar antes de qualquer outro processo
- Fix gerar novos dumps somente uma vez por dia

### [2.0.0] - 2019-02-19

Refatoração da lógica de trabalho para melhorar o desempenho e acelerar a sincronização, agora todos os dados de API são salvos localmente e atualizados de forma asincrona de acordo com a configuração inserida.

- Agendamento dos serviços de sincronização baseado no Cron.
- Organização do projeto para melhor manutenção.
- Armazenamento de informações em json para atualizar somente o for e quando for necessário.
- Algoritmos para vincular as empresas do Movidesk com os Projetos do Gitlab (pode ser melhorado pois algumas empresas ainda não estão vinculadas).
- Serviços separados cada serviço poderá ser executado independentemente em seu horário agendado.
- Data e hora de execução dos processos na tela do programa

### [1.0.0] - 2019-02-13

Funcionalidades básicas que constituem um programa espelho com inteligencia para tentar se antecipar as  necessidades dos desenvolvedores.

- Conexão com a API do Movidesk e listagem das empresas de todos os chamados em "Aberto" ou "Em atendimento"
- Conexão com a API do Gitlab e listagem dos projetos que contenham o nome da empresa no repositório.
- Listagem de jobs execução de novos jobs e retorno do resultado dos jobs.
- Download do arquivo de dump de acordo com o restuldado do job.
- Armazenar resultado em json para gerenciar os jobs e dumps já executados.
- Servidor HTTP de arquivos gerenciavel.