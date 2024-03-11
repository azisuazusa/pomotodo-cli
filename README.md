# Todo-CLI
![Screenshot 2024-03-11 23:59:29](https://github.com/azisuazusa/todo-cli/assets/18085025/085cb209-ae83-43df-9dcc-91848102e0a3)

Todo-CLI is a command-line task management application designed to streamline your productivity and task management. Integrating seamlessly with tools like JIRA and Dropbox, it offers a robust solution for tracking your todos directly from your terminal.

## Features
- CLI-Based Management: Easily add, remove, and update tasks through straightforward commands.
- JIRA Integration: Synchronize your tasks with JIRA to keep all your project management in one place.
- Dropbox Sync: Backup and sync your tasks across devices using Dropbox.

## Installation
### Install from source
```
git clone git@github.com:azisuazusa/todo-cli.git
cd todo-cli
go build -o todo
todo setup
```

## Integrations
### JIRA Integration
```
todo project add
todo project add-integration
```

### Dropbox Integration
```
todo setting sync-integration
```

## TODOs
- [ ] Integrate with GitHub Issue for task synchronization
- [ ] Integrate with Slack for update status whenever with start working on a task
