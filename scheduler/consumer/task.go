package consumer

import (
	"attendance-api/scheduler"
	"fmt"
)

func InitTask(task *scheduler.AddTask) {
	if task.Action == "attendance" {
		TaskAttendance(task)
	}
}

func TaskAttendance(task *scheduler.AddTask) {
	fmt.Println("Execute Task Attendance")
	fmt.Printf("Action: %v\n", task.Action)
	fmt.Printf("Body  : %v\n", task.Body)
	fmt.Printf("Date  : %v\n", task.Date)
	fmt.Printf("TStm  : %v\n", task.TimeStamp)
}
