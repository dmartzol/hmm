package controllers

const (
	apiVersionNumber = "0.0.1"
	hmmmCookieName   = "Hmmm-Cookie"
)

type storage interface {
	sessionStorage
	accountStorage
}

// API represents something
type API struct {
	storage
}

func NewAPI(db storage) (*API, error) {
	return &API{db}, nil
}
