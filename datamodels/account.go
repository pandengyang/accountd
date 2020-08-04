package datamodels

type Account struct {
	Id               int64  `json:"id,omitemtpy"`
	Nickname         string `json:"nickname,omitemtpy"`
	Phone            string `json:"phone,omitemtpy"`
	VerificationCode string `json:verification_code",omitemtpy"`
	Password         string `json:"password,omitemtpy"`
	Salt             string `json:"salt,omitemtpy"`
	State            string `json:"state,omitemtpy"`
	CreatedAt        int64  `json:"created_at,omitemtpy"`
}
