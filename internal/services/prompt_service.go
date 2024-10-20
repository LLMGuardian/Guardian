package services

import (
	"errors"
	"fmt"
	"guardian/configs"
	"guardian/internal/models"
	"guardian/utlis/logger"
	"sync"
)

type PromptServiceInterface interface {
	ProcessPrompt(req models.SendRequest) (string, error)
}

type PromptService struct {
	userTaskService UserTaskServiceInterface
	tasksMap        map[string]ProcessingTask
}

func NewPromptService(userTaskService UserTaskServiceInterface) *PromptService {
	return &PromptService{
		userTaskService: userTaskService,
		tasksMap: map[string]ProcessingTask{
			"external-api": &ExternalHttpServiceTask{ApiUrl: "https://google.com"},
		},
	}
}

func (p *PromptService) ProcessPrompt(req models.SendRequest) (string, error) {
	if req.Prompt == "" {
		return "", errors.New("empty prompt")
	}

	if p.pipeline(req) {
		return "malicious", nil
	}

	return "benign", nil
}

func (p *PromptService) pipeline(req models.SendRequest) bool {
	userTasks, err := p.userTaskService.GetUserTasks(req.UserID)
	if err != nil {
		logger.GetLogger().Error(err)
		return false
	}

	workerPoolSize := configs.LoadConfig().PipelineWorkerPoolSize

	taskChan := make(chan models.UserTask, len(userTasks))
	resultsChan := make(chan models.TaskResult, len(userTasks))
	quit := make(chan struct{})

	var wg sync.WaitGroup
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go p.worker(taskChan, resultsChan, quit, req, &wg)
	}

	for _, userTask := range userTasks {
		taskChan <- userTask
	}
	close(taskChan)

	wg.Wait()
	close(resultsChan)

	for result := range resultsChan {
		if result.Err != nil {
			logger.GetLogger().Errorf("task %s faced error", result.Err)
		}

		if !result.Success {
			logger.GetLogger().Infof("task %s failed:", result.TaskType)
			return false
		}
	}

	return true
}

func (p *PromptService) worker(taskChan chan models.UserTask, resultsChan chan models.TaskResult, quit chan struct{},
	req models.SendRequest, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case userTask, ok := <-taskChan:
			if !ok {
				return
			}
			taskType := userTask.Task.Type

			if task, exists := p.tasksMap[taskType]; exists {
				result, err := task.Process(req)
				if err != nil {
					resultsChan <- models.TaskResult{TaskType: taskType, Success: false, Err: err}
					close(quit)
					return
				}

				if !result {
					resultsChan <- models.TaskResult{TaskType: taskType, Success: false}
					close(quit)
					return
				}

				resultsChan <- models.TaskResult{TaskType: taskType, Success: true}
			} else {
				err := fmt.Errorf("task type %s not found", taskType)
				resultsChan <- models.TaskResult{TaskType: taskType, Success: false, Err: err}
			}

		case <-quit:
			return
		}
	}
}
