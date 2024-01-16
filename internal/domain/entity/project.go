package entity

import "errors"

type Integration struct {
	IsEnabled bool
	Type      IntegrationType
	Details   map[string]string
}

type Project struct {
	ID           string
	Name         string
	Description  string
	IsSelected   bool
	Integrations []Integration
}

type Projects []Project

var ErrNoProjectSelected = errors.New("no project selected")
