package model

type Auth struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	UserID   uint   `json:"user_id" gorm:"not null"`
	AuthUUID string `json:"auth_uuid" gorm:"size:255;not null;"`
	Expired  int64  `json:"expired"`
	TypeAuth string `json:"type_auth" gorm:"type:enum('at','rt');default:'at'"`
}

type Ability struct {
	Action  string `json:"action"`
	Subject string `json:"subject"`
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
