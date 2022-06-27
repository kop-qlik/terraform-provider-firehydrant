package firehydrant

import (
	"context"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

type CreateTaskListRequest struct {
	Name string `json:"name"`
	Description string `json:"description"`

	ListItems []TaskListItems `json:"task_list_items"`
}

type TaskListItems struct {
	Summery string `json:"summery,omitempty"`
	Description string `json:"description"`
}

type UpdateTaskListRequest struct {
	Name string `json:"name"`
	Description string `json:"description"`
	ListItems []TaskListItems `json:"task_list_items,omitempty"`
}

type TaskListCreator struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Source string `json:"source"`
	Email string `json:"email"`
}

type TaskListResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy []TaskListCreator `json:"created_by"`
	ListItems []TaskListItems `json:"task_list_items"`
}

type UpdateTaskListResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy []TaskListCreator `json:"created_by"`
	ListItems []TaskListItems `json:"task_list_items"`
}

type TaskListsClient interface {
	Get(ctx context.Context, id string) (*TaskListResponse, error)
	Create(ctx context.Context, createReq CreateTaskListRequest) (*TaskListResponse, error)
	Update(ctx context.Context, id string, updateReq UpdateTaskListRequest) (*TaskListResponse, error)
	Delete(ctx context.Context, id string) error
}

type RESTTaskListsClient struct {
	client *APIClient
}

var _ TaskListsClient = &RESTTaskListsClient{}

func (c *RESTTaskListsClient) restClient() *sling.Sling {
	return c.client.client()
}

func (c *RESTTaskListsClient) Get(ctx context.Context, id string) (*TaskListResponse, error) {
	taskListReposnse := &TaskListResponse{}
	response, err := c.restClient().Get("task_lists"+id).Receive(taskListReposnse, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not get task list")
	}

	err = checkResponseStatusCode(response)
	if err != nil {
		return nil, err
	}
	
	return taskListReposnse, nil
}

func (c *RESTTaskListsClient) Create(ctx context.Context, createReq CreateTaskListRequest) (*TaskListResponse, error) {
	taskListResponse := &TaskListResponse{}
	response, err := c.restClient().Post("task_lists").BodyJSON(&createReq).Receive(taskListResponse, nil)
	if err != nil {
		return nil, errors.Wrap(err, "cloud not create task list")
	}

	err = checkResponseStatusCode(response)
	if err != nil {
		return nil, err
	}

	taskListResponse, err = c.Update(ctx, taskListResponse.ID, UpdateTaskListRequest {
		Name: createReq.Name,
		Description: createReq.Description,
		ListItems: createReq.ListItems,

	})
	if err != nil {
		return nil, errors.Wrap(err, "could not update created runbook")
	}
	if err != nil {
		return nil, err
	}

	return taskListResponse, nil
}

func (c *RESTTaskListsClient) Update(ctx context.Context, id string, updateReq UpdateTaskListRequest) (*TaskListResponse, error) {
	taskListResponse := &TaskListResponse{}
	response, err := c.restClient().Put("task_lists/"+id).BodyJSON(updateReq).Receive(taskListResponse, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not update ")
	}

	err = checkResponseStatusCode(response)
	if err != nil {
		return nil, err
	}

	return taskListResponse, nil
}

func (c *RESTTaskListsClient) Delete(ctx context.Context, id string) error {
	response, err := c.restClient().Delete("task_lists/"+id).Receive(nil, nil)
	if err != nil {
		return errors.Wrap(err, "could not delete task list")
	}

	err = checkResponseStatusCode(response)
	if err != nil {
		return err
	}

	return nil
}