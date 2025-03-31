package porb

import (
	"time"
)

type post struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	File      struct {
		Width  int    `json:"width"`
		Height int    `json:"height"`
		Ext    string `json:"ext"`
		Size   int    `json:"size"`
		Md5    string `json:"md5"`
		Url    string `json:"url"`
	} `json:"file"`
	Preview struct {
		Width  int    `json:"width"`
		Height int    `json:"height"`
		Url    string `json:"url"`
	} `json:"preview"`
	Sample struct {
		Has        bool   `json:"has"`
		Height     int    `json:"height"`
		Width      int    `json:"width"`
		Url        string `json:"url"`
		Alternates struct {
		} `json:"alternates"`
	} `json:"sample"`
	Score struct {
		Up    int `json:"up"`
		Down  int `json:"down"`
		Total int `json:"total"`
	} `json:"score"`
	Tags struct {
		General     []string      `json:"general"`
		Artist      []string      `json:"artist"`
		Contributor []interface{} `json:"contributor"`
		Copyright   []string      `json:"copyright"`
		Character   []interface{} `json:"character"`
		Species     []string      `json:"species"`
		Invalid     []interface{} `json:"invalid"`
		Meta        []string      `json:"meta"`
		Lore        []interface{} `json:"lore"`
	} `json:"tags"`
	LockedTags []interface{} `json:"locked_tags"`
	ChangeSeq  int           `json:"change_seq"`
	Flags      struct {
		Pending      bool `json:"pending"`
		Flagged      bool `json:"flagged"`
		NoteLocked   bool `json:"note_locked"`
		StatusLocked bool `json:"status_locked"`
		RatingLocked bool `json:"rating_locked"`
		Deleted      bool `json:"deleted"`
	} `json:"flags"`
	Rating        string   `json:"rating"`
	FavCount      int      `json:"fav_count"`
	Sources       []string `json:"sources"`
	Pools         []int    `json:"pools"`
	Relationships struct {
		ParentId          interface{}   `json:"parent_id"`
		HasChildren       bool          `json:"has_children"`
		HasActiveChildren bool          `json:"has_active_children"`
		Children          []interface{} `json:"children"`
	} `json:"relationships"`
	ApproverId   interface{} `json:"approver_id"`
	UploaderId   int         `json:"uploader_id"`
	Description  string      `json:"description"`
	CommentCount int         `json:"comment_count"`
	IsFavorited  bool        `json:"is_favorited"`
	HasNotes     bool        `json:"has_notes"`
	Duration     interface{} `json:"duration"`
}
