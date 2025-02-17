package stat

import (
	"gorm.io/datatypes"
	"short-link/pkg/db"
	"time"
)

type StatRepository struct {
	*db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{
		Db: db,
	}
}

func (repo *StatRepository) AddClick(linkID uint) {
	// Если записи на текущий день нет, то ее необходимо создать
	//Если запись есть, то необходимо добавить 1 клик
	currentDate := datatypes.Date(time.Now())
	var stat Stat
	repo.DB.Find(&stat, "link_id = ? and date = ?", linkID, currentDate)
	if stat.ID == 0 {
		repo.DB.Create(&Stat{
			LinkID: linkID,
			Date:   currentDate,
			Clicks: 1,
		})
	} else {
		stat.Clicks++
		repo.DB.Save(&stat)
	}
}

func (repo *StatRepository) GetClicks(from, to time.Time) []Stat {
	var stats []Stat
	repo.DB.Where("date >= ? AND date <= ?", from, to).Group("updated_at").Find(&stats)
	//repo.DB.Where("date >= ? AND date <= ?", from, to).Find(&stats)
	return stats
}
