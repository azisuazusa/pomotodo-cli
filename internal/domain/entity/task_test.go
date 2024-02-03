package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskTimeSpent(t *testing.T) {
	task := Task{
		Histories: []TaskHistory{
			{
				StartedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				StoppedAt: time.Date(2020, 1, 1, 0, 0, 10, 0, time.UTC),
			},
			{
				StartedAt: time.Date(2020, 1, 1, 0, 0, 10, 0, time.UTC),
				StoppedAt: time.Date(2020, 1, 1, 0, 0, 20, 0, time.UTC),
			},
		},
	}

	expected := "20s"

	actual := task.TimeSpent()

	assert.Equal(t, expected, actual.String())
}

func TestStop(t *testing.T) {
	task := Task{
		IsStarted: true,
		Histories: []TaskHistory{
			{
				StartedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	task.Stop()

	assert.Equal(t, false, task.IsStarted)
}

func TestStart(t *testing.T) {
	task := Task{
		IsStarted: false,
		Histories: []TaskHistory{},
	}

	task.Start()

	assert.Equal(t, true, task.IsStarted)
}

func TestComplete(t *testing.T) {
	task := Task{
		IsStarted: true,
		Histories: []TaskHistory{
			{
				StartedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	task.Complete()

	assert.Equal(t, false, task.IsStarted)
}
