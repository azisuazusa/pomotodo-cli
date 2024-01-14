package main

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/manifoldco/promptui"
)

func SyncDropbox() error {
	projects := Projects{}
	if err := projects.load(); err != nil {
		return err
	}

	projectIndex := projects.getSelectedIndex()
	dropboxCfg := dropbox.Config{
		Token: projects[projectIndex].Dropbox.AccessToken,
	}

	dbx := files.New(dropboxCfg)
	downloadArg := files.NewDownloadArg("/.todo-projects.json")
	_, content, err := dbx.Download(downloadArg)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return err
	}

	contentBytes, err := io.ReadAll(content)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	var projectsFromDropbox Projects
	if err = json.Unmarshal(contentBytes, &projectsFromDropbox); err != nil {
		fmt.Println("Error unmarshalling file:", err)
		return err
	}

	if reflect.DeepEqual(projects, projectsFromDropbox) {
		fmt.Println("No changes to sync")
		return nil
	}

	prompt := promptui.Select{
		Label: "Select an action",
		Items: []string{"Download", "Upload", "Cancel"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return err
	}

	switch result {
	case "Download":
		if err := projectsFromDropbox.save(); err != nil {
			return err
		}
	case "Upload":
		if err := projects.save(); err != nil {
			return err
		}

		jsonData, err := json.MarshalIndent(projects, "", "  ")
		if err != nil {
			fmt.Println("Error marshalling file:", err)
			return err
		}

		deleteArg := files.NewDeleteArg("/.todo-projects.json")
		_, err = dbx.DeleteV2(deleteArg)
		if err != nil {
			fmt.Println("Error deleting file:", err)
			return err
		}

		reader := strings.NewReader(string(jsonData))
		uploadArg := files.NewUploadArg("/.todo-projects.json")
		_, err = dbx.Upload(uploadArg, reader)
		if err != nil {
			fmt.Println("Error uploading file:", err)
			return err
		}

	case "Cancel":
		fmt.Println("Cancelled")
	}

	return nil
}
