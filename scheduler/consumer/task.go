package consumer

import (
	"attendance-api/infra"
	"attendance-api/manager"
	"attendance-api/scheduler"
	"attendance-api/scheduler/consumer/jobs"
)

type Task interface {
	InitTask(task *scheduler.AddTask)
}

type task struct {
	infra   infra.Infra
	service manager.ServiceManager
}

func NewTask(infra infra.Infra) Task {
	return &task{
		infra:   infra,
		service: manager.NewServiceManager(infra),
	}
}

func (t task) InitTask(task *scheduler.AddTask) {
	attendanceJob := jobs.NewAttendanceJob(
		t.service.UserScheduleService(),
		t.service.AttendanceService(),
		t.service.AttendanceLogService(),
		task,
	)

	if task.Action == "attendance" {
		attendanceJob.AutoCreate()
	}
}
