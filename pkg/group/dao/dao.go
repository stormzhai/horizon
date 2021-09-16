package dao

import (
	"context"
	"fmt"
	"strings"

	"g.hz.netease.com/horizon/common"
	"g.hz.netease.com/horizon/lib/orm"
	"g.hz.netease.com/horizon/lib/q"
	"g.hz.netease.com/horizon/pkg/group/models"
	"gorm.io/gorm"
)

type DAO interface {
	CheckUnique(ctx context.Context, group *models.Group) error
	Create(ctx context.Context, group *models.Group) (uint, error)
	Delete(ctx context.Context, id uint) error
	Get(ctx context.Context, id uint) (*models.Group, error)
	GetByPath(ctx context.Context, path string) (*models.Group, error)
	GetByNameFuzzily(ctx context.Context, name string) ([]*models.Group, error)
	GetByFullNamesRegexpFuzzily(ctx context.Context, names *[]string) ([]*models.Group, error)
	Update(ctx context.Context, group *models.Group) error
	ListWithoutPage(ctx context.Context, query *q.Query) ([]*models.Group, error)
	List(ctx context.Context, query *q.Query) ([]*models.Group, int64, error)
}

// New returns an instance of the default DAO
func New() DAO {
	return &dao{}
}

type dao struct{}

func (d *dao) GetByFullNamesRegexpFuzzily(ctx context.Context, names *[]string) ([]*models.Group, error) {
	if names == nil || (len(*names)) == 0 {
		return []*models.Group{}, nil
	}

	db, err := orm.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	var groups []*models.Group
	regNames := make([]string, len(*names))
	for i, n := range *names {
		regNames[i] = "^" + n
	}

	namesRegexp := strings.Join(regNames, "|")

	result := db.Where("full_name regexp ?", "^"+namesRegexp).Order("id desc").Find(&groups)

	return groups, result.Error
}

func (d *dao) GetByNameFuzzily(ctx context.Context, name string) ([]*models.Group, error) {
	db, err := orm.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	var groups []*models.Group
	result := db.Where("name LIKE ?", "%"+name+"%").Find(&groups)

	return groups, result.Error
}

func (d *dao) CheckUnique(ctx context.Context, group *models.Group) error {
	db, err := orm.FromContext(ctx)
	if err != nil {
		return err
	}

	query := map[string]interface{}{
		"parent_id": group.ParentID,
		"name":      group.Name,
	}

	queryResult := &models.Group{}
	result := db.Model(&models.Group{}).Where(query).Find(queryResult)

	// update group conflict, has another record with the same parentId & name
	if group.ID > 0 && queryResult.ID > 0 && queryResult.ID != group.ID {
		return common.ErrNameConflict
	}

	// create group conflict
	if group.ID == 0 && result.RowsAffected > 0 {
		return common.ErrNameConflict
	}

	return nil
}

func (d *dao) Create(ctx context.Context, group *models.Group) (uint, error) {
	db, err := orm.FromContext(ctx)
	if err != nil {
		return 0, err
	}
	// conflict if there's record with the same parentId and name
	err = d.CheckUnique(ctx, group)
	if err != nil {
		return 0, err
	}

	result := db.Create(group)

	return group.ID, result.Error
}

func (d *dao) Delete(ctx context.Context, id uint) error {
	db, err := orm.FromContext(ctx)
	if err != nil {
		return err
	}

	result := db.Where("id = ?", id).Delete(&models.Group{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	// todo delete children

	return result.Error
}

func (d *dao) Get(ctx context.Context, id uint) (*models.Group, error) {
	db, err := orm.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	var group *models.Group
	result := db.First(&group, id)

	return group, result.Error
}

func (d *dao) GetByPath(ctx context.Context, path string) (*models.Group, error) {
	db, err := orm.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	var group *models.Group
	result := db.First(&group, "path = ?", path)

	return group, result.Error
}

func (d *dao) ListWithoutPage(ctx context.Context, query *q.Query) ([]*models.Group, error) {
	db, err := orm.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	var groups []*models.Group

	sort := orm.FormatSortExp(query)
	result := db.Order(sort).Where(query.Keywords).Find(&groups)

	return groups, result.Error
}

func (d *dao) List(ctx context.Context, query *q.Query) ([]*models.Group, int64, error) {
	db, err := orm.FromContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	var groups []*models.Group

	sort := orm.FormatSortExp(query)
	offset := (query.PageNumber - 1) * query.PageSize
	var count int64
	result := db.Order(sort).Where(query.Keywords).Offset(offset).Limit(query.PageSize).Find(&groups).
		Offset(-1).Count(&count)
	return groups, count, result.Error
}

func (d *dao) Update(ctx context.Context, group *models.Group) error {
	db, err := orm.FromContext(ctx)
	if err != nil {
		return err
	}

	// conflict if there's record with the same parentId and name
	err = d.CheckUnique(ctx, group)
	if err != nil {
		return err
	}

	result := db.Model(group).Select("Name", "Description", "VisibilityLevel").Where("deleted_at is null").Updates(group)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	// todo modify children's path

	// update self & children's fullNames
	var oldGroup *models.Group
	result = db.First(&oldGroup, group.ID)
	oldFullName := oldGroup.FullName
	split := strings.Split(strings.ReplaceAll(oldFullName, " ", ""), "/")
	split = append(split[:len(split)-1], group.Name)
	newFullName := strings.Join(split, " / ")
	// update `group` set full_name=replace(full_name, full_name, concat('a / e', substring(full_name, 6))) where full_name regexp '^(a / d)'
	updateSQL := fmt.Sprintf("update `group` set full_name=replace(full_name, full_name, concat('%s', substring(full_name, %d))) where full_name regexp '^(%s)'", newFullName, len(oldFullName)+1, oldFullName)
	result = db.Exec(updateSQL)

	return result.Error
}
