package model

type GetUserProfile struct {
	AvatarPath *string `json:"avatar_path,omitempty"`
	FirstName  *string `json:"first_name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
	Username   string  `json:"username"`
	Phone      string  `json:"phone"`
	Email      *string `json:"email,omitempty"`
}

func (g *GetUserProfile) Sanitize() {
	g.Username = utils.SanitizeString(g.Username)
	g.Phone = utils.SanitizeString(g.Phone)

	if g.FirstName != nil {
		s := utils.SanitizeString(*g.FirstName)
		g.FirstName = &s
	}
	if g.LastName != nil {
		s := utils.SanitizeString(*g.LastName)
		g.LastName = &s
	}
	if g.Email != nil {
		s := utils.SanitizeString(*g.Email)
		g.Email = &s
	}
}