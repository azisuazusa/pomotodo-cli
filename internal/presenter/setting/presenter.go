package setting

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/azisuazusa/todo-cli/internal/domain/syncintegration"
	"github.com/manifoldco/promptui"
)

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
		prompt := promptui.Prompt{
			Label: "Dropbox token",
		}

		token, promptErr := prompt.Run()
		if promptErr != nil {
			err = promptErr
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
