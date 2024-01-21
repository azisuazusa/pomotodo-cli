package dropbox

import (
	"context"
	"fmt"
	"io"
	"os"

	settingDomain "github.com/azisuazusa/todo-cli/internal/domain/setting"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type RepoImpl struct{}

func New() *RepoImpl {
	return &RepoImpl{}
}

func (ri *RepoImpl) Download(ctx context.Context, integrationEntity settingDomain.Integration) error {
	dbxCfg := dropbox.Config{
		Token: integrationEntity.Details["token"],
	}
	fileDbx := files.New(dbxCfg)
	downloadArg := files.NewDownloadArg("/.todo-cli.db")
	_, content, err := fileDbx.Download(downloadArg)
	if err != nil {
		return fmt.Errorf("error while downloading file: %w", err)
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

func (ri *RepoImpl) Upload(ctx context.Context, integrationEntity settingDomain.Integration) error {
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
	_, err = fileDbx.Upload(uploadArg, file)
	if err != nil {
		return fmt.Errorf("error while uploading file: %w", err)
	}

	return nil
}
