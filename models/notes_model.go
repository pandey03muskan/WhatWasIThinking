package models

type GetNotes struct {
	ID      string `json:"id" bson:"_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CreateNoteRequestBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
