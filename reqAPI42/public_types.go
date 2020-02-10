package reqAPI42

import (
	"time"
)

// API42 type which use to interact with 42's API
type API42 struct {
	keys		apiKeys
	rlLastReqSec	time.Time
	rlNbReqSec	uint
	campus		*API42Campus
	cursus		*API42Cursus
	// locations	*[]API42Location
	// projects	*[]API42Project
}

type API42User struct {
	ID	uint	`json:"id"`
	Login	string	`json:"login"`
}

type API42Cursus struct {
	ID	uint	`json:"id"`
	Name	string	`json:"name"`
}

type API42Campus struct {
	ID	uint	`json:"id"`
	Name	string	`json:"name"`
}

type API42Location struct {
	// ID	uint		`json:"id"`
	// EndAt	JSONTime	`json:"end_at"`
	Host	string		`json:"host"`
	User	API42User 	`json:"user"`
}

type API42Project struct {
	ID	uint			`json:"id"`
	Name	string			`json:"name"`
	Parent	*API42ProjectParent	`json:"parent"`
	Campus	[]API42Campus		`json:"campus"`
}

type API42ProjectUser struct {
	ID	uint		`json:"id"`
	User	API42User	`json:"user"`
}

type API42ProjectParent struct {
	ID	uint	`json:"id"`
	Name	string	`json:"name"`
	// Slug	string	`json:"slug"`
}
