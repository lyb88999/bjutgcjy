package ExpertDatabase

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ExpertProfile 专家主档：背景信息画像、学科领域画像、研究主题画像打平字段 + 审核流转/推荐排序所需的状态与缓存字段
type ExpertProfile struct {
	global.GVA_MODEL
	// 一、背景信息画像
	Name             string     `json:"name" form:"name" gorm:"column:name;comment:姓名;" binding:"required"`                     //姓名
	Gender           string     `json:"gender" form:"gender" gorm:"column:gender;comment:性别;"`                                  //性别
	BirthDate        *time.Time `json:"birthDate" form:"birthDate" gorm:"column:birth_date;comment:出生年月;"`                      //出生年月
	Ethnicity        string     `json:"ethnicity" form:"ethnicity" gorm:"column:ethnicity;comment:民族;"`                         //民族
	PoliticalStatus  string     `json:"politicalStatus" form:"politicalStatus" gorm:"column:political_status;comment:政治面貌;"`    //政治面貌
	UnitName         string     `json:"unitName" form:"unitName" gorm:"column:unit_name;comment:所在单位;"`                         //所在单位
	Department       string     `json:"department" form:"department" gorm:"column:department;comment:所属院系或部门;"`                 //所属院系或部门
	AdminTitle       string     `json:"adminTitle" form:"adminTitle" gorm:"column:admin_title;comment:行政职务;"`                   //行政职务
	TechTitle        string     `json:"techTitle" form:"techTitle" gorm:"column:tech_title;comment:专业技术职称;"`                    //专业技术职称
	Phone            string     `json:"phone" form:"phone" gorm:"column:phone;comment:办公电话;"`                                   //办公电话
	Mobile           string     `json:"mobile" form:"mobile" gorm:"column:mobile;comment:手机号码;"`                                //手机号码
	Email            string     `json:"email" form:"email" gorm:"column:email;comment:电子邮箱;"`                                   //电子邮箱
	Address          string     `json:"address" form:"address" gorm:"column:address;comment:通讯地址;"`                             //通讯地址
	HighestDegree    string     `json:"highestDegree" form:"highestDegree" gorm:"column:highest_degree;comment:最高学位;"`          //最高学位
	HighestEducation string     `json:"highestEducation" form:"highestEducation" gorm:"column:highest_education;comment:最高学历;"` //最高学历
	GraduateSchool   string     `json:"graduateSchool" form:"graduateSchool" gorm:"column:graduate_school;comment:毕业院校;"`       //毕业院校
	Major            string     `json:"major" form:"major" gorm:"column:major;comment:所学专业;"`                                   //所学专业
	TalentProgram    string     `json:"talentProgram" form:"talentProgram" gorm:"column:talent_program;comment:人才计划;"`          //人才计划
	HonorTitle       string     `json:"honorTitle" form:"honorTitle" gorm:"column:honor_title;comment:荣誉称号;"`                   //荣誉称号

	// 二、学科领域画像
	DisciplineL1       string `json:"disciplineL1" form:"disciplineL1" gorm:"column:discipline_l1;comment:一级学科;"`                          //一级学科
	DisciplineL2       string `json:"disciplineL2" form:"disciplineL2" gorm:"column:discipline_l2;comment:二级学科;"`                          //二级学科
	CrossDiscipline    string `json:"crossDiscipline" form:"crossDiscipline" gorm:"column:cross_discipline;comment:交叉学科领域;"`               //交叉学科领域
	DisciplinePlatform string `json:"disciplinePlatform" form:"disciplinePlatform" gorm:"column:discipline_platform;comment:所属学科平台/研究基地;"` //所属学科平台/研究基地

	// 三、研究主题画像
	ResearchDirections string `json:"researchDirections" form:"researchDirections" gorm:"column:research_directions;comment:近五年重点研究方向;"` //近五年重点研究方向
	ResearchKeywords   string `json:"researchKeywords" form:"researchKeywords" gorm:"column:research_keywords;comment:研究关键词;"`           //研究关键词，逗号分隔
	PolicyFields       string `json:"policyFields" form:"policyFields" gorm:"column:policy_fields;comment:所属政策领域;"`                      //所属政策领域
	ResearchObjects    string `json:"researchObjects" form:"researchObjects" gorm:"column:research_objects;comment:主要研究对象;"`             //主要研究对象
	FocusTopics        string `json:"focusTopics" form:"focusTopics" gorm:"column:focus_topics;comment:当前关注议题;"`                         //当前关注议题
	RegionExpertise    string `json:"regionExpertise" form:"regionExpertise" gorm:"column:region_expertise;comment:区域或国别研究专长;"`          //区域或国别研究专长
	MethodExpertise    string `json:"methodExpertise" form:"methodExpertise" gorm:"column:method_expertise;comment:研究方法专长;"`             //研究方法专长

	// 审核工作流字段（状态机逻辑在里程碑 B 补充，这里先建列）
	Status      string `json:"status" form:"status" gorm:"column:status;default:draft;comment:审核状态;"`      //draft/pending_org_review/org_rejected/pending_city_review/city_rejected/published
	OrgId       *uint  `json:"orgId" form:"orgId" gorm:"column:org_id;index;comment:申报单位ID;"`              //申报单位ID
	SubmittedBy *uint  `json:"submittedBy" form:"submittedBy" gorm:"column:submitted_by;comment:申报人用户ID;"` //申报人用户ID
	CurrentStep string `json:"currentStep" form:"currentStep" gorm:"column:current_step;comment:当前审核环节;"`  //当前审核环节

	// 推荐排序缓存字段（算分逻辑在里程碑 C 补充，这里先建列）
	AchievementScore float64    `json:"achievementScore" form:"achievementScore" gorm:"column:achievement_score;comment:成果得分缓存;"` //成果得分缓存
	InfluenceScore   float64    `json:"influenceScore" form:"influenceScore" gorm:"column:influence_score;comment:决策影响得分缓存;"`     //决策影响得分缓存
	CompositeScore   float64    `json:"compositeScore" form:"compositeScore" gorm:"column:composite_score;comment:综合排序得分缓存;"`     //综合排序得分缓存
	ScoreUpdatedAt   *time.Time `json:"scoreUpdatedAt" form:"scoreUpdatedAt" gorm:"column:score_updated_at;comment:得分缓存更新时间;"`    //得分缓存更新时间

	CreatedBy uint `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 专家主档 ExpertProfile自定义表名 expert_profile
func (ExpertProfile) TableName() string {
	return "expert_profile"
}
