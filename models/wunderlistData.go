package models

type WunderlistListData struct {
	Name string `json:"title"`
	Id   int64  `json:"id"`
}

type WunderlistTaskData struct {
	Id     int64  `json:"id"`
	ListId int64  `json:"list_id"`
	Name   string `json:"title"`
	Type   string `json:"type"`
}

type WunderlistDto struct {
	Name  string
	Tasks []WunderlistTaskDto
}

type WunderlistTaskDto struct {
	Name string
	Type string
}

func (list WunderlistListData) MapToDto() WunderlistDto {

}
