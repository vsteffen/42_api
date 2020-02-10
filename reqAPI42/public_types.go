package reqAPI42

import (
	"time"
)

// API42 is the object used to communicate with the 42's API
type API42 struct {
	keys         apiKeys
	rlLastReqSec time.Time
	rlNbReqSec   uint
	campus       *API42Campus
	cursus       *API42Cursus
	// locations	*[]API42Location
	// projects	*[]API42Project
}

// API42User object from API 42 (https://api.intra.42.fr/apidoc/2.0/users.html)
type API42User struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
}

// API42Cursus object from API 42 (https://api.intra.42.fr/apidoc/2.0/cursus.html)
type API42Cursus struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// API42Campus object from API 42 (https://api.intra.42.fr/apidoc/2.0/campus.html)
type API42Campus struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// API42Location object from API 42 (https://api.intra.42.fr/apidoc/2.0/locations.html)
type API42Location struct {
	// ID	uint		`json:"id"`
	// EndAt	JSONTime	`json:"end_at"`
	Host string    `json:"host"`
	User API42User `json:"user"`
}

// API42Project object from API 42 (https://api.intra.42.fr/apidoc/2.0/projects.html)
type API42Project struct {
	ID     uint                `json:"id"`
	Name   string              `json:"name"`
	Parent *API42ProjectParent `json:"parent"`
	Campus []API42Campus       `json:"campus"`
}

// API42ProjectUser object from API 42 (https://api.intra.42.fr/apidoc/2.0/projects_users.html)
type API42ProjectUser struct {
	ID   uint      `json:"id"`
	User API42User `json:"user"`
}

// API42ProjectParent object from API 42 (https://api.intra.42.fr/apidoc/2.0/projects.html)
type API42ProjectParent struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	// Slug	string	`json:"slug"`
}
