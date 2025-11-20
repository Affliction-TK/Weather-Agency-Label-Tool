package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// BaiduOCRClient 百度OCR客户端
type BaiduOCRClient struct {
	APIKey      string
	SecretKey   string
	AccessToken string
	TokenExpiry time.Time
}

// BaiduTokenResponse 百度token响应
type BaiduTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Error       string `json:"error,omitempty"`
	ErrorDesc   string `json:"error_description,omitempty"`
}

// BaiduOCRResponse 百度OCR响应
type BaiduOCRResponse struct {
	LogID          uint64         `json:"log_id"`
	WordsResultNum int            `json:"words_result_num"`
	WordsResult    []BaiduOCRWord `json:"words_result"`
	ErrorCode      int            `json:"error_code,omitempty"`
	ErrorMsg       string         `json:"error_msg,omitempty"`
}

// BaiduOCRWord OCR识别的单个文字块
type BaiduOCRWord struct {
	Words    string           `json:"words"`
	Location BaiduOCRLocation `json:"location"`
}

// BaiduOCRLocation 文字位置信息
type BaiduOCRLocation struct {
	Left   int `json:"left"`
	Top    int `json:"top"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// OCRResult OCR识别结果
type OCRResult struct {
	Time       string // 识别到的时间
	Location   string // 识别到的地点
	IsStandard bool   // 是否为标准图片（同时有时间和地点）
}

// NewBaiduOCRClient 创建百度OCR客户端
func NewBaiduOCRClient(apiKey, secretKey string) *BaiduOCRClient {
	return &BaiduOCRClient{
		APIKey:    apiKey,
		SecretKey: secretKey,
	}
}

// GetAccessToken 获取访问令牌
func (c *BaiduOCRClient) GetAccessToken() error {
	// 如果token还有效，直接返回
	if c.AccessToken != "" && time.Now().Before(c.TokenExpiry) {
		return nil
	}

	tokenURL := fmt.Sprintf(
		"https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s",
		c.APIKey, c.SecretKey,
	)

	resp, err := http.Get(tokenURL)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp BaiduTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenResp.Error != "" {
		return fmt.Errorf("token error: %s - %s", tokenResp.Error, tokenResp.ErrorDesc)
	}

	c.AccessToken = tokenResp.AccessToken
	// 提前5分钟过期，确保安全
	c.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-300) * time.Second)

	return nil
}

// RecognizeText 识别图片中的文字
func (c *BaiduOCRClient) RecognizeText(imagePath string) (*BaiduOCRResponse, error) {
	// 确保有有效的token
	if err := c.GetAccessToken(); err != nil {
		return nil, err
	}

	// 读取图片并转换为base64
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// URL encode
	encodedImage := url.QueryEscape(base64Image)

	// 构建请求
	ocrURL := fmt.Sprintf(
		"https://aip.baidubce.com/rest/2.0/ocr/v1/general?access_token=%s",
		c.AccessToken,
	)

	data := fmt.Sprintf("image=%s", encodedImage)
	req, err := http.NewRequest("POST", ocrURL, bytes.NewBufferString(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send OCR request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var ocrResp BaiduOCRResponse
	if err := json.Unmarshal(body, &ocrResp); err != nil {
		return nil, fmt.Errorf("failed to decode OCR response: %w", err)
	}

	if ocrResp.ErrorCode != 0 {
		return nil, fmt.Errorf("OCR error %d: %s", ocrResp.ErrorCode, ocrResp.ErrorMsg)
	}

	return &ocrResp, nil
}

// ExtractTimeAndLocation 从OCR结果中提取时间和地点
func ExtractTimeAndLocation(ocrResp *BaiduOCRResponse) *OCRResult {
	result := &OCRResult{}

	if ocrResp == nil || len(ocrResp.WordsResult) == 0 {
		return result
	}

	// 更宽松的时间正则表达式，支持更多格式
	timePatterns := []*regexp.Regexp{
		// 标准格式：2024-01-15 14:30:45 或 2024/01/15 14:30:45
		regexp.MustCompile(`\d{4}[-/年]\d{1,2}[-/月]\d{1,2}[日\s]*\d{1,2}:\d{2}:\d{2}`),
		// 没有秒：2024-01-15 14:30
		regexp.MustCompile(`\d{4}[-/年]\d{1,2}[-/月]\d{1,2}[日\s]*\d{1,2}:\d{2}`),
		// 紧凑格式：20240115143045 或 20240115 143045
		regexp.MustCompile(`\d{4}\d{2}\d{2}\s*\d{2}\d{2}\d{2}`),
		// 更短的紧凑格式：20240115 1430
		regexp.MustCompile(`\d{4}\d{2}\d{2}\s*\d{2}\d{2}`),
		// 只包含日期和部分数字（宽松匹配）
		regexp.MustCompile(`\d{4}[-/年]\d{1,2}[-/月]\d{1,2}`),
	}

	// 地点关键词
	locationKeywords := []string{"省", "市", "县", "区", "站", "路", "街", "镇", "乡", "村"}

	// 分离上部区域（右上角）和下部区域（右下角）的文字块
	// 按照 top 位置排序，前30%为上部区域，后30%为下部区域
	var topWords, bottomWords []BaiduOCRWord
	if len(ocrResp.WordsResult) > 0 {
		// 找出最大和最小的top值
		minTop, maxTop := ocrResp.WordsResult[0].Location.Top, ocrResp.WordsResult[0].Location.Top
		for _, word := range ocrResp.WordsResult {
			if word.Location.Top < minTop {
				minTop = word.Location.Top
			}
			if word.Location.Top > maxTop {
				maxTop = word.Location.Top
			}
		}

		// 计算阈值
		heightRange := maxTop - minTop
		upperThreshold := minTop + heightRange/3 // 上部1/3区域
		lowerThreshold := maxTop - heightRange/3 // 下部1/3区域

		for _, word := range ocrResp.WordsResult {
			if word.Location.Top <= upperThreshold {
				topWords = append(topWords, word)
			} else if word.Location.Top >= lowerThreshold {
				bottomWords = append(bottomWords, word)
			}
		}
	}

	// 优先在上部区域查找时间（右上角）
	for _, word := range topWords {
		text := strings.TrimSpace(word.Words)
		for _, pattern := range timePatterns {
			if match := pattern.FindString(text); match != "" {
				result.Time = normalizeTime(match)
				break
			}
		}
		if result.Time != "" {
			break
		}
	}

	// 如果上部没找到，在所有文字块中查找时间
	if result.Time == "" {
		for _, word := range ocrResp.WordsResult {
			text := strings.TrimSpace(word.Words)
			for _, pattern := range timePatterns {
				if match := pattern.FindString(text); match != "" {
					result.Time = normalizeTime(match)
					break
				}
			}
			if result.Time != "" {
				break
			}
		}
	}

	// 优先在下部区域查找地点（右下角）
	for _, word := range bottomWords {
		text := strings.TrimSpace(word.Words)
		for _, keyword := range locationKeywords {
			if strings.Contains(text, keyword) {
				result.Location = cleanLocationText(text)
				break
			}
		}
		if result.Location != "" {
			break
		}
	}

	// 如果下部没找到，在所有文字块中查找地点
	if result.Location == "" {
		for _, word := range ocrResp.WordsResult {
			text := strings.TrimSpace(word.Words)
			for _, keyword := range locationKeywords {
				if strings.Contains(text, keyword) {
					result.Location = cleanLocationText(text)
					break
				}
			}
			if result.Location != "" {
				break
			}
		}
	}

	// 判断是否为标准图片
	result.IsStandard = result.Time != "" && result.Location != ""

	return result
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
	// 获取OCR配置
	apiKey := getEnv("BAIDU_OCR_API_KEY", "")
	secretKey := getEnv("BAIDU_OCR_SECRET_KEY", "")

	if apiKey == "" || secretKey == "" {
		log.Println("Warning: Baidu OCR credentials not configured, skipping OCR processing")
		return &OCRResult{IsStandard: false}, nil
	}

	// 创建OCR客户端
	client := NewBaiduOCRClient(apiKey, secretKey)

	// 调用OCR识别
	ocrResp, err := client.RecognizeText(imagePath)
	if err != nil {
		return nil, fmt.Errorf("OCR recognition failed: %w", err)
	}

	// 提取时间和地点
	result := ExtractTimeAndLocation(ocrResp)

	log.Printf("OCR Result - Time: %s, Location: %s, IsStandard: %v",
		result.Time, result.Location, result.IsStandard)

	return result, nil
}
