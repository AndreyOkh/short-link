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

func (repo *StatRepository) GetStats(by string, from, to time.Time) []GetStatResponse {
	var stats []GetStatResponse
	var selectQuery string
	switch by {
	case GroupByMonth:
		selectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	case GroupByDay:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	}
	repo.DB.Table("stats").
		Select(selectQuery).
		Where("date BETWEEN ? AND ?", from, to).
		Group("period").
		Order("period").
		Scan(&stats)
	return stats
}
