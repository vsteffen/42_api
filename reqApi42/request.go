package reqApi42

// API42 type which use to interact with 42's API
type API42 struct {
	keys apiKeys
}

// New create new API42 obj
func New() *API42 {
	tmp := API42{}
	initKeys(&tmp)
	return &tmp
}
