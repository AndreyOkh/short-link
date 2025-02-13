package link

import (
	"gorm.io/gorm/clause"
	"short-link/pkg/db"
)

type LinkRepository struct {
	Database *db.Db
}

func NewLinkRepository(database *db.Db) *LinkRepository {
	return &LinkRepository{
		Database: database,
	}
}

func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	if result := repo.Database.DB.Create(&link); result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repo *LinkRepository) FindByHash(hash string) (*Link, error) {
	var link Link
	result := repo.Database.DB.Where("hash = ?", hash).First(&link)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (repo *LinkRepository) DeleteByID(linkID uint64) error {
	err := repo.CheckExistByID(linkID)
	if err != nil {
		return err
	} else {
		result := repo.Database.DB.Delete(&Link{}, linkID)
		if result.Error != nil {
			return result.Error
		}
		return nil
	}
}

func (repo *LinkRepository) Update(link *Link) (*Link, error) {
	if result := repo.Database.DB.Clauses(clause.Returning{}).Updates(&link); result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repo *LinkRepository) CheckExistByID(linkID uint64) error {
	result := repo.Database.DB.Model(Link{}).Where("id = ?", linkID).First(&Link{})
	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}
