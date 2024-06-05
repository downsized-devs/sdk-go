package audit

type auditKey string

const (
	eventTypeHttp   auditKey = "http"
	eventTypeDomain auditKey = "domain"
	logType         auditKey = "application/audit-trail"
)

type Collection struct {
	EventName        string
	EventDescription string
	RequestBody      any
	InsertParam      any
	SelectParam      any
	UpdateParam      any
	Error            error
}

type inputs struct {
	Insert any `json:"insertParam,omitempty"`
	Select any `json:"selectParam,omitempty"`
	Update any `json:"updateParam,omitempty"`
}
