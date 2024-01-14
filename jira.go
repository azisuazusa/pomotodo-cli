package main

import (
	"fmt"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/manifoldco/promptui"
)

func initJIRAClient(project Project) (client *jira.Client, err error) {
	jiraAuth := jira.BasicAuthTransport{
		Username: project.JIRA.Username,
		Password: project.JIRA.Token,
	}

	client, err = jira.NewClient(jiraAuth.Client(), project.JIRA.URL)
	if err != nil {
		fmt.Println("Error creating JIRA client: ", err)
		return
	}

	return
}

func AddWorklogToJIRAIssue(task Task) (err error) {
	projects := Projects{}
	if err = projects.load(); err != nil {
		return
	}

	projectIndex := projects.getSelectedIndex()
	client, err := initJIRAClient(projects[projectIndex])
	if err != nil {
		fmt.Println("Error initializing JIRA client: ", err)
		return
	}

	var timeSpent time.Duration
	for _, history := range task.TaskHistories {
		timeSpent += history.StoppedAt.Sub(history.StartedAt)
	}

	prompt := promptui.Prompt{
		Label:     "Add worklog to JIRA issue, time spent in hours",
		Default:   fmt.Sprintf("%.0f", timeSpent.Hours()),
		AllowEdit: true,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Println("Error adding worklog to JIRA issue: ", err)
		return
	}

	if result != fmt.Sprintf("%.0f", timeSpent.Hours()) {
		timeSpent, err = time.ParseDuration(fmt.Sprintf("%sh", result))
		if err != nil {
			fmt.Println("Error adding worklog to JIRA issue: ", err)
			return
		}
	}

	worklog := jira.WorklogRecord{
		TimeSpentSeconds: int(timeSpent.Seconds()),
		Comment:          task.Name,
	}

	_, _, err = client.Issue.AddWorklogRecord(task.ID, &worklog)
	if err != nil {
		fmt.Println("Error adding worklog to JIRA issue: ", err)
		return
	}

	return
}

func SyncJIRAIssues() (err error) {
	projects := Projects{}
	if err = projects.load(); err != nil {
		return
	}

	projectIndex := projects.getSelectedIndex()
	client, err := initJIRAClient(projects[projectIndex])
	if err != nil {
		fmt.Println("Error initializing JIRA client: ", err)
		return
	}

	last := 0
	var issues []jira.Issue
	for {
		opt := &jira.SearchOptions{
			StartAt:    last,
			MaxResults: 1000,
		}

		chunk, resp, errSearch := client.Issue.Search(projects[projectIndex].JIRA.JQL, opt)
		if errSearch != nil {
			err = errSearch
			fmt.Println("Error searching JIRA: ", err)
			return
		}

		total := resp.Total
		if issues == nil {
			issues = make([]jira.Issue, 0, total)
		}

		issues = append(issues, chunk...)
		last = resp.StartAt + len(chunk)
		if last >= resp.Total {
			break
		}
	}

	for _, issue := range issues {
		isFoundIndex := -1
		for i, task := range projects[projectIndex].Tasks {
			if task.ID == issue.Key {
				isFoundIndex = i
				break
			}
		}

		if isFoundIndex != -1 {
			projects[projectIndex].Tasks[isFoundIndex].Name = issue.Fields.Summary
			projects[projectIndex].Tasks[isFoundIndex].Description = issue.Fields.Description
			continue
		}

		task := Task{
			ID:          issue.Key,
			Name:        issue.Fields.Summary,
			Description: issue.Fields.Description,
			IsJIRATask:  true,
		}

		projects[projectIndex].Tasks = append(projects[projectIndex].Tasks, task)
	}

	if err = projects.save(); err != nil {
		return
	}

	return nil

}
