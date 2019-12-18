package reqApi42

import (

)

type Api42 struct {
	keys	apiKeys
}

func New() (*Api42) {
	tmp := Api42{}
	initKeys(&tmp)
	return &tmp
}
