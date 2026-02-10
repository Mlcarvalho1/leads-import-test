package models

type ImportRequest struct {
	Name      string `json:"name"`
	AccountID int    `json:"account_id"`
	SourceID  int    `json:"source_id"`
	TagIDs    []int  `json:"tag_ids"`
}
