package types

type battery struct {
	SerialId  string `json:"serialid"`
	ModelName string `json:"modelname"`
	Specs     string `json:"specs"`
	Status    string `json:"status"`
}
