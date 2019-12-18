package organization

import (
	"github.com/ATechnoHazard/hades-2/pkg"
	"github.com/ATechnoHazard/hades-2/pkg/user"
	"github.com/jinzhu/gorm"
)

type repo struct {
	DB *gorm.DB
}

func (r *repo) Save(organization *Organization) error {
	err := r.DB.Save(organization).Error
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return pkg.ErrNotFound
	default:
		return pkg.ErrDatabase
	}
}

func (r *repo) Find(orgID uint) (*Organization, error) {
	org := &Organization{ID: orgID}
	err := r.DB.Find(org).Association("Events").Find(&org.Events).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	case nil:
		return org, err
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) FindAll() ([]Organization, error) {
	var organizations []Organization
	err := r.DB.Model(Organization{}).Find(&organizations).Error
	return organizations, err
}

func (r *repo) Delete(orgID uint) error {
	err := r.DB.Delete(&Organization{ID: orgID}).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return pkg.ErrNotFound
	case nil:
		return nil
	default:
		return pkg.ErrDatabase
	}
}

func (r *repo) SaveJoinReq(request *JoinRequest) error {
	tx := r.DB.Begin()
	u := &user.User{Email: request.Email}
	if tx.Find(u).Error == gorm.ErrRecordNotFound {
		tx.Rollback()
		return pkg.ErrNotFound
	}

	for _, org := range u.Organizations {
		if org.ID == request.OrganizationID {
			tx.Rollback()
			return pkg.ErrAlreadyExists
		}
	}

	err := tx.Find(request).Error
	switch err {
	case gorm.ErrRecordNotFound:
		tx.Save(request)
		tx.Commit()
		return nil
	case nil:
		tx.Rollback()
		return pkg.ErrAlreadyExists
	default:
		tx.Rollback()
		return pkg.ErrDatabase
	}
}

func (r *repo) FindAllJoinReq(orgID uint) ([]JoinRequest, error) {
	org := &Organization{ID: orgID}
	err := r.DB.Find(org).Association("JoinRequests").Find(&org.JoinRequests).Error
	switch err {
	case nil:
		return org.JoinRequests, nil
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) AcceptJoinReq(request *JoinRequest) error {
	tx := r.DB.Begin()
	u := &user.User{Email: request.Email}
	org := &Organization{ID: request.OrganizationID}
	err := tx.Find(request).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	err = tx.Find(org).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	err = tx.Find(u).Association("Organizations").Append(org).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrAlreadyExists
		default:
			return pkg.ErrDatabase
		}
	}

	err = tx.Delete(request).Error
	switch err {
	case nil:
		tx.Commit()
		return nil
	case gorm.ErrRecordNotFound:
		tx.Rollback()
		return pkg.ErrAlreadyExists
	default:
		tx.Rollback()
		return pkg.ErrDatabase
	}
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{DB: db}
}
