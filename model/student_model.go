package model

import "strings"

type Student struct {
	GormCustom
	UserID         uint         `json:"user_id"`
	User           User         `json:"user"`
	NIM            string       `json:"nim" gorm:"type:varchar(20);unique"`
	DOB            string       `json:"dob" gorm:"type:date"`
	FacultyID      uint         `json:"faculty_id"`
	Faculty        Faculty      `json:"faculty"`
	MajorID        uint         `json:"major_id"`
	Major          Major        `json:"major"`
	StudyProgramID uint         `json:"study_program_id"`
	StudyProgram   StudyProgram `json:"study_program"`
	Address        string       `json:"address" gorm:"type:varchar(255)"`
	Gender         string       `json:"gender" gorm:"type:enum('laki-laki','perempuan');default:'laki-laki'"`
}

type UserStudent struct {
	Username       string       `json:"username" gorm:"unique"`
	Password       string       `json:"password"`
	FirstName      string       `json:"first_name"`
	LastName       string       `json:"last_name"`
	Handphone      string       `json:"handphone" gorm:"unique"`
	Email          string       `json:"email" gorm:"unique"`
	Intro          string       `json:"intro" gorm:"type:varchar(255)"`
	Profile        string       `json:"profile" gorm:"type:varchar(255)"`
	UserID         uint         `json:"user_id"`
	NIM            string       `json:"nim" gorm:"type:varchar(20);unique"`
	DOB            string       `json:"dob" gorm:"type:date"`
	FacultyID      uint         `json:"faculty_id"`
	Faculty        Faculty      `json:"faculty"`
	MajorID        uint         `json:"major_id"`
	Major          Major        `json:"major"`
	StudyProgramID uint         `json:"study_program_id"`
	StudyProgram   StudyProgram `json:"study_program"`
	Address        string       `json:"address" gorm:"type:varchar(255)"`
	Gender         string       `json:"gender" gorm:"type:enum('laki-laki','perempuan');default:'laki-laki'"`
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
