# CLI Convert – Conversor Universal de Arquivos

**CLI Convert** é uma ferramenta de linha de comando (CLI) poderosa e flexível, construída com **Go (Golang)**, projetada para converter arquivos de dados entre vários formatos. Suporta transformações entre **JSON, CSV, XML e YAML**, com foco em performance e facilidade de uso.

Ideal para desenvolvedores, engenheiros de dados e qualquer pessoa que trabalhe frequentemente com diferentes formatos de dados.

---

## ✨ Funcionalidades

* **Conversão Bidirecional:** Converta entre JSON, CSV, XML e YAML.
  * `JSON <-> CSV`
  * `JSON <-> XML`
  * `JSON <-> YAML`
  * `CSV <-> XML`
  * `CSV <-> YAML`
  * `XML <-> YAML`
* **Auto-detecção de Formato:** Detecta automaticamente o formato de entrada (não precisa de `--from`).
* **Validação Robusta:** Garante que arquivos de entrada existem, não estão vazios e seguem o formato especificado.
* **Tratamento Inteligente de Tipos:** Detecta e converte automaticamente valores numéricos e booleanos.
* **Preservação de Ordem XML:** Mantém a ordem dos elementos ao converter XML para JSON.
* **🤖 Integração com IA (Opcional):** Gere schemas, pergunte sobre dados e detecte formatos usando IA.

---

## 🚀 Instalação

Certifique-se de ter o Go instalado. Então:

```bash
# Clone o repositório
git clone https://github.com/alvarossantos/cli_convert.git
cd cli_build

# Compile o executável
go build

# (Opcional) Mova para seu PATH
mv cli-convert /usr/local/bin/
```

---

## 📖 Uso

### Comando Principal: `convert`

```bash
cli-convert convert --from <formato> --to <formato> --input <arquivo> --output <arquivo> [opções]
```

### Flags

| Flag | Obrigatório | Descrição |
|------|:-----------:|-----------|
| `--input` | ✅ | Caminho do arquivo de entrada |
| `--output` | ✅ | Caminho do arquivo de saída |
| `--from` | ❌ | Formato de origem (detectado automaticamente se omitido) |
| `--to` | ✅ | Formato de destino |
| `--delimiter` | ❌ | Delimitador CSV (padrão: `,`) |
| `--root` | ❌ | Nome do elemento raiz para XML (padrão: `root`) |

### Exemplos de Conversão

```bash
# JSON para CSV
cli-convert convert --from json --to csv --input data.json --output data.csv

# CSV para JSON (com delimitador ponto-e-vírgula)
cli-convert convert --from csv --to json --input dados.csv --output dados.json --delimiter ';'

# Auto-detectar formato e converter para JSON
cli-convert convert --to json --input dados.csv --output dados.json

# XML para YAML
cli-convert convert --from xml --to yaml --input config.xml --output config.yaml

# YAML para XML (com elemento raiz customizado)
cli-convert convert --from yaml --to xml --input dados.yaml --output dados.xml --root MeusDados
```

---

## 🤖 Comandos de IA

Os comandos de IA requerem configuração no arquivo `.env`. Copie `.env.example` para `.env` e configure sua chave de API.

### `detect` — Auto-detectar Formato

```bash
cli-convert detect --input arquivo.json
# Saída: Detected format: json
```

### `schema` — Gerar JSON Schema

Para arquivos JSON, o schema é gerado localmente. Para outros formatos, usa IA.

```bash
cli-convert schema --input dados.csv
```

### `ask` — Perguntar sobre os Dados

Faça perguntas em linguagem natural sobre o conteúdo de um arquivo:

```bash
cli-convert ask --input vendas.csv --question "Qual o total de vendas?"
cli-convert ask --input usuarios.json --question "Quantos usuários são maiores de 18 anos?"
```

### Configuração de IA

Copie `.env.example` para `.env`:

```bash
cp .env.example .env
```

Configure o provedor e a chave de API:

| Provedor | Variável | Modelo Padrão |
|----------|----------|---------------|
| OpenRouter | `OPENROUTER_API_KEY` | `inclusionai/ling-3.0-flash:free` |
| OpenAI | `OPENAI_API_KEY` | `gpt-4o-mini` |
| Gemini | `GEMINI_API_KEY` | `gemini-2.0-flash` |

---

## ⚠️ Tratamento de Erros

A ferramenta fornece mensagens de erro informativas para:

* Flags obrigatórias ausentes (`--input`, `--output`, `--to`)
* Arquivos de entrada inválidos ou inexistentes
* Formatos de conversão não suportados
* Dados de entrada malformados
* Delimitador com mais de um caractere

---

## 🤝 Contribuições

Contribuições são bem-vindas! Sinta-se à livre para abrir issues ou enviar pull requests.

---

## 📄 Licença

Este projeto está licenciado sob a MIT License.
