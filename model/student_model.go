package model

import "strings"

type Student struct {
	GormCustom
	UserID         uint         `json:"user_id" query:"user_id" form:"user_id"`
	User           User         `json:"user" query:"user" form:"user"`
	NIM            string       `json:"nim" gorm:"type:varchar(20);unique" query:"nim" form:"nim"`
	DOB            string       `json:"dob" gorm:"type:date" query:"dob" form:"dob"`
	FacultyID      uint         `json:"faculty_id" query:"faculty_id" form:"faculty_id"`
	Faculty        Faculty      `json:"faculty" query:"faculty" form:"faculty"`
	MajorID        uint         `json:"major_id" query:"major_id" form:"major_id"`
	Major          Major        `json:"major" query:"major" form:"major"`
	StudyProgramID uint         `json:"study_program_id" query:"study_program_id" form:"study_program_id"`
	StudyProgram   StudyProgram `json:"study_program" query:"study_program" form:"study_program"`
	Address        string       `json:"address" gorm:"type:varchar(255)" query:"address" form:"address"`
	Gender         string       `json:"gender" gorm:"type:enum('laki-laki','perempuan');default:'laki-laki'" query:"gender" form:"gender"`
	Avatar         string       `json:"avatar" gorm:"-" query:"avatar" form:"avatar"`
	ScheduleID     int          `json:"schedule_id" gorm:"-" query:"schedule_id" form:"schedule_id"`
	OwnerID        int          `json:"owner_id" gorm:"-" query:"owner_id" form:"owner_id"`
}

type UserStudent struct {
	Username       string       `json:"username" gorm:"unique" query:"username" form:"username"`
	Password       string       `json:"password" query:"password" form:"password"`
	FirstName      string       `json:"first_name" query:"first_name" form:"first_name"`
	LastName       string       `json:"last_name" query:"last_name" form:"last_name"`
	Handphone      string       `json:"handphone" gorm:"unique" query:"handphone" form:"handphone"`
	Email          string       `json:"email" gorm:"unique" query:"email" form:"email"`
	Intro          string       `json:"intro" gorm:"type:varchar(255)" query:"intro" form:"intro"`
	Profile        string       `json:"profile" gorm:"type:varchar(255)" query:"profile" form:"profile"`
	UserID         uint         `json:"user_id" query:"user_id" form:"user_id"`
	NIM            string       `json:"nim" gorm:"type:varchar(20);unique" query:"nim" form:"nim"`
	DOB            string       `json:"dob" gorm:"type:date" query:"dob" form:"dob"`
	FacultyID      uint         `json:"faculty_id" query:"faculty_id" form:"faculty_id"`
	Faculty        Faculty      `json:"faculty" query:"faculty" form:"faculty"`
	MajorID        uint         `json:"major_id" query:"major_id" form:"major_id"`
	Major          Major        `json:"major" query:"major" form:"major"`
	StudyProgramID uint         `json:"study_program_id" query:"study_program_id" form:"study_program_id"`
	StudyProgram   StudyProgram `json:"study_program" query:"study_program" form:"study_program"`
	Address        string       `json:"address" gorm:"type:varchar(255)" query:"address" form:"address"`
	Gender         string       `json:"gender" gorm:"type:enum('laki-laki','perempuan');default:'laki-laki'" query:"gender" form:"gender"`
	Avatar         string       `json:"avatar" gorm:"-" query:"avatar" form:"avatar"`
}

func (data Student) GetAvatar() (url string) {
	if data.Gender == "laki-laki" {
		return "https://cdn-icons-png.flaticon.com/512/8348/8348118.png"
	} else {
		return "https://cdn-icons-png.flaticon.com/512/8348/8348099.png"
	}
}

func (data Student) GeneratePassword() (passwordGenrate string) {
	firstName := strings.ToLower(data.User.FirstName)
	firstName = strings.Replace(firstName, " ", "_", -1)
	dob := strings.Replace(data.DOB, "-", "", -1)
	passwordGenrate = firstName + "@" + dob
	// dewa_ketut@19970621
	return
}

func (data UserStudent) GetUser() (user User) {
	return User{
		Username:  data.Username,
		Password:  data.Password,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Handphone: data.Handphone,
		Email:     data.Email,
		Intro:     data.Intro,
		Profile:   data.Profile,
	}
}

func (data UserStudent) GetStudent() (student Student) {
	return Student{
		UserID:         data.UserID,
		NIM:            data.NIM,
		DOB:            data.DOB,
		FacultyID:      data.FacultyID,
		Faculty:        data.Faculty,
		MajorID:        data.MajorID,
		Major:          data.Major,
		StudyProgramID: data.StudyProgramID,
		Address:        data.Address,
		Gender:         data.Gender,
	}
}
