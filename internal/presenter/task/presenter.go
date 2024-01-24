package task

import (
	"context"
	"fmt"
	"os"

	"github.com/azisuazusa/todo-cli/internal/domain/entity"
	"github.com/azisuazusa/todo-cli/internal/domain/syncintegration"
	"github.com/azisuazusa/todo-cli/internal/domain/task"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/manifoldco/promptui"
)

type Presenter struct {
	taskUseCase    task.UseCase
	settingUseCase syncintegration.UseCase
}

func New(taskUseCase task.UseCase, settingUseCase syncintegration.UseCase) *Presenter {
	return &Presenter{
		taskUseCase:    taskUseCase,
		settingUseCase: settingUseCase,
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
	t.AppendHeader(table.Row{"#", "Name", "Description", ""})

	taskNumber := 1
	for _, task := range tasks {
		isStarted := ""
		if task.IsStarted {
			isStarted = "Started"
		}

		description := task.Description
		if len(description) > 50 {
			description = description[:50] + "..."
		}

		if task.ParentTaskID == "" {
			t.AppendRow(table.Row{taskNumber, task.Name, description, isStarted})
			taskNumber++
			continue
		}

		t.AppendRow(table.Row{"", task.Name, description, isStarted})
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
	for i, t := range parentTasks {
		views = append(views, CreateTaskView(i+1, t))
	}

	selectPrompt := promptui.Select{
		Label:     "Select Parent Task",
		Items:     views,
		Templates: taskSelectTemplate,
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
	for i, task := range tasks {
		views = append(views, CreateTaskView(i+1, task))
	}

	selectPrompt := promptui.Select{
		Label:     "Select Task",
		Items:     views,
		Templates: taskSelectTemplate,
	}

	taskIndex, _, err := selectPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
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
	for i, task := range task {
		views = append(views, CreateTaskView(i+1, task))
	}

	selectPrompt := promptui.Select{
		Label:     "Select Task",
		Items:     views,
		Templates: taskSelectTemplate,
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

	fmt.Println("Task completed successfully")

	return nil
}

func (p *Presenter) Remove(ctx context.Context) error {
	tasks, err := p.taskUseCase.GetUncompleteTasks(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	views := []TaskView{}
	for i, task := range tasks {
		views = append(views, CreateTaskView(i+1, task))
	}

	selectPrompt := promptui.Select{
		Label:     "Select Task",
		Items:     views,
		Templates: taskSelectTemplate,
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
	if err := p.taskUseCase.Stop(ctx); err != nil {
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
