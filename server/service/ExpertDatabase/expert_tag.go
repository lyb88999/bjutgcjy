package ExpertDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	"gorm.io/gorm"
)

type ExpertTagService struct{}

// CreateExpertTag 创建标签
func (expertTagService *ExpertTagService) CreateExpertTag(tag *ExpertDatabase.ExpertTag) (err error) {
	return global.GVA_DB.Create(tag).Error
}

// DeleteExpertTag 删除标签
func (expertTagService *ExpertTagService) DeleteExpertTag(ID string, userID uint) (err error) {
	if err = global.GVA_DB.Model(&ExpertDatabase.ExpertTag{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
		return err
	}
	return global.GVA_DB.Delete(&ExpertDatabase.ExpertTag{}, "id = ?", ID).Error
}

// UpdateExpertTag 更新标签
func (expertTagService *ExpertTagService) UpdateExpertTag(tag ExpertDatabase.ExpertTag) (err error) {
	return global.GVA_DB.Save(&tag).Error
}

// GetExpertTag 根据ID获取标签
func (expertTagService *ExpertTagService) GetExpertTag(ID string) (tag ExpertDatabase.ExpertTag, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&tag).Error
	return
}

// GetExpertTagInfoList 分页获取标签列表
func (expertTagService *ExpertTagService) GetExpertTagInfoList(info ExpertDatabaseReq.ExpertTagSearch) (list []ExpertDatabase.ExpertTag, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&ExpertDatabase.ExpertTag{})
	var tags []ExpertDatabase.ExpertTag

	if info.TagType != "" {
		db = db.Where("tag_type = ?", info.TagType)
	}
	if info.TagValue != "" {
		db = db.Where("tag_value LIKE ?", "%"+info.TagValue+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	db = db.Order("id desc")
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Find(&tags).Error
	return tags, total, err
}

// SetExpertTags 覆盖式设置某专家的标签关联（先清空再写入）
func (expertTagService *ExpertTagService) SetExpertTags(expertId uint, tagIds []uint) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("expert_id = ?", expertId).Delete(&ExpertDatabase.ExpertTagRelation{}).Error; err != nil {
			return err
		}
		if len(tagIds) == 0 {
			return nil
		}
		relations := make([]ExpertDatabase.ExpertTagRelation, 0, len(tagIds))
		for _, tagId := range tagIds {
			relations = append(relations, ExpertDatabase.ExpertTagRelation{ExpertId: expertId, TagId: tagId})
		}
		return tx.Create(&relations).Error
	})
}

// GetExpertTagsByExpertId 获取某专家关联的全部标签
func (expertTagService *ExpertTagService) GetExpertTagsByExpertId(expertId uint) (tags []ExpertDatabase.ExpertTag, err error) {
	err = global.GVA_DB.
		Joins("JOIN expert_tag_relation ON expert_tag_relation.tag_id = expert_tag.id").
		Where("expert_tag_relation.expert_id = ?", expertId).
		Find(&tags).Error
	return
}
