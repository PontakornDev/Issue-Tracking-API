package entities

import "time"

// User represents a user who can report issues
type User struct {
	UserID    uint       `gorm:"primaryKey;column:user_id;autoIncrement" json:"user_id"`
	FullName  string     `gorm:"column:full_name;not null" json:"full_name" validate:"required,min=2,max=255"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	// Relations
	Issues   []Issue   `gorm:"foreignKey:ReporterID;references:UserID" json:"issues,omitempty"`
	Comments []Comment `gorm:"foreignKey:UserID;references:UserID" json:"comments,omitempty"`
}

func (User) TableName() string {
	return "users"
}

// Officer represents an officer who can handle issues
type Officer struct {
	OfficerID uint       `gorm:"primaryKey;column:officer_id;autoIncrement" json:"officer_id"`
	FullName  string     `gorm:"column:full_name;not null" json:"full_name" validate:"required,min=2,max=255"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	// Relations
	AssignedIssues  []Issue              `gorm:"foreignKey:AssigneeID;references:OfficerID" json:"assigned_issues,omitempty"`
	StatusHistories []IssueStatusHistory `gorm:"foreignKey:ChangedBy;references:OfficerID" json:"status_histories,omitempty"`
}

func (Officer) TableName() string {
	return "officer"
}

// IssueStatus represents the status of an issue
type IssueStatus struct {
	StatusID     uint      `gorm:"primaryKey;column:status_id;autoIncrement" json:"status_id"`
	StatusCode   string    `gorm:"column:status_code;type:varchar(50);unique;not null;index" json:"status_code" validate:"required,min=2,max=50"`
	DisplayName  string    `gorm:"column:display_name;type:varchar(100);not null" json:"display_name" validate:"required,min=2,max=100"`
	Description  string    `gorm:"column:description;type:text" json:"description" validate:"max=1000"`
	Color        string    `gorm:"column:color;type:varchar(7);not null" json:"color" validate:"required,len=7"`
	DisplayOrder int       `gorm:"column:display_order;not null;default:0;index" json:"display_order"`
	IsActive     bool      `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`

	// Relations
	Issues             []Issue              `gorm:"foreignKey:StatusID;references:StatusID" json:"issues,omitempty"`
	StatusHistoriesNew []IssueStatusHistory `gorm:"foreignKey:NewStatusID;references:StatusID" json:"status_histories_new,omitempty"`
	StatusHistoriesOld []IssueStatusHistory `gorm:"foreignKey:OldStatusID;references:StatusID" json:"status_histories_old,omitempty"`
}

func (IssueStatus) TableName() string {
	return "issue_statuses"
}

// Issue represents a support ticket or issue
type Issue struct {
	IssueID     uint      `gorm:"primaryKey;column:issue_id;autoIncrement" json:"issue_id"`
	ReporterID  uint      `gorm:"column:reporter_id;not null;index" json:"reporter_id" validate:"required"`
	AssigneeID  *uint     `gorm:"column:assignee_id;index" json:"assignee_id,omitempty"`
	StatusID    uint      `gorm:"column:status_id;not null;index" json:"status_id" validate:"required"`
	Title       string    `gorm:"column:title;type:varchar(255);not null" json:"title" validate:"required,min=3,max=255"`
	Description string    `gorm:"column:description;type:text" json:"description" validate:"max=5000"`
	Priority    string    `gorm:"column:priority;type:varchar(20);not null;default:'medium';index" json:"priority" validate:"required,oneof=low medium high critical"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime;index" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relations
	Reporter      User                 `gorm:"foreignKey:ReporterID;references:UserID" json:"reporter,omitempty"`
	Assignee      *Officer             `gorm:"foreignKey:AssigneeID;references:OfficerID" json:"assignee,omitempty"`
	Status        IssueStatus          `gorm:"foreignKey:StatusID;references:StatusID" json:"status,omitempty"`
	StatusHistory []IssueStatusHistory `gorm:"foreignKey:IssueID;references:IssueID" json:"status_history,omitempty"`
	Comments      []Comment            `gorm:"foreignKey:IssueID;references:IssueID" json:"comments,omitempty"`
}

func (Issue) TableName() string {
	return "issues"
}

// IssueStatusHistory tracks status changes for issues
type IssueStatusHistory struct {
	HistoryID   uint      `gorm:"primaryKey;column:history_id;autoIncrement" json:"history_id"`
	IssueID     uint      `gorm:"column:issue_id;not null;index" json:"issue_id"`
	OldStatusID *uint     `gorm:"column:old_status_id;index" json:"old_status_id,omitempty"`
	NewStatusID uint      `gorm:"column:new_status_id;not null;index" json:"new_status_id"`
	ChangedBy   uint      `gorm:"column:changed_by;not null;index" json:"changed_by"`
	Comment     string    `gorm:"column:comment;type:text" json:"comment"`
	ChangedAt   time.Time `gorm:"column:changed_at;autoCreateTime;index" json:"changed_at"`

	// Relations
	Issue            Issue        `gorm:"foreignKey:IssueID;references:IssueID" json:"issue,omitempty"`
	OldStatus        *IssueStatus `gorm:"foreignKey:OldStatusID;references:StatusID" json:"old_status,omitempty"`
	NewStatus        IssueStatus  `gorm:"foreignKey:NewStatusID;references:StatusID" json:"new_status,omitempty"`
	ChangedByOfficer Officer      `gorm:"foreignKey:ChangedBy;references:OfficerID;foreignKeyConstraint:OnDelete:RESTRICT,OnUpdate:CASCADE" json:"changed_by_officer,omitempty"`
}

func (IssueStatusHistory) TableName() string {
	return "issue_status_history"
}

// Comment represents a comment on an issue
type Comment struct {
	CommentID uint      `gorm:"primaryKey;column:comment_id;autoIncrement" json:"comment_id"`
	IssueID   uint      `gorm:"column:issue_id;not null;index" json:"issue_id" validate:"required"`
	UserID    uint      `gorm:"column:user_id;not null" json:"user_id" validate:"required"`
	Content   string    `gorm:"column:content;type:text;not null" json:"content" validate:"required,min=1,max=2000"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;index" json:"created_at"`

	// Relations
	Issue Issue `gorm:"foreignKey:IssueID;references:IssueID" json:"issue,omitempty"`
	User  User  `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}
