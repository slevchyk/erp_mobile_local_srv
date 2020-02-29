package models

type Config struct {
	Auth       AuthConfig `json:"auth"`
	DB         DBConfig   `json:"db"`
	WinService winService `json:"win_service"`
}

type AuthConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type DBConfig struct {
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type winService struct {
	Name        string `json:"name"`
	LongName    string `json:"long_name"`
	Description string `json:"description"`
}

type Channel struct {
	ID       int    `json:"id"`
	UserID   string `json:"user_id"`
	UpdateID int    `json:"update_id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	News     string `json:"news"`
	Date     string `json:"date"`
}

type FirebaseTokens struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

type Timing struct {
	ID          int64              `json:"id"`
	MobID       int64              `json:"mob_id"`
	AccID       string             `json:"acc_id"`
	UserID      string             `json:"user_id"`
	Date        NullTime           `json:"date"`
	Status      string             `json:"status"`
	IsTurnstile ConvertibleBoolean `json:"is_turnstile"`
	StartedAt   NullTime           `json:"started_at"`
	EndedAt     NullTime           `json:"ended_at"`
	CreatedAt   NullTime           `json:"created_at"`
	UpdatedAt   NullTime           `json:"updated_at"`
	DeletedAt   NullTime           `json:"deleted_at"`
}

type Profile struct {
	ID                  int64              `json:"-"`
	Blocked             ConvertibleBoolean `json:"blocked"`
	UserID              string             `json:"user_id"`
	Pin                 string             `json:"pin"`
	InfoCard            int                `json:"info_card"`
	LastName            string             `json:"last_name"`
	FirstName           string             `json:"first_name"`
	MiddleName          string             `json:"middle_name"`
	ITN                 string             `json:"itn"`
	Phone               string             `json:"phone"`
	Birthday            NullTime           `json:"birthday"`
	Email               string             `json:"email"`
	Gender              string             `json:"gender"`
	Address             string             `json:"address"`
	PassportType        string             `json:"passport_type"`
	PassportSeries      string             `json:"passport_series"`
	PassportNumber      string             `json:"passport_number"`
	PassportIssued      string             `json:"passport_issued"`
	PassportDate        NullTime           `json:"passport_date"`
	PassportExpiry      NullTime           `json:"passport_expiry"`
	CivilStatus         string             `json:"civil_status"`
	Children            string             `json:"children"`
	JobPosition         string             `json:"job_position"`
	Education           int                `json:"education"`
	Specialty           string             `json:"specialty"`
	AdditionalEducation string             `json:"additional_education"`
	LastWorkPlace       string             `json:"last_work_place"`
	Skills              string             `json:"skills"`
	Languages           string             `json:"languages"`
	Disability          ConvertibleBoolean `json:"disability"`
	Pensioner           ConvertibleBoolean `json:"pensioner"`
	Photo               string             `json:"photo"`
	PhotoData           string             `json:"photo_data"`
}
