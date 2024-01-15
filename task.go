package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
)

type TaskHistory struct {
	StartedAt time.Time `json:"started_at"`
	StoppedAt time.Time `json:"stopped_at"`
}

type Task struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	IsStarted     bool          `json:"is_started"`
	IsCompleted   bool          `json:"is_completed"`
	TaskHistories []TaskHistory `json:"task_histories"`
	ParentTaskID  string        `json:"parent_task_id"`
	IsJIRATask    bool          `json:"is_jira_task"`
}

type Tasks []Task

func (t *Tasks) getIndexByID(id string) int {
	for i, task := range *t {
		if task.ID == id {
			return i
		}
	}

	return -1
}

func (t *Tasks) getStartedTaskIndex() int {
	for i, task := range *t {
		if task.IsStarted {
			return i
		}
	}

	return -1
}

func (t *Tasks) filterCompletedTasks() {
	filteredTasks := Tasks{}
	for _, task := range *t {
		if !task.IsCompleted {
			filteredTasks = append(filteredTasks, task)
		}
	}

	*t = filteredTasks
}

func (t *Task) Add() error {
	projects := Projects{}
	if err := projects.load(); err != nil {
		return err
	}

	projectIndex := projects.getSelectedIndex()
	if projectIndex == -1 {
		fmt.Println("No project selected, please select a project first")
		return nil
	}

	id := uuid.New().String()
	t.ID = id
	projects[projectIndex].Tasks = append(projects[projectIndex].Tasks, *t)

	if err := projects.save(); err != nil {
		return err
	}

	fmt.Println("Successfully added task:", t.Name)

	return nil
}

func (t *Task) AddSubTask() error {
	projects := Projects{}
	if err := projects.load(); err != nil {
		return err
	}

	projectIndex := projects.getSelectedIndex()
	if projectIndex == -1 {
		fmt.Println("No project selected, please select a project first")
		return nil
	}

	prompt := promptui.Select{
		Label:     "Select a task to add a subtask to",
		Items:     projects[projectIndex].Tasks,
		Templates: taskSelectTemplate(),
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return err
	}

	id := uuid.New().String()
	t.ID = id
	t.ParentTaskID = projects[projectIndex].Tasks[i].ID

	parentTaskIndex := projects[projectIndex].Tasks.getIndexByID(projects[projectIndex].Tasks[i].ID)
	projects[projectIndex].Tasks = append(projects[projectIndex].Tasks[:parentTaskIndex+1], append(Tasks{*t}, projects[projectIndex].Tasks[parentTaskIndex+1:]...)...)

	if err := projects.save(); err != nil {
		return err
	}

	fmt.Println("Successfully added subtask:", t.Name)

	return nil
}

func TaskComplete() error {
	projects := Projects{}
	if err := projects.load(); err != nil {
		return err
	}

	projectIndex := projects.getSelectedIndex()
	if projectIndex == -1 {
		fmt.Println("No project selected, please select a project first")
		return nil
	}

	prompt := promptui.Select{
		Label:     "Select a task to complete",
		Items:     projects[projectIndex].Tasks,
		Templates: taskSelectTemplate(),
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return err
	}

	projects[projectIndex].Tasks[i].IsCompleted = true
	if projects[projectIndex].Tasks[i].IsStarted {
		projects[projectIndex].Tasks[i].IsStarted = false
		projects[projectIndex].Tasks[i].TaskHistories[len(projects[projectIndex].Tasks[i].TaskHistories)-1].StoppedAt = time.Now()
	}

	project := projects[projectIndex]
	if project.JIRA.URL == "" {
		if err = projects.save(); err != nil {
			return err
		}

		fmt.Println("Successfully completed task:", projects[projectIndex].Tasks[i].Name)
		return nil
	}

	taskID := ""
	task := project.Tasks[i]
	if task.IsJIRATask {
		taskID = task.ID
	}

	if task.ParentTaskID != "" {
		parentTask := project.Tasks[project.Tasks.getIndexByID(project.Tasks[i].ParentTaskID)]
		if parentTask.IsJIRATask {
			taskID = parentTask.ID
		}
	}

	if taskID != "" {
		err = AddWorklogToJIRAIssue(taskID, task.Name, task.TaskHistories)
		if err != nil {
			fmt.Println("Failed to add worklog to JIRA issue:", err)
		}
	}

	if err = projects.save(); err != nil {
		return err
	}

	fmt.Println("Successfully completed task:", task.Name)

	return nil
}

func TaskList() error {
	projects := Projects{}
	if err := projects.load(); err != nil {
		return err
	}

	projectIndex := projects.getSelectedIndex()
	if projectIndex == -1 {
		fmt.Println("No project selected, please select a project first")
		return nil
	}

	fmt.Println("Tasks for project:", projects[projectIndex].Name)
	taskNumber := 1
	projects[projectIndex].Tasks.filterCompletedTasks()
	for _, task := range projects[projectIndex].Tasks {
		if task.ParentTaskID != "" {
			fmt.Printf("   - %s\n", task.Name)
			continue
		}
		fmt.Printf("%d. %s", taskNumber, task.Name)
		if task.IsJIRATask {
			fmt.Printf(" (%s)", task.ID)
		}
		fmt.Println()
		taskNumber++
	}

	return nil
}

func TaskRemove() error {
	projects := Projects{}
	if err := projects.load(); err != nil {
		return err
	}

	projectIndex := projects.getSelectedIndex()
	if projectIndex == -1 {
		fmt.Println("No project selected, please select a project first")
		return nil
	}

	projects[projectIndex].Tasks.filterCompletedTasks()
	prompt := promptui.Select{
		Label:     "Select a task to remove",
		Items:     projects[projectIndex].Tasks,
		Templates: taskSelectTemplate(),
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return err
	}

	removedTaskName := projects[projectIndex].Tasks[i].Name
	projects[projectIndex].Tasks.remove(i)
	if err := projects.save(); err != nil {
		return err
	}

	fmt.Println("Successfully removed task:", removedTaskName)

	return nil
}

func (t *Tasks) remove(index int) {
	*t = append((*t)[:index], (*t)[index+1:]...)
}

func taskSelectTemplate() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Active:   `{{ if .ParentTaskID }}  ▸ {{ .Name | cyan }}{{ else }}▸ {{ .Name | cyan }}{{ end }}{{ if .IsJIRATask }} ({{ .ID }}){{ end }}`,
		Inactive: `{{ if .ParentTaskID }}  ▸ {{ .Name }}{{ else }}▸ {{ .Name }}{{ end }}{{ if .IsJIRATask }} ({{ .ID }}){{ end }}`,
		Selected: `{{ "✔" | green }} {{ "Selected" | bold }}: {{ .Name | cyan }}{{ if .IsJIRATask }} ({{ .ID }}){{ end }}`,
		Details:  `{{ "Description:" }} {{ .Description }}`,
	}
}

func TaskStart() error {
	projects := Projects{}
	if err := projects.load(); err != nil {
		return err
	}

	projectIndex := projects.getSelectedIndex()
	if projectIndex == -1 {
		fmt.Println("No project selected, please select a project first")
		return nil
	}

	projects[projectIndex].Tasks.filterCompletedTasks()
	prompt := promptui.Select{
		Label:     "Select a task to start",
		Items:     projects[projectIndex].Tasks,
		Templates: taskSelectTemplate(),
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return err
	}

	projects[projectIndex].Tasks[i].IsStarted = true
	projects[projectIndex].Tasks[i].TaskHistories = append(projects[projectIndex].Tasks[i].TaskHistories, TaskHistory{
		StartedAt: time.Now(),
	})
	if err := projects.save(); err != nil {
		return err
	}

	fmt.Println("Started task:", projects[projectIndex].Tasks[i].Name)

	return nil
}

func TaskStop() error {
	projects := Projects{}
	if err := projects.load(); err != nil {
		return err
	}

	projectIndex := projects.getSelectedIndex()
	if projectIndex == -1 {
		fmt.Println("No project selected, please select a project first")
		return nil
	}

	i := projects[projectIndex].Tasks.getStartedTaskIndex()
	projects[projectIndex].Tasks[i].IsStarted = false
	projects[projectIndex].Tasks[i].TaskHistories[len(projects[projectIndex].Tasks[i].TaskHistories)-1].StoppedAt = time.Now()
	if err := projects.save(); err != nil {
		return err
	}

	fmt.Println("Stopped task:", projects[projectIndex].Tasks[i].Name)

	return nil
}
