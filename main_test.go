package main

import (
	"reflect"
	"testing"
)

var (
	testSearchList    []string
	testSearchStrings []string
)

func init() {
	fileLoc = "./test-assets/test-in.txt"
	testSearchList = []string{
		"Deoksugung Palace",
		"National Museum of Korea",
		"Jade Emperor Pagoda",
		"War Remnants Museum Ho Chi Minh City",
		"HCMC Museum",
		"Central Post Office Ho Chi Minh City",
		"Hoan Kiem Lake",
		"Temple of Literature Hanoi",
		"Imperial Citadel of Thang Long",
		"National Museum of Vietnamese History Hanoi",
		"Bach Ma Temple",
		"Vietnam Military History Museum Hanoi",
		"Dong Xuan Market",
		"Martyrs' Monument Hanoi",
		"Tran Quoc Pagoda",
		"Fine Arts Museum of Vietnam",
		"Vietnamese Women’s Museum Hanoi",
		"Đường Tàu",
		"Bun Cha Huong Lien(Restuarante)",
		"Ho Chi Minh Mausoleum",
		"MF Spa Hanoi",
	}
	testSearchStrings = []string{
		"https://www.google.com/maps/search/Deoksugung+Palace",
		"https://www.google.com/maps/search/National+Museum+of+Korea",
		"https://www.google.com/maps/search/Jade+Emperor+Pagoda",
		"https://www.google.com/maps/search/War+Remnants+Museum+Ho+Chi+Minh+City",
		"https://www.google.com/maps/search/HCMC+Museum",
		"https://www.google.com/maps/search/Central+Post+Office+Ho+Chi+Minh+City",
		"https://www.google.com/maps/search/Hoan+Kiem+Lake",
		"https://www.google.com/maps/search/Temple+of+Literature+Hanoi",
		"https://www.google.com/maps/search/Imperial+Citadel+of+Thang+Long",
		"https://www.google.com/maps/search/National+Museum+of+Vietnamese+History+Hanoi",
		"https://www.google.com/maps/search/Bach+Ma+Temple",
		"https://www.google.com/maps/search/Vietnam+Military+History+Museum+Hanoi",
		"https://www.google.com/maps/search/Dong+Xuan+Market",
		"https://www.google.com/maps/search/Martyrs'+Monument+Hanoi",
		"https://www.google.com/maps/search/Tran+Quoc+Pagoda",
		"https://www.google.com/maps/search/Fine+Arts+Museum+of+Vietnam",
		"https://www.google.com/maps/search/Vietnamese+Women’s+Museum+Hanoi",
		"https://www.google.com/maps/search/Đường+Tàu",
		"https://www.google.com/maps/search/Bun+Cha+Huong+Lien(Restuarante)",
		"https://www.google.com/maps/search/Ho+Chi+Minh+Mausoleum",
		"https://www.google.com/maps/search/MF+Spa+Hanoi",
	}

}

// Unit test for getSearchlist
func TestGetSearchList(t *testing.T) {
	err := getSearchList()
	if !reflect.DeepEqual(searchList, testSearchList) || err != nil {
		t.Error("Failed to get correct search items: ", err)
	}
}

// Unit test for createSearchStrings
func TestCreateSearchStrings(t *testing.T) {
	err := getSearchList()
	if err != nil {
		t.Error("Error in getSearchList: ", err)
	}

	createSearchStrings()

	for index, url := range urlArr {
		shouldBe := testSearchStrings[index]
		if url != shouldBe {
			t.Errorf("%v is incorrect, should be: %v", url, shouldBe)
			break
		}
	}
}
