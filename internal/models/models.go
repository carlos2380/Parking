package models

type Car struct {
	PlateNumber string `json:"plate_number"`
	EntryDate   string `json:"entry_date"`
}

type Ticket struct {
	PlateNumber string `json:"plate_number"`
	EntryDate   string `json:"entry_date"`
	ExitDate    string `json:"exit_date"`
	Price       string `json:"price"`
}
