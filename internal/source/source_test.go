package source

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

// Тест 1: Очистка от мусора
func TestCleanText(t *testing.T) {
	// Грязная строка
	raw := "vless://test1\r\n\n#комментарий\r\n  trojan://test2  \n"
	expected := []string{"vless://test1", "trojan://test2"}

	actual := cleanText(raw)

	// reflect.DeepEqual сравнивает массивы поэлементно
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Ожидали %v, получили %v", expected, actual)
	}
}

// Тест 2: Декодирование Base64
func TestCleanTextBase64(t *testing.T) {
	raw := "dmxlc3M6Ly90ZXN0MQ=="
	expected := []string{"vless://test1"}

	actual := cleanText(raw)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Ожидали %v, получили %v", expected, actual)
	}
}

// Тест 3: Чтение локального файла
func TestFetchFromFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("vless://file-test\n")
	tmpFile.Close()

	links, err := Fetch(tmpFile.Name())
	if err != nil {
		t.Fatalf("Ошибка: %v", err)
	}

	if len(links) != 1 || links[0] != "vless://file-test" {
		t.Errorf("Неверный результат: %v", links)
	}
}

// Тест 4: Скачивание по URL
func TestFetchFromURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ss://url-test\n"))
	}))
	defer server.Close()

	links, err := Fetch(server.URL)
	if err != nil {
		t.Fatalf("Ошибка: %v", err)
	}

	if len(links) != 1 || links[0] != "ss://url-test" {
		t.Errorf("Неверный результат: %v", links)
	}
}
