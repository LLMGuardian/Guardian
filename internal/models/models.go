package models

import (
	"github.com/google/uuid"
)

// SendRequest represents a request to send a prompt.
type SendRequest struct {
	UserID uuid.UUID     `json:"user_id"`
	ChatID *uuid.UUID    `json:"chat_id,omitempty"`
	Prompt string        `json:"prompt"`
	Target []TargetModel `json:"targets"`
}

// SendResponse represents the response from a send operation.
type SendResponse struct {
	Status string  `json:"status"`
	Target AIModel `json:"target,omitempty"`
}

// Group represents a group of users.
type Group struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Status int       `json:"status"`
	Users  []User    `json:"users"` // Handling user associations differently
}

// User represents a user of the system.
type User struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Status int       `json:"status"`
	Groups []Group   `json:"groups"` // Same as above
}

// AIModel represents an AI model's metadata.
type AIModel struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Status  int    `json:"status"`
}

// RefereeModel represents a model used by referees.
type RefereeModel struct {
	ModelID uuid.UUID `json:"model_id"`
	Token   string    `json:"token,omitempty"`
}

// TargetModel represents the target model for processing.
type TargetModel struct {
	ModelID uuid.UUID `json:"model_id"`
	Token   string    `json:"token"`
}

// Usage records token consumption for users.
type Usage struct {
	UserID                 uuid.UUID `json:"user_id"`
	TargetModelID          uuid.UUID `json:"target_model_id"`
	InputTokenConsumption  int       `json:"input_token_consumption"`
	OutputTokenConsumption int       `json:"output_token_consumption"`
}

// Task represents a task that can be used in the pipeline.
type Task struct {
	Type string `json:"type"`
}

// Pipeline contains the related processing tasks for each
type Pipeline struct {
	UserTasks  *[]UserTask  `json:"user_tasks,omitempty"`
	GroupTasks *[]GroupTask `json:"group_tasks,omitempty"`
}

// UserTask links a user to a task.
type UserTask struct {
	UserID uuid.UUID `json:"user_id"`
	Task   Task      `json:"task"`
}

// GroupTask links a group to a task.
type GroupTask struct {
	GroupID uuid.UUID `json:"group_id"`
	Task    Task      `json:"task"`
}