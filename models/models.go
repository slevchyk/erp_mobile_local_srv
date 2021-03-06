package models

type Config struct {
	Auth       AuthConfig `json:"auth"`
	MainSrv    string     `json:"main_srv"`
	MainAuth   AuthConfig `json:"main_auth"`
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
	PhotoName           string             `json:"photo_name"`
	PhotoData           string             `json:"photo_data"`
}

type HelpDesk struct {
	ID              int                `json:"id"`
	UserID          string             `json:"user_id"`
	Date            NullTime           `json:"date"`
	Title           string             `json:"title"`
	Body            string             `json:"body"`
	Status          string             `json:"status"`
	Answer          string             `json:"answer"`
	AnsweredBy      string             `json:"answered_by"`
	AnsweredAt      NullTime           `json:"answered_at"`
	IsModifiedByMob ConvertibleBoolean `json:"-"`
	IsModifiedByAcc ConvertibleBoolean `json:"-"`
}

type PayDesk struct {
	ID                 int                `json:"id"`
	PayDeskType        int                `json:"pay_desk_type"`
	UserID             string             `json:"user_id"`
	CurrencyAccID      string             `json:"currency_acc_id"`
	CostItemAccID      string             `json:"cost_item_acc_id"`
	IncomeItemAccID    string             `json:"income_item_acc_id"`
	FromPayOfficeAccID string             `json:"from_pay_office_acc_id"`
	ToPayOfficeAccID   string             `json:"to_pay_office_acc_id"`
	Amount             float32            `json:"amount"`
	Payment            string             `json:"payment"`
	DocumentNumber     string             `json:"document_number"`
	DocumentDate       NullTime           `json:"document_date"`
	FilePaths          string             `json:"file_paths"`
	FilesQuantity      int                `json:"files_quantity"`
	IsChecked          ConvertibleBoolean `json:"is_checked"`
	IsReadOnly         ConvertibleBoolean `json:"is_read_only"`
	CreatedAt          NullTime           `json:"created_at"`
	UpdatedAt          NullTime           `json:"updated_at"`
	IsDeleted          ConvertibleBoolean `json:"is_deleted"`
	IsModifiedByMob    ConvertibleBoolean `json:"-"`
	IsModifiedByAcc    ConvertibleBoolean `json:"-"`
}

type CostItem struct {
	ID        int                `json:"id"`
	AccID     string             `json:"acc_id"`
	Name      string             `json:"name"`
	CreatedAt NullTime           `json:"created_at"`
	UpdatedAt NullTime           `json:"updated_at"`
	IsDeleted ConvertibleBoolean `json:"is_deleted"`
}

type IncomeItem struct {
	ID        int                `json:"id"`
	AccID     string             `json:"acc_id"`
	Name      string             `json:"name"`
	CreatedAt NullTime           `json:"created_at"`
	UpdatedAt NullTime           `json:"updated_at"`
	IsDeleted ConvertibleBoolean `json:"is_deleted"`
}

type PayOffice struct {
	ID            int                `json:"id"`
	AccID         string             `json:"acc_id"`
	Name          string             `json:"name"`
	CurrencyAccID string             `json:"currency_acc_id"`
	CreatedAt     NullTime           `json:"created_at"`
	UpdatedAt     NullTime           `json:"updated_at"`
	IsDeleted     ConvertibleBoolean `json:"is_deleted"`
}

type PayOfficeBalance struct {
	AccID     string   `json:"acc_id"`
	Balance   float32  `json:"balance"`
	UpdatedAt NullTime `json:"updated_at"`
}

type Currency struct {
	ID        int                `json:"id"`
	AccID     string             `json:"acc_id"`
	Code      int                `json:"code"`
	Name      string             `json:"name"`
	CreatedAt NullTime           `json:"created_at"`
	UpdatedAt NullTime           `json:"updated_at"`
	IsDeleted ConvertibleBoolean `json:"is_deleted"`
}

type UserGrants struct {
	UserID      string             `json:"user_id"`
	ObjectType  int                `json:"odject_type"`
	ObjectAccID string             `json:"odject_acc_id"`
	IsVisible   ConvertibleBoolean `json:"is_visible"`
	IsAvailable ConvertibleBoolean `json:"is_available"`
	IsReceiver  ConvertibleBoolean `json:"is_receiver"`
}

type WebServerLog struct {
	Time      string `json:"time"`
	ID        string `json:"id"`
	MobID     string `json:"mob_id"`
	AccID     string `json:"acc_id"`
	IsError   bool   `json:"is_error"`
	IsWarning bool   `json:"is_warning"`
	Message   string `json:"message"`
}

type PayDeskImage struct {
	PID       int                `json:"pid"`
	ImageName string             `json:"image_name"`
	File      string             `json:"file"`
	Sha256    string             `json:"sha256"`
	IsDeleted ConvertibleBoolean `json:"is_deleted"`
}

type PayDeskImageSha256 struct {
	ImageName string `json:"image_name"`
	Sha256    string `json:"sha256"`
}

type LogInfo struct {
	UserID    string    `json:"user_id"`
	FileName  string    `json:"file_name"`
	File      string    `json:"file"`
	Date      NullTime  `json:"date"`
}