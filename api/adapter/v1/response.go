package v1

// NewGetPageResponseSuccess returns a GetPageResponse with the given page.
func NewGetPageResponseSuccess(page *Page) *GetPageResponse {
	return &GetPageResponse{
		Response: &GetPageResponse_Success{
			Success: page,
		},
	}
}

// NewGetPageResponseError returns a GetPageResponse with the given error.
func NewGetPageResponseError(err *Error) *GetPageResponse {
	return &GetPageResponse{
		Response: &GetPageResponse_Error{
			Error: err,
		},
	}
}
