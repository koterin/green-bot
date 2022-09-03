package entity

import "strconv"

type ResponseData struct {
	Status   int     `json:"status"`
	Response string  `json:"Response"`
	Origins  []Origs `json:"origins,omitempty"`
}

type Origs struct {
	Origin string `json:"origin,omitempty"`
}

type Recipient struct {
	ID int
}

func (user Recipient) Recipient() string {
	return strconv.Itoa(user.ID)
}
