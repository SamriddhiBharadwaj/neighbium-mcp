package main

// These structs match the JSON that Neighbium returns.
type DirectoryEntry struct {
	Member     Member          `json:"member"`
	Apartments []ApartmentLink `json:"apartments"`
}

type Member struct {
	UserInfo   UserInfo   `json:"userInfo"`
	MemberInfo MemberInfo `json:"memberInfo"`
}

type UserInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MemberInfo struct {
	AddedDate  int64  `json:"addedDate"`
	MemberType string `json:"memberType"`
	Role       string `json:"role"`
}

type ApartmentLink struct {
	Apartment  Apartment `json:"apartment"`
	MemberType string    `json:"memberType"`
}

type Apartment struct {
	ID               int    `json:"id"`
	FlatNo           string `json:"flatNo"`
	AllocationStatus string `json:"allocationStatus"`
	Tags             []any  `json:"tags"`
}

// These structs are what the MCP tools send back.
type GetMembersResponse struct {
	Tool        string          `json:"tool"`
	SocietyID   int             `json:"societyId"`
	MemberCount int             `json:"memberCount"`
	Members     []MemberDetails `json:"members"`
}

type MemberDetails struct {
	ID              int          `json:"id"`
	Name            string       `json:"name"`
	Flats           []MemberFlat `json:"flats"`
	OccupancyStatus string       `json:"occupancyStatus"`
	Role            string       `json:"role"`
	PendingDues     int          `json:"pendingDues"`
	Tags            []string     `json:"tags"`
}

type MemberFlat struct {
	FlatID           int    `json:"flatId"`
	FlatNumber       string `json:"flatNumber"`
	RelationshipType string `json:"relationshipType"`
}

type GetMemberByIDResponse struct {
	Tool      string     `json:"tool"`
	Status    string     `json:"status"`
	SocietyID int        `json:"societyId"`
	Member    MemberByID `json:"member,omitempty"`
	MemberID  int        `json:"memberId,omitempty"`
}

type MemberByID struct {
	ID         int          `json:"id"`
	Name       string       `json:"name"`
	Flats      []SimpleFlat `json:"flats"`
	MemberType string       `json:"memberType"`
	JoinedOn   string       `json:"joinedOn"`
}

type SimpleFlat struct {
	FlatID           int    `json:"flatId"`
	RelationshipType string `json:"relationshipType"`
}

type GetLedgersResponse struct {
	Tool      string `json:"tool"`
	Status    string `json:"status"`
	SocietyID int    `json:"societyId"`
	Data      any    `json:"data"`
}
