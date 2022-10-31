# Oráculo

Aqui você encontra um espelho dos arquivos e repositórios utilizados para acelerar o desenvolvimento.

Esses estão armazenados em uma maquina na rede, acelerando o download e consequentemente o seu desenvolvimento.

## Fucionamento

**Pasta "dumps":**

​	Resumindo, um programa fica verificando todos os chamados em aberto, inclusive os que ainda estão somente com o atendimento, lista todos os projetos no Gitlab com o nome da empresa que abriu o chamado, em seguida verifica se o dump do banco ainda não foi gerado no dia, caso ainda não tenha sido gerado, executa a tarefa de dump e faz o download do mesmo.

​	Caso algum colaborador já tenha executado a tarefa e o arquivo não está aqui ou não tem tarefa na proxima verificação ele será baixado.

​	As sincronizações estão configuradas para executarem:

| Sincronização                           | Quando                          |
| :-------------------------------------- | :------------------------------ |
| Empresas do Movidesk                    | Seg-Sex às 6h                   |
| Projetos do Gitlab                      | Seg-Sex às 7,12,18h             |
| Vinculação das empresas com os projetos | Seg-Sex às 7,12,18h 15m         |
| Jobs dos projetos no Gitlab             | Seg-Sex das 7 às 19h a cada 10m |
| Dumps gerados                           | Seg-Sex das 6 às 18h a cada 5m  |
| Tickets do Movidesk                     | Seg-Sex das 6 às 18h a cada 19m |
| Download dos dumps                      | Seg-Sex das 7 às 19h a cada 5m  |
| Limpeza da pasta today                  | Seg-Sex às 4h                   |
| Criação da pasta de downloads           | Seg-Sex às 0h                   |

**Pasta "dumps/today":**

​	Todos os bancos gerados no dia serão copiados para esta pasta com um nome amigavel para facilitar o uso de **"scripts"** de instalação de lojas

## TODO

- [ ] Espelho dos repositórios do Gitlab

  - [ ] Instruções de como e quando usar
- [x] Limpar pasta **"today"** automaticamente
- [x] README nos respositórios
- [x] Corrigir loop de falha de download
- [x] Adicionar data e hora nos logs em tela no sincronizador
- [x] Adicionar um cache permanente das lojas e dos projetos para diminuir a consulta na API
- [ ] Exibir debug dos processos na tela do navegador
- [ ] Opção para solicitar a geração de um dump
- [ ] Pasta scripts com scripts uteis onde as pessoas poderão fazer contribuições com exemplo de como contribuir
- [ ] Apagar dumps antigos

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