package domain

type Tour struct {
	ID                 int
	AgencyID           int
	Name               string
	StartDate          string
	EndDate            string
	Description        string
	LastEnrollmentDate string
	Price              int
	Discount           int
	Status             string
}
