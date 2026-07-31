package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Provider representa um provedor de IA configurado.
type Provider struct {
	Name         string
	BaseURL      string
	EnvKey       string
	DefaultModel string
}

var providers = map[string]Provider{
	"openai": {
		Name:         "OpenAI",
		BaseURL:      "https://api.openai.com/v1",
		EnvKey:       "OPENAI_API_KEY",
		DefaultModel: "gpt-4o-mini",
	},
	"gemini": {
		Name:         "Google Gemini",
		BaseURL:      "https://generativelanguage.googleapis.com/v1beta/openai/",
		EnvKey:       "GEMINI_API_KEY",
		DefaultModel: "gemini-2.0-flash",
	},
	"openrouter": {
		Name:         "OpenRouter",
		BaseURL:      "https://openrouter.ai/api/v1",
		EnvKey:       "OPENROUTER_API_KEY",
		DefaultModel: "inclusionai/ling-3.0-flash:free",
	},
}

// GetProvider retorna o provedor ativo com base na variável AI_PROVIDER.
func GetProvider() Provider {
	name := strings.ToLower(os.Getenv("AI_PROVIDER"))
	if name == "" {
		name = "openrouter"
	}
	p, ok := providers[name]
	if !ok {
		p = providers["openrouter"]
	}
	return p
}

// GetModel retorna o modelo ativo (AI_MODEL override > padrão do provedor).
func GetModel() string {
	if m := os.Getenv("AI_MODEL"); m != "" {
		return m
	}
	return GetProvider().DefaultModel
}

// GetApiKey retorna a chave de API do provedor ativo.
func GetApiKey() string {
	return os.Getenv(GetProvider().EnvKey)
}

// chatMessage representa uma mensagem na conversa com a IA.
type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatRequest representa o corpo da requisição à API de IA.
type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
}

// chatChoice representa uma escolha na resposta da IA.
type chatChoice struct {
	Message chatMessage `json:"message"`
}

// chatResponse representa a resposta da API de IA.
type chatResponse struct {
	Choices []chatChoice `json:"choices"`
}

// CallAI faz uma chamada à API de IA e retorna a resposta como string.
func CallAI(messages []chatMessage) (string, error) {
	provider := GetProvider()
	apiKey := GetApiKey()
	model := GetModel()

	if apiKey == "" {
		return "", fmt.Errorf("chave da API de IA não configurada. Defina %s no .env", provider.EnvKey)
	}

	body := chatRequest{
		Model:       model,
		Messages:    messages,
		MaxTokens:   500,
		Temperature: 0.7,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar request: %w", err)
	}

	url := provider.BaseURL + "/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("erro ao criar request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Headers específicos do OpenRouter
	if provider.Name == "OpenRouter" {
		req.Header.Set("HTTP-Referer", "https://github.com/alvarossantos/cli_convert")
		req.Header.Set("X-Title", "cli-convert")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro na chamada HTTP: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("erro da API de IA (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("erro ao parsear resposta: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("resposta vazia da API de IA")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// ──────────────────────────────────────────────
//  Features de IA
// ──────────────────────────────────────────────

// DetectFormat detecta o formato de um arquivo usando assinatura de conteúdo.
// Retorna o formato (json, csv, xml, yaml) sem precisar de IA.
func DetectFormat(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	trimmed := strings.TrimSpace(string(data))

	if len(trimmed) == 0 {
		return "", fmt.Errorf("arquivo vazio")
	}

	// JSON: começa com { ou [
	if trimmed[0] == '{' || trimmed[0] == '[' {
		if json.Valid(data) {
			return "json", nil
		}
	}

	// XML: começa com < (possivelmente após <?xml ...)
	if trimmed[0] == '<' {
		return "xml", nil
	}

	// YAML: verifica padrões típicos (key: value, - item, etc.)
	lines := strings.Split(trimmed, "\n")
	yamlScore := 0
	csvScore := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Padrões YAML
		if strings.Contains(line, ": ") || strings.HasPrefix(line, "- ") {
			yamlScore++
		}

		// Padrões CSV: múltiplas vírgulas por linha
		commas := strings.Count(line, ",")
		if commas >= 1 {
			csvScore++
		}
	}

	// CSV: se todas as linhas têm o mesmo número de vírgulas
	if csvScore > 0 && csvScore == yamlScore {
		// Mais provável ser CSV se há header consistente
		commasPerLine := make([]int, 0)
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			commasPerLine = append(commasPerLine, strings.Count(line, ","))
		}
		if len(commasPerLine) > 1 {
			allSame := true
			for _, c := range commasPerLine[1:] {
				if c != commasPerLine[0] {
					allSame = false
					break
				}
			}
			if allSame && commasPerLine[0] > 0 {
				return "csv", nil
			}
		}
	}

	if yamlScore > csvScore {
		return "yaml", nil
	}
	if csvScore > 0 {
		return "csv", nil
	}

	return "", fmt.Errorf("não foi possível detectar o formato automaticamente")
}

// InferSchema gera um JSON Schema a partir dos dados de um arquivo.
func InferSchema(filePath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	schema := map[string]interface{}{
		"title":    "Schema gerado automaticamente",
		"source":   filePath,
		"type":     "object",
		"properties": map[string]interface{}{},
	}

	trimmed := strings.TrimSpace(string(data))

	// JSON
	if trimmed[0] == '{' || trimmed[0] == '[' {
		var parsed interface{}
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, fmt.Errorf("JSON inválido: %w", err)
		}
		schema["properties"] = generateJSONSchemaProperties(parsed)
		schema["format"] = "json"
		return schema, nil
	}

	// Para outros formatos, delega à IA
	format := "desconhecido"
	if trimmed[0] == '<' {
		format = "xml"
	} else {
		// Tenta detectar
		format, _ = DetectFormat(filePath)
	}

	// Amostra dos primeiros 500 chars
	sample := trimmed
	if len(sample) > 500 {
		sample = sample[:500] + "..."
	}

	prompt := fmt.Sprintf(`Analise este arquivo %s e gere um JSON Schema descrevendo a estrutura dos dados.
Retorne APENAS o JSON válido, sem markdown, sem explicação.
Arquivo:
%s`, format, sample)

	messages := []chatMessage{
		{Role: "system", Content: "Você é um gerador de schemas. Retorne APENAS JSON válido."},
		{Role: "user", Content: prompt},
	}

	response, err := CallAI(messages)
	if err != nil {
		return nil, err
	}

	// Tenta parsear a resposta como JSON
	response = strings.TrimSpace(response)
	// Remove possíveis wrappers markdown
	if strings.HasPrefix(response, "```") {
		lines := strings.Split(response, "\n")
		var jsonLines []string
		inBlock := false
		for _, line := range lines {
			if strings.HasPrefix(line, "```") {
				inBlock = !inBlock
				continue
			}
			if inBlock {
				jsonLines = append(jsonLines, line)
			}
		}
		response = strings.Join(jsonLines, "\n")
	}

	var aiSchema map[string]interface{}
	if err := json.Unmarshal([]byte(response), &aiSchema); err != nil {
		// Se não conseguiu parsear, retorna o schema básico com a resposta da IA
		schema["ai_analysis"] = response
		schema["format"] = format
		return schema, nil
	}

	aiSchema["source"] = filePath
	aiSchema["format"] = format
	return aiSchema, nil
}

// AskQuestion permite fazer perguntas em linguagem natural sobre um arquivo.
func AskQuestion(filePath string, question string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	trimmed := strings.TrimSpace(string(data))

	// Detecta formato
	format, _ := DetectFormat(filePath)

	// Amostra dos dados (máx 1000 chars)
	sample := trimmed
	if len(sample) > 1000 {
		sample = sample[:1000] + "\n... (arquivo truncado)"
	}

	systemPrompt := `Você é um assistente especializado em análise de dados.
O usuário tem um arquivo de dados e quer entender melhor seu conteúdo.
Responda de forma objetiva e direta em português brasileiro.
Se o arquivo contiver dados tabulares, use tabelas quando apropriado.`

	userPrompt := fmt.Sprintf(`Tenho um arquivo %s com os seguintes dados:

%s

Minha pergunta é: %s`, format, sample, question)

	messages := []chatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	return CallAI(messages)
}

// generateJSONSchemaProperties gera propriedades de schema recursivamente.
func generateJSONSchemaProperties(data interface{}) map[string]interface{} {
	props := make(map[string]interface{})

	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			prop := map[string]interface{}{}
			switch val := value.(type) {
			case nil:
				prop["type"] = "null"
			case bool:
				prop["type"] = "boolean"
			case float64:
				prop["type"] = "number"
			case string:
				prop["type"] = "string"
			case []interface{}:
				prop["type"] = "array"
				if len(val) > 0 {
					prop["items"] = generateJSONSchemaProperties(val[0])
				}
			case map[string]interface{}:
				prop["type"] = "object"
				prop["properties"] = generateJSONSchemaProperties(val)
			default:
				prop["type"] = "string"
			}
			props[key] = prop
		}
	case []interface{}:
		props["items"] = map[string]interface{}{"type": "object"}
	}

	return props
}
