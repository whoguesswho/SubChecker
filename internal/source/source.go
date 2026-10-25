package source

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Fetch решает, откуда брать данные (сеть / файл)
func Fetch(pathOrURL string) ([]string, error) {
	var rawData []byte
	var err error

	// Проверка строки, является ли она юрлом
	if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
		rawData, err = downloadFromURL(pathOrURL)
	} else {
		// Иначе - локальный путь
		rawData, err = os.ReadFile(pathOrURL)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	// Отдаем сырой текст в функцию очистки
	return cleanText(string(rawData)), nil
}

func downloadFromURL(url string) ([]byte, error) {
	// ! Таймаут, на случай зависания сервера
	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	// defer закрывает тело ответа при выходе из функции
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	// Читаем всё в память
	return io.ReadAll(resp.Body)
}

// cleanText чистит текст + декод из base64 при необходимости
func cleanText(text string) []string {
	// Шаг 1. Пробуем декодировать Base64
	decodedText := tryDecodeBase64(text)
	// Если внутри есть "://", значит это - закодированные ссылки
	if strings.Contains(decodedText, "://") {
		text = decodedText
	}

	// Шаг 2. Разбиваем на строки
	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		result = append(result, line)
	}

	return result
}

func tryDecodeBase64(text string) string {
	text = strings.TrimSpace(text)

	encodings := []*base64.Encoding{
		base64.StdEncoding,
		base64.RawStdEncoding,
		base64.URLEncoding,
		base64.RawURLEncoding,
	}

	for _, enc := range encodings {
		if decoded, err := enc.DecodeString(text); err == nil {
			return string(decoded)
		}
	}
	return text
}
