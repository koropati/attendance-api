package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type DashboardService interface {
	RetrieveDashboardAcademic() (result model.DashboardAcademic, err []error)
	RetrieveDashboardUser() (result model.DashboardUser, err error)
	RetrieveDashboardStudent() (result model.DashboardStudent, err error)
	RetrieveDashboardTeacher() (result model.DashboardTeacher, err error)
	RetrieveDashboardAttendance(month, year int) (results []model.DashboardAttendance, err error)
	RetrieveDashboardAttendanceSeries(month, year int) (results []model.AttendanceSeries, err error)
}

type dashboardService struct {
	dashboardRepo repo.DashboardRepo
}

func NewDashboardService(dashboardRepo repo.DashboardRepo) DashboardService {
	return &dashboardService{dashboardRepo: dashboardRepo}
}

func (s dashboardService) RetrieveDashboardAcademic() (result model.DashboardAcademic, err []error) {
	data, err := s.dashboardRepo.RetrieveDashboardAcademic()
	if err != nil {
		return model.DashboardAcademic{}, err
	}
	return data, nil
}

func (s dashboardService) RetrieveDashboardUser() (result model.DashboardUser, err error) {
	data, err := s.dashboardRepo.RetrieveDashboardUser()
	if err != nil {
		return model.DashboardUser{}, err
	}
	return data, nil
}

func (s dashboardService) RetrieveDashboardStudent() (result model.DashboardStudent, err error) {
	data, err := s.dashboardRepo.RetrieveDashboardStudent()
	if err != nil {
		return model.DashboardStudent{}, err
	}
	return data, nil
}

func (s dashboardService) RetrieveDashboardTeacher() (result model.DashboardTeacher, err error) {
	data, err := s.dashboardRepo.RetrieveDashboardTeacher()
	if err != nil {
		return model.DashboardTeacher{}, err
	}
	return data, nil
}

func (s dashboardService) RetrieveDashboardAttendance(month, year int) (results []model.DashboardAttendance, err error) {
	data, err := s.dashboardRepo.RetrieveDashboardAttendance(month, year)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s dashboardService) RetrieveDashboardAttendanceSeries(month, year int) (results []model.AttendanceSeries, err error) {
	datas, err := s.dashboardRepo.RetrieveDashboardAttendanceSeries(month, year)
	if err != nil {
		return nil, err
	}
	return datas, nil
}
