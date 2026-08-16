package model

type Category struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name"`
}
