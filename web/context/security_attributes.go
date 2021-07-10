package context

type SecurityAttributes struct {
	UserId            string `json:"user_id"`
	TechnicalUsername string `json:"technical_username"`
}
