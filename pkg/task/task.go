package task

import (
	log "github.com/sirupsen/logrus"
)

// Task defines task to be performed in a job
type Task interface {
	Run() error
}

type SeqTask struct {
	Tasks []Task
	log   log.FieldLogger
}

func NewSeqTask(log log.FieldLogger, tasks ...Task) *SeqTask {
	return &SeqTask{
		Tasks: tasks,
		log:   log,
	}
}

func (st *SeqTask) Run() error {
	for _, task := range st.Tasks {
		st.log.Infof("Running sequential task %v", task)
		if err := task.Run(); err != nil {
			st.log.Errorf("error running sequential task %v", task, err)
			return err
		}
	}
	return nil
}
