package user

type CreateRequest struct {
	Name string

	GEWISID uint

	OrganIDs []uint
} // @name UserCreateRequest

type FilterParams struct {
	ID      *uint `form:"id"`
	GEWISID *uint `form:"gewisId"`
	OrganID *uint `form:"organId"`
} // @name UserFilterParams
