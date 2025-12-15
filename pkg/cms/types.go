package cms

import (
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/settings"
	"github.com/dtekltd/common/pkg/site"
	"github.com/dtekltd/common/pkg/users"
	"github.com/dtekltd/common/types"
	"gorm.io/datatypes"
)

const MODULE = "cms"

type CarouselItem struct {
	Title  string `json:"title,omitempty"`
	Source string `json:"source,omitempty"`
}

type Carousel struct {
	ID    string          `json:"id,omitempty"`
	Items *[]CarouselItem `json:"items,omitempty"`
}

type SwiperItem struct {
	Title  string `json:"title,omitempty"`
	Source string `json:"source,omitempty"`
}

type Swiper struct {
	ID    string        `json:"id,omitempty"`
	Items *[]SwiperItem `json:"items,omitempty"`
}

type Category struct {
	ID          uint64 `json:"id,omitempty" gorm:"primaryKey"`
	Type        string `json:"type,omitempty" gorm:"type:varchar(16);index"`
	Name        string `json:"name,omitempty" gorm:"type:varchar(32);"`
	Alias       string `json:"alias,omitempty" gorm:"type:varchar(32);uniqueIndex"`
	Description string `json:"description,omitempty" gorm:"type:varchar(128)"`
	Frequency   int    `json:"frequency"`
	Ordering    int    `json:"ordering"`
	CreatedAt   uint64 `json:"created_at,omitempty"`
	UpdatedAt   uint64 `json:"updated_at,omitempty"`
	PublishedAt uint64 `json:"published_at,omitempty"`
	DeletedAt   uint64 `json:"-"`
}

type Meta struct {
	View     string     `json:"view,omitempty"`
	Layout   string     `json:"layout,omitempty"`
	PageMeta *site.Meta `json:"pageMeta,omitempty"`
}

type Author struct {
	ID        uint64 `json:"id,omitempty" gorm:"primaryKey"`
	Name      string `json:"name,omitempty" gorm:"type:varchar(64);"`
	Email     string `json:"email,omitempty" gorm:"type:varchar(64);uniqueIndex"`
	Phone     string `json:"phone,omitempty" gorm:"type:varchar(16)"`
	Username  string `json:"username,omitempty" gorm:"type:varchar(64);uniqueIndex"`
	AvatarUrl string `json:"avatarUrl,omitempty" gorm:"type:varchar(256)"`
	Posts     []Post `json:"posts,omitempty" gorm:"foreignKey:CreatedBy"`
}

type Post struct {
	ID          uint64                    `json:"id,omitempty" gorm:"primaryKey"`
	ParentID    *database.NullableUint64  `json:"parentId,omitempty"`
	Type        string                    `json:"type,omitempty" gorm:"type:varchar(16);index"`
	Name        string                    `json:"name,omitempty" gorm:"type:varchar(128);"`
	Path        string                    `json:"path,omitempty" gorm:"type:varchar(128);index"`
	Content     string                    `json:"content,omitempty"`
	ContentType string                    `json:"contentType,omitempty" gorm:"type:varchar(8)"`
	Ordering    int                       `json:"ordering,omitempty"`
	Meta        *datatypes.JSONType[Meta] `json:"meta,omitempty" gorm:"type:text"`
	CreatedBy   uint64                    `json:"createdBy,omitempty"`
	CreatedAt   uint64                    `json:"created_at,omitempty"`
	UpdatedAt   uint64                    `json:"updated_at,omitempty"`
	PublishedAt uint64                    `json:"published_at,omitempty"`
	DeletedAt   uint64                    `json:"-"`
	Author      *users.Account            `json:"author,omitempty" gorm:"foreignKey:CreatedBy"`
	Categories  []Category                `json:"categories,omitempty" gorm:"many2many:cms_post_categories"`
}

type PostCategory struct {
	PostID     uint64 `gorm:"primaryKey"`
	CategoryID uint64 `gorm:"primaryKey"`
	Post       Post
	Category   Category
}

type CmsParams struct {
	HomePageID   int           `json:"homePageId"`
	BlogPageSize int           `json:"blogPageSize"`
	CustomParams *types.Params `json:"customParams,omitempty"`
}

type CmsSetting struct {
	*settings.Setting
	Params *datatypes.JSONType[CmsParams] `json:"params" gorm:"type:text"`
}

// TableName overrides the table name
func (CmsSetting) TableName() string {
	return "com_settings"
}

func (Category) TableName() string {
	return "cms_categories"
}

func (Post) TableName() string {
	return "cms_posts"
}

func (Author) TableName() string {
	return "user_accounts"
}

func (PostCategory) TableName() string {
	return "cms_post_categories"
}

func Migrate() {
	database.DB.AutoMigrate(&Category{}, &Post{})
}
