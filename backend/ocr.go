package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	defaultQwenBaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
	defaultQwenModel   = "qwen3-vl-plus"
)

const (
	vlmSystemPrompt    = "你是一名结构化信息抽取助手，负责从气象监测照片中的水印或字幕里提取时间与地点。务必只输出严格符合要求的 JSON。"
	vlmUserInstruction = `请阅读这张气象监测照片右上或右下角的文字水印，提取拍摄时间和测站地点。如果无法确定某个字段，请填空字符串并将 confidence 设置为 0。
返回 JSON，字段说明：
{
  "time": "24小时制时间戳，格式为 YYYY-MM-DD HH:MM[:SS]，若无法确定则为空字符串",
  "location": "地点中文名称，包含省市县或测站名称，无法确定则为空字符串",
  "confidence": 小数，0-1 之间，表示整体提取置信度,
  "notes": "可选，说明判断依据，若无可留空"
}
仅返回 JSON，不要添加其它文字。`
)

// QwenVLMClient 封装通义千问 VLM 的调用
type QwenVLMClient struct {
	APIKey         string
	BaseURL        string
	Model          string
	EnableThinking bool
	ThinkingBudget int
	httpClient     *http.Client
}

type qwenMessage struct {
	Role    string        `json:"role"`
	Content []qwenContent `json:"content"`
}

type qwenContent struct {
	Type     string        `json:"type"`
	Text     string        `json:"text,omitempty"`
	ImageURL *qwenImageURL `json:"image_url,omitempty"`
}

type qwenImageURL struct {
	URL string `json:"url"`
}

type qwenChatRequest struct {
	Model          string        `json:"model"`
	Messages       []qwenMessage `json:"messages"`
	EnableThinking *bool         `json:"enable_thinking,omitempty"`
	ThinkingBudget int           `json:"thinking_budget,omitempty"`
}

type qwenChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// vlmStructuredResult 定义模型返回的结构化内容
type vlmStructuredResult struct {
	Time       string  `json:"time"`
	Location   string  `json:"location"`
	Confidence float64 `json:"confidence"`
	Notes      string  `json:"notes,omitempty"`
}

// OCRResult OCR识别结果
type OCRResult struct {
	Time       string // 识别到的时间
	Location   string // 识别到的地点
	IsStandard bool   // 是否为标准图片（同时有时间和地点）
}

// NewQwenVLMClient 创建 VLM 客户端
func NewQwenVLMClient(apiKey, baseURL, model string, enableThinking bool, thinkingBudget int) *QwenVLMClient {
	client := &QwenVLMClient{
		APIKey:         apiKey,
		BaseURL:        baseURL,
		Model:          model,
		EnableThinking: enableThinking,
		ThinkingBudget: thinkingBudget,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}

	if client.BaseURL == "" {
		client.BaseURL = defaultQwenBaseURL
	}
	if client.Model == "" {
		client.Model = defaultQwenModel
	}

	return client
}

// ExtractMetadata 调用 VLM 模型识别时间、地点
func (c *QwenVLMClient) ExtractMetadata(imagePath string) (*vlmStructuredResult, error) {
	dataURL, err := encodeImageToDataURL(imagePath)
	if err != nil {
		return nil, err
	}

	request := qwenChatRequest{
		Model: c.Model,
		Messages: []qwenMessage{
			{
				Role: "system",
				Content: []qwenContent{{
					Type: "text",
					Text: vlmSystemPrompt,
				}},
			},
			{
				Role: "user",
				Content: []qwenContent{
					{Type: "text", Text: vlmUserInstruction},
					{Type: "image_url", ImageURL: &qwenImageURL{URL: dataURL}},
				},
			},
		},
	}

	if c.EnableThinking {
		enable := true
		request.EnableThinking = &enable
		if c.ThinkingBudget > 0 {
			request.ThinkingBudget = c.ThinkingBudget
		}
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal VLM request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create VLM request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call VLM API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read VLM response: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("VLM API error: status %d, body %s", resp.StatusCode, string(body))
	}

	var chatResp qwenChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to decode VLM response: %w", err)
	}

	if chatResp.Error != nil {
		return nil, fmt.Errorf("VLM API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("VLM API returned no choices")
	}

	rawContent := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	if rawContent == "" {
		return nil, fmt.Errorf("VLM API returned empty content")
	}

	return parseVLMJSON(rawContent)
}

func encodeImageToDataURL(imagePath string) (string, error) {
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(imagePath))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	encoded := base64.StdEncoding.EncodeToString(imageData)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}

func parseVLMJSON(raw string) (*vlmStructuredResult, error) {
	clean := sanitizeJSONBlock(raw)
	if clean == "" {
		return nil, fmt.Errorf("no JSON block found in VLM response")
	}

	var result vlmStructuredResult
	if err := json.Unmarshal([]byte(clean), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VLM JSON: %w", err)
	}

	return &result, nil
}

func sanitizeJSONBlock(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(trimmed, "```") {
		trimmed = strings.TrimPrefix(trimmed, "```json")
		trimmed = strings.TrimPrefix(trimmed, "```")
		trimmed = strings.TrimSpace(trimmed)
		if idx := strings.LastIndex(trimmed, "```"); idx != -1 {
			trimmed = trimmed[:idx]
		}
		trimmed = strings.TrimSpace(trimmed)
	}

	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start >= 0 && end > start {
		return trimmed[start : end+1]
	}

	return trimmed
}

// normalizeTime 标准化时间格式，转换为 YYYY-MM-DD HH:MM:SS 格式
func normalizeTime(timeStr string) string {
	// 移除中文字符并替换为标准分隔符
	timeStr = strings.ReplaceAll(timeStr, "年", "-")
	timeStr = strings.ReplaceAll(timeStr, "月", "-")
	timeStr = strings.ReplaceAll(timeStr, "日", "")

	// 清理多余空格
	timeStr = strings.TrimSpace(timeStr)

	// 处理常见的OCR错误：日期和时间之间缺少空格
	// 例如: "2023-07-2015:47" -> "2023-07-20 15:47"
	// 匹配 YYYY-MM-DD或YYYY-MM-D后面直接跟数字的情况
	timeStr = regexp.MustCompile(`^(\d{4}-\d{2}-\d{1,2})(\d{2}:\d{2})`).ReplaceAllString(timeStr, "$1 $2")
	timeStr = regexp.MustCompile(`^(\d{4}/\d{2}/\d{1,2})(\d{2}:\d{2})`).ReplaceAllString(timeStr, "$1 $2")

	// 清理多余空格
	timeStr = regexp.MustCompile(`\s+`).ReplaceAllString(timeStr, " ")

	// 处理紧凑格式 YYYYMMDDHHMMSS (14位)
	if matched, _ := regexp.MatchString(`^\d{14}$`, timeStr); matched {
		return fmt.Sprintf("%s-%s-%s %s:%s:%s",
			timeStr[0:4], timeStr[4:6], timeStr[6:8],
			timeStr[8:10], timeStr[10:12], timeStr[12:14])
	}

	// 处理紧凑格式 YYYYMMDD HHMMSS (带空格的14位)
	if matched, _ := regexp.MatchString(`^\d{8}\s+\d{6}$`, timeStr); matched {
		cleaned := strings.ReplaceAll(timeStr, " ", "")
		return fmt.Sprintf("%s-%s-%s %s:%s:%s",
			cleaned[0:4], cleaned[4:6], cleaned[6:8],
			cleaned[8:10], cleaned[10:12], cleaned[12:14])
	}

	// 处理紧凑格式 YYYYMMDDHHMM (12位，没有秒)
	if matched, _ := regexp.MatchString(`^\d{12}$`, timeStr); matched {
		return fmt.Sprintf("%s-%s-%s %s:%s:00",
			timeStr[0:4], timeStr[4:6], timeStr[6:8],
			timeStr[8:10], timeStr[10:12])
	}

	// 处理紧凑格式 YYYYMMDD HHMM (带空格的12位)
	if matched, _ := regexp.MatchString(`^\d{8}\s+\d{4}$`, timeStr); matched {
		cleaned := strings.ReplaceAll(timeStr, " ", "")
		return fmt.Sprintf("%s-%s-%s %s:%s:00",
			cleaned[0:4], cleaned[4:6], cleaned[6:8],
			cleaned[8:10], cleaned[10:12])
	}

	// 处理只有日期的格式 YYYYMMDD (8位)
	if matched, _ := regexp.MatchString(`^\d{8}$`, timeStr); matched {
		return fmt.Sprintf("%s-%s-%s 00:00:00",
			timeStr[0:4], timeStr[4:6], timeStr[6:8])
	}

	// 处理标准格式带秒 YYYY-MM-DD HH:MM:SS (已经是标准格式)
	if matched, _ := regexp.MatchString(`^\d{4}-\d{1,2}-\d{1,2}\s+\d{1,2}:\d{2}:\d{2}$`, timeStr); matched {
		// 补齐单位数的月、日、时
		parts := regexp.MustCompile(`^(\d{4})-(\d{1,2})-(\d{1,2})\s+(\d{1,2}):(\d{2}):(\d{2})$`).FindStringSubmatch(timeStr)
		if len(parts) == 7 {
			return fmt.Sprintf("%s-%02s-%02s %02s:%s:%s",
				parts[1], parts[2], parts[3], parts[4], parts[5], parts[6])
		}
		return timeStr
	}

	// 处理标准格式但没有秒 YYYY-MM-DD HH:MM 或 YYYY-MM-D HH:MM
	if matched, _ := regexp.MatchString(`^\d{4}-\d{1,2}-\d{1,2}\s+\d{1,2}:\d{2}$`, timeStr); matched {
		parts := regexp.MustCompile(`^(\d{4})-(\d{1,2})-(\d{1,2})\s+(\d{1,2}):(\d{2})$`).FindStringSubmatch(timeStr)
		if len(parts) == 6 {
			return fmt.Sprintf("%s-%02s-%02s %02s:%s:00",
				parts[1], parts[2], parts[3], parts[4], parts[5])
		}
		return timeStr + ":00"
	}

	// 处理只有日期 YYYY-MM-DD 或 YYYY-MM-D
	if matched, _ := regexp.MatchString(`^\d{4}-\d{1,2}-\d{1,2}$`, timeStr); matched {
		parts := regexp.MustCompile(`^(\d{4})-(\d{1,2})-(\d{1,2})$`).FindStringSubmatch(timeStr)
		if len(parts) == 4 {
			return fmt.Sprintf("%s-%02s-%02s 00:00:00", parts[1], parts[2], parts[3])
		}
		return timeStr + " 00:00:00"
	}

	// 处理斜杠格式带秒 YYYY/MM/DD HH:MM:SS
	if matched, _ := regexp.MatchString(`^\d{4}/\d{1,2}/\d{1,2}\s+\d{1,2}:\d{2}:\d{2}$`, timeStr); matched {
		parts := regexp.MustCompile(`^(\d{4})/(\d{1,2})/(\d{1,2})\s+(\d{1,2}):(\d{2}):(\d{2})$`).FindStringSubmatch(timeStr)
		if len(parts) == 7 {
			return fmt.Sprintf("%s-%02s-%02s %02s:%s:%s",
				parts[1], parts[2], parts[3], parts[4], parts[5], parts[6])
		}
		timeStr = strings.ReplaceAll(timeStr, "/", "-")
		return timeStr
	}

	// 处理斜杠格式但没有秒 YYYY/MM/DD HH:MM
	if matched, _ := regexp.MatchString(`^\d{4}/\d{1,2}/\d{1,2}\s+\d{1,2}:\d{2}$`, timeStr); matched {
		parts := regexp.MustCompile(`^(\d{4})/(\d{1,2})/(\d{1,2})\s+(\d{1,2}):(\d{2})$`).FindStringSubmatch(timeStr)
		if len(parts) == 6 {
			return fmt.Sprintf("%s-%02s-%02s %02s:%s:00",
				parts[1], parts[2], parts[3], parts[4], parts[5])
		}
		timeStr = strings.ReplaceAll(timeStr, "/", "-")
		return timeStr + ":00"
	}

	// 处理斜杠格式只有日期 YYYY/MM/DD 或 YYYY/MM/D
	if matched, _ := regexp.MatchString(`^\d{4}/\d{1,2}/\d{1,2}`, timeStr); matched {
		parts := regexp.MustCompile(`^(\d{4})/(\d{1,2})/(\d{1,2})`).FindStringSubmatch(timeStr)
		if len(parts) == 4 {
			return fmt.Sprintf("%s-%02s-%02s 00:00:00", parts[1], parts[2], parts[3])
		}
		timeStr = strings.ReplaceAll(timeStr, "/", "-")
	}

	// 如果已经是标准格式，直接返回
	if matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}$`, timeStr); matched {
		return timeStr
	}

	return timeStr
}

// cleanLocationText 清理地点文字
func cleanLocationText(text string) string {
	// 移除可能的特殊字符和多余空格
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	return text
}

// ProcessImageOCR 处理图片OCR（主入口函数）
func ProcessImageOCR(imagePath string) (*OCRResult, error) {
	apiKey := getEnv("QWEN_VLM_API_KEY", "")
	if apiKey == "" {
		log.Println("Warning: Qwen VLM API key not configured, skipping OCR processing")
		return &OCRResult{IsStandard: false}, nil
	}

	baseURL := getEnv("QWEN_VLM_BASE_URL", defaultQwenBaseURL)
	model := getEnv("QWEN_VLM_MODEL", defaultQwenModel)
	enableThinking := strings.EqualFold(getEnv("QWEN_VLM_ENABLE_THINKING", ""), "true")
	thinkingBudget := getEnvInt("QWEN_VLM_THINKING_BUDGET", 0)

	client := NewQwenVLMClient(apiKey, baseURL, model, enableThinking, thinkingBudget)
	structured, err := client.ExtractMetadata(imagePath)
	if err != nil {
		return nil, fmt.Errorf("VLM extraction failed: %w", err)
	}

	result := &OCRResult{}
	if strings.TrimSpace(structured.Time) != "" {
		result.Time = normalizeTime(structured.Time)
	}
	if strings.TrimSpace(structured.Location) != "" {
		result.Location = cleanLocationText(structured.Location)
	}
	result.IsStandard = result.Time != "" && result.Location != ""

	log.Printf("VLM Result - Time: %s, Location: %s, Confidence: %.2f, IsStandard: %v",
		result.Time, result.Location, structured.Confidence, result.IsStandard)

	return result, nil
}
