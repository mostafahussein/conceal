package config

import (
	b64 "encoding/base64"
	"encoding/json"
	"os"
	"time"

	"conceal/internal/logging"
	gopenpgp "github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/miquella/ask"
)

func FetchPasswordManager() (string, string, error) {
	jsonParsed, err := LoadKoanf()
	if err != nil {
		return "", "", err
	}
	passwordManager := jsonParsed.Search("provider", "0", "name").Data().(string)
	vaultLocation := jsonParsed.Search("provider", "0", "vault_location").Data().(string)

	return passwordManager, vaultLocation, nil

}

func fetchAuthSecret() (string, error) {
	jsonParsed, err := LoadKoanf()
	if err != nil {
		return "", err
	}
	authSecret := jsonParsed.Search("provider", "0", "auth_secret").Data().(string)

	return authSecret, nil

}

func fetchSessionTimeOut() (float64, error) {
	jsonParsed, err := LoadKoanf()
	if err != nil {
		return 0, err
	}
	sessionTimeout := jsonParsed.Search("provider", "0", "session_timeout").Data().(float64)

	return sessionTimeout, nil

}

func passwordExpiration(timeout float64) time.Time {
	expirationTimestamp := time.Duration(timeout) * time.Minute
	passwordExpiration := time.Now().Local().Add(expirationTimestamp - 1*time.Minute)
	return passwordExpiration
}

func encryptPassword(password string) (string, error) {
	authSecret, _ := fetchAuthSecret()
	gpgSecret := []byte(authSecret)
	encryptedPassword, err := gopenpgp.EncryptMessageWithPassword(gpgSecret, password)
	encodedString := b64.StdEncoding.EncodeToString([]byte(encryptedPassword))
	return encodedString, err
}

func savePassword(payload string, expiration string) {
	data := map[string]string{"password": payload, "expires_at": expiration}
	file, _ := json.MarshalIndent(data, "", " ")
	if _, err := os.Stat(AppPath); os.IsNotExist(err) {
		_ = os.Mkdir(AppPath, 0755)
	}
	_ = os.WriteFile(credentialsFilePath, file, 0644)
}

func decryptPassword(encryptedPassword string) (string, error) {
	authSecret, _ := fetchAuthSecret()
	gpgSecret := []byte(authSecret)
	decodedString, _ := b64.StdEncoding.DecodeString(encryptedPassword)
	gpgValue := string(decodedString[:])
	decryptedPassword, err := gopenpgp.DecryptMessageWithPassword(gpgSecret, gpgValue)
	return decryptedPassword, err
}

func PasswordManagerPrompt() string {
	var isCredentialsExist bool
	var fileData map[string]string
	var isValidSession bool
	var vaultPassword string
	if _, err := os.Stat(credentialsFilePath); !os.IsNotExist(err) {
		isCredentialsExist = true
		file, _ := os.ReadFile(credentialsFilePath)
		_ = json.Unmarshal(file, &fileData)
	}
	if isCredentialsExist {
		expirationTime := fileData["expires_at"]
		currentTime := time.Now().Local()
		parsedExpirationTime, _ := time.Parse(time.RFC3339, expirationTime)
		if currentTime.After(parsedExpirationTime) {
			isValidSession = false
		} else {
			isValidSession = true
		}
	}
	if isCredentialsExist && isValidSession {
		decryptedPassword, _ := decryptPassword(fileData["password"])
		vaultPassword = decryptedPassword
	} else {
		if response, err := ask.HiddenAsk("Unlock your vault: "); err != nil {
			logging.Logger.Fatal().Msg("unable to get vault password")
		} else {
			vaultPassword = response
		}
	}
	return vaultPassword
}

func SavePassword(password string) {
	timeout, _ := fetchSessionTimeOut()
	expiration := passwordExpiration(timeout).Format(time.RFC3339)
	encryptedPassword, _ := encryptPassword(password)
	savePassword(encryptedPassword, expiration)
}
