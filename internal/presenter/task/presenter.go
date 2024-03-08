package task

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	"github.com/azisuazusa/todo-cli/internal/domain/jira"
	"github.com/azisuazusa/todo-cli/internal/domain/syncintegration"
	"github.com/azisuazusa/todo-cli/internal/domain/task"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/manifoldco/promptui"
)

type Presenter struct {
	taskUseCase    task.UseCase
	settingUseCase syncintegration.UseCase
	jiraUseCase    jira.UseCase
}

func New(taskUseCase task.UseCase, settingUseCase syncintegration.UseCase, jiraUseCase jira.UseCase) *Presenter {
	return &Presenter{
		taskUseCase:    taskUseCase,
		settingUseCase: settingUseCase,
		jiraUseCase:    jiraUseCase,
	}
}

func (p *Presenter) GetUncompleteTasks(ctx context.Context) error {
	tasks, err := p.taskUseCase.GetUncompleteTasks(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Todo List:")
	t.AppendHeader(table.Row{"#", "Name", ""})

	taskNumber := 1
	for _, task := range tasks {
		isStarted := ""
		if task.IsStarted {
			isStarted = "Started"
		}

		if task.ParentTaskID == "" {
			t.AppendRow(table.Row{taskNumber, task.Name, isStarted})
			taskNumber++
			continue
		}

		t.AppendRow(table.Row{"-", task.Name, isStarted})
	}

	t.SetStyle(table.StyleLight)
	t.Render()

	return nil
}

func (p *Presenter) Add(ctx context.Context) error {
	prompt := promptui.Prompt{
		Label: "Name",
	}

	name, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	prompt = promptui.Prompt{
		Label: "Description",
	}

	description, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	task := entity.Task{
		Name:        name,
		Description: description,
	}

	if err = p.taskUseCase.Add(ctx, task); err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	if err := p.settingUseCase.Upload(ctx); err != nil {
		fmt.Printf("Upload error: %v\n", err)
		return err
	}

	fmt.Println("Task added successfully")

	return nil
}

func (p *Presenter) AddSubTask(ctx context.Context) error {
	prompt := promptui.Prompt{
		Label: "Name",
	}

	name, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	prompt = promptui.Prompt{
		Label: "Description",
	}

	description, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	parentTasks, err := p.taskUseCase.GetUncompleteTasks(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	views := []TaskView{}
	number := 1
	for _, t := range parentTasks {
		if t.ParentTaskID == "" {
			views = append(views, CreateTaskView(fmt.Sprintf("%d.", number), t))
			number++
			continue
		}

		views = append(views, CreateTaskView("-", t))
	}

	selectPrompt := promptui.Select{
		Label:     "Select Parent Task",
		Items:     views,
		Templates: taskSelectTemplate,
		Size:      10,
	}

	parentTaskIndex, _, err := selectPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	task := entity.Task{
		Name:         name,
		Description:  description,
		ParentTaskID: parentTasks[parentTaskIndex].ID,
	}

	if err = p.taskUseCase.Add(ctx, task); err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	if err := p.settingUseCase.Upload(ctx); err != nil {
		fmt.Printf("Upload error: %v\n", err)
		return err
	}

	fmt.Println("Subtask added successfully")

	return nil
}

func (p *Presenter) Start(ctx context.Context) error {
	tasks, err := p.taskUseCase.GetUncompleteTasks(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	views := []TaskView{}
	number := 1
	for _, task := range tasks {
		if task.ParentTaskID == "" {
			views = append(views, CreateTaskView(fmt.Sprintf("%d.", number), task))
			number++
			continue
		}

		views = append(views, CreateTaskView("-", task))
	}

	selectPrompt := promptui.Select{
		Label:     "Select Task",
		Items:     views,
		Templates: taskSelectTemplate,
		Size:      10,
	}

	taskIndex, _, err := selectPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	// Makesure that no other task is running
	if err = p.taskUseCase.Stop(ctx); err != nil && errors.Unwrap(err) != task.ErrTaskNotFound {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	if err = p.taskUseCase.Start(ctx, tasks[taskIndex].ID); err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	if err := p.settingUseCase.Upload(ctx); err != nil {
		fmt.Printf("Upload error: %v\n", err)
		return err
	}

	fmt.Println("Task started successfully")

	return nil
}

func (p *Presenter) Complete(ctx context.Context) error {
	task, err := p.taskUseCase.GetUncompleteTasks(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	views := []TaskView{}
	number := 1
	for _, task := range task {
		if task.ParentTaskID == "" {
			views = append(views, CreateTaskView(fmt.Sprintf("%d.", number), task))
			number++
			continue
		}

		views = append(views, CreateTaskView("-", task))
	}

	selectPrompt := promptui.Select{
		Label:     "Select Task",
		Items:     views,
		Templates: taskSelectTemplate,
		Size:      10,
	}

	taskIndex, _, err := selectPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	if err = p.taskUseCase.Complete(ctx, task[taskIndex].ID); err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	if err := p.settingUseCase.Upload(ctx); err != nil {
		fmt.Printf("Upload error: %v\n", err)
		return err
	}

	currentTask := task[taskIndex]
	integrationType := currentTask.Integration.Type
	if currentTask.ParentTaskID != "" {
		parentTask, err := p.taskUseCase.GetByID(ctx, currentTask.ParentTaskID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		integrationType = parentTask.Integration.Type
	}

	if integrationType == entity.IntegrationTypeJIRA {
		if err := p.addWorklogJIRA(ctx, task[taskIndex].ID); err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
	}

	fmt.Println("Task completed successfully")

	return nil
}

func (p *Presenter) addWorklogJIRA(ctx context.Context, taskID string) error {
	task, err := p.taskUseCase.GetByID(ctx, taskID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	prompt := promptui.Prompt{
		Label:   "Add Worklog",
		Default: task.TimeSpent().Round(time.Minute).String(),
	}

	timeSpentStr, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	timeSpent, err := time.ParseDuration(timeSpentStr)
	if err != nil {
		fmt.Printf("Parse duration failed %v\n", err)
		return err
	}

	if err := p.jiraUseCase.AddWorklog(ctx, task, timeSpent); err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	return nil
}

func (p *Presenter) Remove(ctx context.Context) error {
	tasks, err := p.taskUseCase.GetUncompleteTasks(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	views := []TaskView{}
	number := 1
	for _, task := range tasks {
		if task.ParentTaskID == "" {
			views = append(views, CreateTaskView(fmt.Sprintf("%d.", number), task))
			number++
			continue
		}

		views = append(views, CreateTaskView("-", task))
	}

	selectPrompt := promptui.Select{
		Label:     "Select Task",
		Items:     views,
		Templates: taskSelectTemplate,
		Size:      10,
	}

	taskIndex, _, err := selectPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	if err = p.taskUseCase.Remove(ctx, tasks[taskIndex].ID); err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	if err := p.settingUseCase.Upload(ctx); err != nil {
		fmt.Printf("Upload error: %v\n", err)
		return err
	}

	fmt.Println("Task removed successfully")

	return nil
}

func (p *Presenter) Stop(ctx context.Context) error {
	err := p.taskUseCase.Stop(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	if err := p.settingUseCase.Upload(ctx); err != nil {
		fmt.Printf("Upload error: %v\n", err)
		return err
	}

	fmt.Println("Task stopped successfully")

	return nil
}
