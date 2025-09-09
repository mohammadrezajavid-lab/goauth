package goauth

type GenerateOTPRequest struct {
	PhoneNumber string `json:"phone_number"`
}

type GenerateOTPResponse struct {
	OTP string `json:"otp"`
}

type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number"`
	OTP         string `json:"otp"`
}

type VerifyOTPResponse struct {
	Token string `json:"token"`
	IsNew bool   `json:"is_new"`
}

type GetUserRequest struct {
	ID int64
}

type ListUsersRequest struct {
	Page     int
	PageSize int
	Search   string
}

type PaginationMetadata struct {
	CurrentPage  int `json:"currentPage"`
	PageSize     int `json:"pageSize"`
	TotalRecords int `json:"totalRecords"`
	TotalPages   int `json:"totalPages"`
}

type ListUsersResponse struct {
	Users    []User             `json:"users"`
	Metadata PaginationMetadata `json:"metadata"`
}
