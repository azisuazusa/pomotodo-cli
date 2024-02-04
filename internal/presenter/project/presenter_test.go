package project

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	projectMocks "github.com/azisuazusa/todo-cli/internal/domain/project/mocks"
	synintegrationMocks "github.com/azisuazusa/todo-cli/internal/domain/syncintegration/mocks"
	taskMocks "github.com/azisuazusa/todo-cli/internal/domain/task/mocks"
	"github.com/stretchr/testify/suite"
)

type PresenterTestSuite struct {
	suite.Suite
	projectUseCase         *projectMocks.UseCase
	syncintegrationUseCase *synintegrationMocks.UseCase
	taskUseCase            *taskMocks.UseCase
	presenter              *Presenter
}

func (t *PresenterTestSuite) SetupTest() {
	t.projectUseCase = new(projectMocks.UseCase)
	t.syncintegrationUseCase = new(synintegrationMocks.UseCase)
	t.taskUseCase = new(taskMocks.UseCase)
	t.presenter = New(t.projectUseCase, t.syncintegrationUseCase, t.taskUseCase)
}

func TestPresenterTestSuite(t *testing.T) {
	suite.Run(t, new(PresenterTestSuite))
}

func (t *PresenterTestSuite) TestGetProjects() {
	tests := []struct {
		name        string
		expectedRes string
		expectedErr error
		mockFunc    func()
	}{
		{
			name:        "failed to get projects",
			expectedRes: "any-error\n",
			expectedErr: errors.New("any-error"),
			mockFunc: func() {
				t.projectUseCase.On("GetAll", context.Background()).Return(nil, errors.New("any-error")).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func() {
			test.mockFunc()

			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() {
				w.Close()
				os.Stdout = originalStdout
			}()

			err := t.presenter.GetProjects(context.Background())

			// Capture the output
			output, _ := io.ReadAll(r)
			t.Equal(test.expectedRes, string(output))
			t.Equal(test.expectedErr, err)
		})
	}
}
