package link

type CreateRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type CreateResponse struct {
	Status string `json:"status"`
}

type UpdateRequest struct{}
type UpdateResponse struct{}
type DeleteRequest struct{}
type DeleteResponse struct{}
type GetRequest struct{}
type GetResponse struct{}
