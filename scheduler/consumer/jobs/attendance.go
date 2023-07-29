package jobs

import (
	"attendance-api/model"
	"attendance-api/scheduler"
	"attendance-api/service"
	"fmt"
	"log"
	"sync"
	"time"
)

type AttendanceJob interface {
	AutoCreate()
}

type attendanceJob struct {
	userScheduleService  service.UserScheduleService
	attendanceService    service.AttendanceService
	attendanceLogService service.AttendanceLogService
	task                 *scheduler.AddTask
}

func NewAttendanceJob(
	userScheduleService service.UserScheduleService,
	attendanceService service.AttendanceService,
	attendanceLogService service.AttendanceLogService,
	task *scheduler.AddTask,
) AttendanceJob {
	return &attendanceJob{
		userScheduleService:  userScheduleService,
		attendanceService:    attendanceService,
		attendanceLogService: attendanceLogService,
		task:                 task,
	}
}

func (j attendanceJob) AutoCreate() {
	fmt.Println("Execute Task Attendance")
	fmt.Printf("Action: %v\n", j.task.Action)
	fmt.Printf("Body  : %v\n", j.task.Body)
	fmt.Printf("Date  : %v\n", j.task.Date)
	fmt.Printf("TStm  : %v\n", j.task.TimeStamp)

	userSchedules, err := j.userScheduleService.GetAllByTodayRange()
	if err != nil {
		log.Printf("[Scheduler] [Error] [Attendance-GetAllByTodayRange] E: %v\n", err)
	}

	wg := sync.WaitGroup{}
	for _, userSchedule := range userSchedules {
		wg.Add(1)
		go func(userSchedule model.UserSchedule) {
			// Cek apakah jadwal user memang di hari ini
			log.Printf("IS Today schedule: %v\n", userSchedule.Schedule.IsTodaySchedule())
			if userSchedule.Schedule.IsTodaySchedule() {
				// Buat Data Presensi kosong / tidak hadir secara default terlebih dahulu
				dataAttendance := model.Attendance{
					UserID:         userSchedule.UserID,
					ScheduleID:     userSchedule.ScheduleID,
					Date:           time.Now().Format("2006-01-02"),
					ClockIn:        0,
					ClockOut:       0,
					Status:         "-",
					StatusPresence: "not_presence",
				}
				_, err := j.attendanceService.CreateAttendance(dataAttendance)
				if err != nil {
					log.Printf("[Scheduler] [Error] [Attendance-CreateAttendance] E: %v\n", err)
				}
			}
			wg.Done()
		}(userSchedule)
	}
	wg.Wait()
}
