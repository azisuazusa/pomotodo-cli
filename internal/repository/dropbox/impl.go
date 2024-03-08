package dropbox

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	syncintegrationDomain "github.com/azisuazusa/todo-cli/internal/domain/syncintegration"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

const DROPBOX_API_KEY = "aau1qknsx7d5zq7"

type SyncIntegrationRepo interface {
	SetSyncIntegration(ctx context.Context, integration syncintegrationDomain.SyncIntegration) error
}

type RepoImpl struct {
	syncIntegrationRepo SyncIntegrationRepo
}

func New(syncIntegrationRepo SyncIntegrationRepo) *RepoImpl {
	return &RepoImpl{
		syncIntegrationRepo: syncIntegrationRepo,
	}
}

func (ri *RepoImpl) Download(ctx context.Context, integrationEntity syncintegrationDomain.SyncIntegration) error {
	dbxCfg := dropbox.Config{
		Token: integrationEntity.Details["token"],
	}
	fileDbx := files.New(dbxCfg)
	downloadArg := files.NewDownloadArg("/.todo-cli.db")
	_, content, err := fileDbx.Download(downloadArg)
	if err != nil && !strings.Contains(err.Error(), "expired") {
		return fmt.Errorf("error while downloading file: %w", err)
	}

	if err != nil && strings.Contains(err.Error(), "expired") {
		err = ri.refreshToken(ctx, integrationEntity)
		if err != nil {
			return fmt.Errorf("error while refreshing token: %w", err)
		}

		return ri.Download(ctx, integrationEntity)
	}

	contentBytes, err := io.ReadAll(content)
	if err != nil {
		return fmt.Errorf("error while reading file content: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error while getting user home dir: %w", err)
	}

	err = os.WriteFile(homeDir+"/.todo-cli-remote.db", contentBytes, 0644)
	if err != nil {
		return fmt.Errorf("error while writing file: %w", err)
	}

	return nil
}

func (ri *RepoImpl) Upload(ctx context.Context, integrationEntity syncintegrationDomain.SyncIntegration) error {
	dbxCfg := dropbox.Config{
		Token: integrationEntity.Details["token"],
	}
	fileDbx := files.New(dbxCfg)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error while getting user home dir: %w", err)
	}

	file, err := os.Open(homeDir + "/.todo-cli.db")
	if err != nil {
		return fmt.Errorf("error while opening file: %w", err)
	}

	uploadArg := files.NewUploadArg("/.todo-cli.db")
	uploadArg.CommitInfo.Mode.Tag = files.WriteModeOverwrite
	_, err = fileDbx.Upload(uploadArg, file)
	if err != nil && !strings.Contains(err.Error(), "expired") {
		return fmt.Errorf("error while uploading file: %w", err)
	}

	if err != nil && strings.Contains(err.Error(), "expired") {
		err = ri.refreshToken(ctx, integrationEntity)
		if err != nil {
			return fmt.Errorf("error while refreshing token: %w", err)
		}

		return ri.Upload(ctx, integrationEntity)
	}

	return err
}

func (ri *RepoImpl) refreshToken(ctx context.Context, integrationEntity syncintegrationDomain.SyncIntegration) error {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", integrationEntity.Details["refresh_token"])
	data.Set("client_id", DROPBOX_API_KEY)

	req, err := http.NewRequest("POST", "https://api.dropbox.com/oauth2/token", nil)
	if err != nil {
		return fmt.Errorf("error while creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Body = io.NopCloser(strings.NewReader(data.Encode()))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error while making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error while reading response: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error getting access token: %s", string(body))
	}

	var dropboxTokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err = json.Unmarshal(body, &dropboxTokenResponse); err != nil {
		return fmt.Errorf("error while unmarshalling response: %w", err)
	}

	integrationEntity.Details["token"] = dropboxTokenResponse.AccessToken
	integrationEntity.Details["refresh_token"] = integrationEntity.Details["refresh_token"]
	integrationEntity.Details["token_type"] = dropboxTokenResponse.TokenType
	integrationEntity.Details["expires_in"] = fmt.Sprintf("%d", dropboxTokenResponse.ExpiresIn)

	err = ri.syncIntegrationRepo.SetSyncIntegration(ctx, integrationEntity)
	if err != nil {
		return fmt.Errorf("error while setting dropbox token: %w", err)
	}

	return nil
}
