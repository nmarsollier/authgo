package routes

type permission struct {
	Permissions []string `json:"permissions" binding:"required"`
}
