package entities

type People struct {
	ID             int    `json:"id"`
	PassportSeries int    `json:"passport_series"`
	PassportNumber int    `json:"passport_number"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}
