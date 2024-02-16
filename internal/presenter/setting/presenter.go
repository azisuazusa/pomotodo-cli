package setting

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/azisuazusa/todo-cli/internal/domain/syncintegration"
	"github.com/manifoldco/promptui"
)

const DROPBOX_API_KEY = "aau1qknsx7d5zq7"

type Presenter struct {
	settingUseCase syncintegration.UseCase
}

func New(settingUseCase syncintegration.UseCase) *Presenter {
	return &Presenter{
		settingUseCase: settingUseCase,
	}
}

func (p *Presenter) SetSyncIntegration(ctx context.Context) error {
	prompt := promptui.Select{
		Label: "Sync Integration",
		Items: []string{string(syncintegration.Dropbox)},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	integration := syncintegration.SyncIntegration{
		Type:    syncintegration.SyncIntegrationType(result),
		Details: map[string]string{},
	}

	if integration.Type == syncintegration.Dropbox {

		token, err := p.dropboxIntegration(ctx)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}

		integration.Details["token"] = token
	}

	if err = p.settingUseCase.SetSyncIntegration(ctx, integration); err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	if err = p.settingUseCase.Upload(ctx); err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	fmt.Println("Integration set successfully!")

	return nil
}

func (p *Presenter) Sync(ctx context.Context) error {
	err := p.settingUseCase.Download(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	currentFile, err := os.ReadFile(homeDir + "/.todo-cli.db")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	remoteFile, err := os.ReadFile(homeDir + "/.todo-cli-remote.db")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	currentChecksum, err := calculateChecksum(currentFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	remoteChecksum, err := calculateChecksum(remoteFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	if currentChecksum == remoteChecksum {
		fmt.Println("Already up to date!")
		return nil
	}

	prompt := promptui.Select{
		Label: "There are differences between local and remote. What do you want to do?",
		Items: []string{"Upload", "Download", "Cancel"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	if result == "Upload" {
		err = p.settingUseCase.Upload(ctx)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}

		if err = os.Remove(homeDir + "/.todo-cli-remote.db"); err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}

		fmt.Println("Upload successful!")
		return nil
	}

	if result == "Download" {
		if err = os.Remove(homeDir + "/.todo-cli.db"); err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}

		if err = os.Rename(homeDir+"/.todo-cli-remote.db", homeDir+"/.todo-cli.db"); err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}

		fmt.Println("Download successful!")
		return nil
	}

	fmt.Println("Sync cancelled!")
	return nil

}

func calculateChecksum(file []byte) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, bytes.NewReader(file)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (p *Presenter) dropboxIntegration(ctx context.Context) (token string, err error) {
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	codeChallenge := generateCodeChallenge(codeVerifier)
	authorizeURL := constructAuthURL(codeChallenge)

	fmt.Printf("1. Go to %s\n", authorizeURL)
	fmt.Println("2. Click on 'Allow' (you might have to log in first)")
	fmt.Println("3. Copy the code that appears on the next page")
	fmt.Println("4. Paste the code here")

	prompt := promptui.Prompt{
		Label: "Dropbox Code",
	}

	code, promptErr := prompt.Run()
	if promptErr != nil {
		err = promptErr
		fmt.Printf("Error: %v\n", err)
		return
	}

	data := url.Values{}
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", DROPBOX_API_KEY)
	data.Set("code_verifier", codeVerifier)
	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/oauth2/token", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Body = io.NopCloser(strings.NewReader(data.Encode()))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Error: %v\n", string(body))
		err = errors.New("Error getting access token")
		return
	}

	var dropboxTokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err = json.Unmarshal(body, &dropboxTokenResponse); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	token = dropboxTokenResponse.AccessToken

	return
}
func generateCodeVerifier() (string, error) {
	verifier := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, verifier)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(verifier), nil
}

func generateCodeChallenge(verifier string) string {
	hashed := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hashed[:])
}

func constructAuthURL(codeChallenge string) (url string) {
	authorizeURL := "https://www.dropbox.com/oauth2/authorize"
	urlParams := fmt.Sprintf("response_type=code&client_id=%s&code_challenge=%s&code_challenge_method=S256&token_access_type=offline", DROPBOX_API_KEY, codeChallenge)
	url = fmt.Sprintf("%s?%s", authorizeURL, urlParams)
	return
}
