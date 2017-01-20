package structs

type Country struct {
	Id        string `json:"-"`
	CodeA2    string `json:"codeA2"`
	CodeA3    string `json:"codeA3"`
	CodePhone string `json:"phoneCode"`
	En        string `json:"en"`
	Ru        string `json:"ru"`
	Enabled   bool   `json:"-"`
}
