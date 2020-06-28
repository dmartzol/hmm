package mockdb

import (
	"time"

	"github.com/dmartzol/hmmm/internal/models"
)

// MockDB represents a database
type MockDB struct{}

func NewMockDB() (*MockDB, error) {
	return &MockDB{}, nil
}

func (db *MockDB) AccountExists(email string) (bool, error) {
	if email == "registered@email.com" {
		return true, nil
	}
	return false, nil
}

func (db *MockDB) Account(id int64) (*models.Account, error) {
	var a models.Account
	return &a, nil
}

func (db *MockDB) Accounts() (models.Accounts, error) {
	var accs []*models.Account
	return accs, nil
}

func (db *MockDB) AccountWithCredentials(email, allegedPassword string) (*models.Account, error) {
	var a models.Account
	return &a, nil
}

func (db *MockDB) CreateAccount(first, last, email, password, confirmationCode string, dob time.Time, gender, phone *string) (*models.Account, *models.Confirmation, error) {
	a := models.Account{
		Row: models.Row{
			ID:         1,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		},
		FirstName: first,
		LastName:  last,
		Email:     email,
		DOB:       dob,
		Gender:    gender,
	}
	var cc models.Confirmation
	return &a, &cc, nil
}

func (db *MockDB) SessionFromIdentifier(identifier string) (*models.Session, error) {
	var s models.Session
	return &s, nil
}

func (db *MockDB) CreateSession(accountID int64) (*models.Session, error) {
	var s models.Session
	return &s, nil
}

func (db *MockDB) DeleteSession(identifier string) error {
	return nil
}

func (db *MockDB) CleanSessionsOlderThan(age time.Duration) (int64, error) {
	return 2, nil
}

func (db *MockDB) UpdateSession(sessionToken string) (*models.Session, error) {
	var s models.Session
	return &s, nil
}

func (db *MockDB) CreateConfirmation(accountID int64, t models.ConfirmationType) (*models.Confirmation, error) {
	var cc models.Confirmation
	return &cc, nil
}

func (db *MockDB) PendingConfirmationByKey(key string) (*models.Confirmation, error) {
	var c models.Confirmation
	return &c, nil
}

func (db *MockDB) FailedConfirmationIncrease(id int64) (*models.Confirmation, error) {
	var c models.Confirmation
	return &c, nil
}

func (db *MockDB) Confirm(id int64) (*models.Confirmation, error) {
	var c models.Confirmation
	return &c, nil
}

func (db *MockDB) Role(roleID int64) (*models.Role, error) {
	return nil, nil
}

func (db *MockDB) RoleExists(name string) (bool, error) {
	return false, nil
}

func (db *MockDB) CreateRole(name string) (*models.Role, error) {
	var r models.Role
	return &r, nil
}

func (db *MockDB) RolesForAccount(accountID int64) (models.Roles, error) {
	return nil, nil
}

func (db *MockDB) AddAccountRole(roleID, accountID int64) (*models.AccountRole, error) {
	return nil, nil
}
