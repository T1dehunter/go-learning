package weboscket

type UserAuthMessage struct {
	Name    string `json:"name"`
	Payload struct {
		UserID      int32  `json:"userID"`
		UserName    string `json:"userName"`
		AccessToken string `json:"accessToken"`
	} `json:"payload"`
}
