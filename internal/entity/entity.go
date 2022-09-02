package entity

import "strconv"

type ResponseData struct {
	Status   int    `json:"status"`
	Origins  string `json:"origins"`
	Response string `json:"Response"`
}

type Recipient struct {
	ID int
}

func (user Recipient) Recipient() string {
	return strconv.Itoa(user.ID)
}
