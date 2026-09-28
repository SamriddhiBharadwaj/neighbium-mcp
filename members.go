package main

import (
	"fmt"
	"strings"
	"time"
)

func buildGetMembersResponse(societyID int, entries []DirectoryEntry) GetMembersResponse {
	members := []MemberDetails{}

	for _, entry := range entries {
		members = append(members, toMemberDetails(entry))
	}

	return GetMembersResponse{
		Tool:        "get_member_details",
		SocietyID:   societyID,
		MemberCount: len(members),
		Members:     members,
	}
}

func toMemberDetails(entry DirectoryEntry) MemberDetails {
	flats := []MemberFlat{}
	allTags := []string{}
	occupied := false

	for _, link := range entry.Apartments {
		flat := MemberFlat{
			FlatID:           link.Apartment.ID,
			FlatNumber:       link.Apartment.FlatNo,
			RelationshipType: strings.ToUpper(link.MemberType),
		}
		flats = append(flats, flat)

		if link.Apartment.AllocationStatus != "Vacant" {
			occupied = true
		}

		for _, tag := range link.Apartment.Tags {
			allTags = append(allTags, fmt.Sprint(tag))
		}
	}

	role := entry.Member.MemberInfo.Role
	if role == "" {
		role = entry.Member.MemberInfo.MemberType
	}

	occupancy := "Vacant"
	if occupied {
		occupancy = "Occupied"
	}

	return MemberDetails{
		ID:              entry.Member.UserInfo.ID,
		Name:            entry.Member.UserInfo.Name,
		Flats:           flats,
		OccupancyStatus: occupancy,
		Role:            role,
		PendingDues:     0,
		Tags:            uniqueStrings(allTags),
	}
}

func findMemberByID(entries []DirectoryEntry, memberID int) (DirectoryEntry, bool) {
	for _, entry := range entries {
		if entry.Member.UserInfo.ID == memberID {
			return entry, true
		}
	}
	return DirectoryEntry{}, false
}

func buildGetMemberByIDResponse(societyID int, entry DirectoryEntry) GetMemberByIDResponse {
	flats := []SimpleFlat{}
	for _, link := range entry.Apartments {
		flats = append(flats, SimpleFlat{
			FlatID:           link.Apartment.ID,
			RelationshipType: strings.ToUpper(link.MemberType),
		})
	}

	memberType := entry.Member.MemberInfo.MemberType
	if len(entry.Apartments) > 0 {
		memberType = entry.Apartments[0].MemberType
	}

	joinedOn := time.UnixMilli(entry.Member.MemberInfo.AddedDate).UTC().Format("2006-01-02")

	return GetMemberByIDResponse{
		Tool:      "get_member_by_id",
		Status:    "success",
		SocietyID: societyID,
		Member: MemberByID{
			ID:         entry.Member.UserInfo.ID,
			Name:       entry.Member.UserInfo.Name,
			Flats:      flats,
			MemberType: memberType,
			JoinedOn:   joinedOn,
		},
	}
}

func uniqueStrings(values []string) []string {
	result := []string{}

	for _, value := range values {
		alreadyAdded := false
		for _, existing := range result {
			if existing == value {
				alreadyAdded = true
				break
			}
		}
		if !alreadyAdded {
			result = append(result, value)
		}
	}

	return result
}
