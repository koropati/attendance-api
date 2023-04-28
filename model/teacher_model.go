package model

type Teacher struct {
	GormCustom
	UserID         uint         `json:"user_id"`
	User           User         `json:"user"`
	NIP            string       `json:"nip" gorm:"type:varchar(100)"`
	DOB            string       `json:"dob" gorm:"type:date"`
	FacultyID      uint         `json:"faculty_id"`
	Faculty        Faculty      `json:"faculty"`
	MajorID        uint         `json:"major_id"`
	Major          Major        `json:"major"`
	StudyProgramID uint         `json:"study_program_id"`
	StudyProgram   StudyProgram `json:"study_program"`
	Address        string       `json:"address" gorm:"type:varchar(100)"`
	Gender         string       `json:"gender" gorm:"type:enum('laki-laki','perempuan');default:'laki-laki'"`
}
