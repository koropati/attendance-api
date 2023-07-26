package model

type Auth struct {
	ID       uint   `json:"id" gorm:"primary_key" query:"id" form:"id"`
	UserID   uint   `json:"user_id" gorm:"not null" query:"user_id" form:"user_id"`
	AuthUUID string `json:"auth_uuid" gorm:"size:255;not null;" query:"auth_uuid" form:"auth_uuid"`
	Expired  int64  `json:"expired" query:"expired" form:"expired"`
	TypeAuth string `json:"type_auth" gorm:"type:enum('at','rt');default:'at'" query:"type_auth" form:"type_auth"`
}

type Ability struct {
	Action  string `json:"action"`
	Subject string `json:"subject"`
}

type ForgotPassword struct {
	Email string `json:"email" query:"email"`
}

type ConfirmForgotPassword struct {
	Token           string `json:"token" query:"token" form:"token"`
	Password        string `json:"password" query:"password" form:"password"`
	ConfirmPassword string `json:"confirm_password" query:"confirm_password" form:"confirm_password"`
}

func GetDefaultAbility() (results []Ability) {
	results = append(results, Ability{
		Action:  "auth",
		Subject: "read",
	})
	return
}

func GetSuperAdminAbility() (results []Ability) {
	results = append(results, Ability{
		Action:  "manage",
		Subject: "all",
	})
	return
}

func GetAdminAbility() (results []Ability) {
	menus := []string{"student", "teacher", "attendance", "schedule", "daily_schedule", "subject", "user_schedule"}
	actions := []string{"read", "read", "manage", "manage", "manage", "manage", "manage"}
	for i, menu := range menus {
		results = append(results, Ability{
			Action:  actions[i],
			Subject: menu,
		})
	}
	return
}

func GetUserAbility() (results []Ability) {
	menus := []string{"student", "teacher", "attendance", "schedule", "daily_schedule", "subject", "user_schedule"}
	actions := []string{"read", "read", "create", "read", "read", "read", "read"}
	for i, menu := range menus {
		results = append(results, Ability{
			Action:  actions[i],
			Subject: menu,
		})
	}
	return
}
