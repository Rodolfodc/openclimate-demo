package database

import (
	"github.com/Varunram/essentials/utils"
	globals "github.com/YaleOpenLab/openclimate/globals"
	"github.com/pkg/errors"
)

type Actor interface {
	GetPledges() ([]Pledge, error)
	AddPledges(pledgeIDs ...int) error
	UpdateMRV(MRV string)
	GetID() int
}

type BucketItem interface {
	SetID(id int)
	GetID() int
	Save() error
}

type Location struct {
	Name            string
	Latitude        string
	Longitude       string
	NationStaeID    int
	RegionalStateID int
}

type DistributionRecord struct {
	Title  string
	Weight int
}

// type RepData struct {
// 	// emissions, mitigation, adaption, etc.
// 	ReportType string
// 	Year       int
// 	IpfsHash   string
// }

/*
	Given the type of actor (company, city, state, region, country, etc.) and
	the ID of the actor, return the entity (all actor types implement the
	Actor interface, so the function returns the interface).
*/
func RetrieveActor(actorType string, actorID int) (Actor, error) {

	var actor Actor
	var err error

	switch actorType {
	case "company":
		var x Company
		x, err = RetrieveCompany(actorID)
		actor = &x
	case "city":
		var x City
		x, err = RetrieveCity(actorID)
		actor = &x
	case "state":
		var x State
		x, err = RetrieveState(actorID)
		actor = &x
	case "region":
		var x Region
		x, err = RetrieveRegion(actorID)
		actor = &x
	case "country":
		var x Country
		x, err = RetrieveCountry(actorID)
		actor = &x
	case "oversight":
		var x Oversight
		x, err = RetrieveOsOrg(actorID)
		actor = &x
	default:
		return actor, errors.New("User's actor type is not valid.")
	}

	if err != nil {
		return actor, nil
	}

	return actor, nil
}

// Puts asset object in assets bucket. Called by NewAsset
func (x *Asset) Save() error {
	return Save(globals.DbPath, AssetBucket, x)
}

// Saves city object in cities bucket. Called by NewCity
func (x *City) Save() error {
	x.LastUpdated = utils.Timestamp()
	return Save(globals.DbPath, CityBucket, x)
}

// Saves country object in countries bucket. Called by NewCountry
func (x *Country) Save() error {
	x.LastUpdated = utils.Timestamp()
	return Save(globals.DbPath, CountryBucket, x)
}

func (x *Oversight) Save() error {
	return Save(globals.DbPath, OversightBucket, x)
}

func (x *Pledge) Save() error {
	return Save(globals.DbPath, PledgeBucket, x)
}

// Saves region object in regions bucket. Called by NewRegion
func (x *Region) Save() error {
	x.LastUpdated = utils.Timestamp()
	return Save(globals.DbPath, RegionBucket, x)
}

func (x *ConnectRequest) Save() error {
	return Save(globals.DbPath, RequestBucket, x)
}

// Saves state object in states bucket. Called by NewState
func (x *State) Save() error {
	x.LastUpdated = utils.Timestamp()
	return Save(globals.DbPath, StateBucket, x)
}

// Save inserts a passed User object into the database
func (x *User) Save() error {
	return Save(globals.DbPath, UserBucket, x)
}

// Saves company object in companies bucket. Called by NewCompany
func (x *Company) Save() error {
	x.LastUpdated = utils.Timestamp()
	return Save(globals.DbPath, CompanyBucket, x)
}

/* 	BucketItem interface method:

SetID() is a common method between all structs that qualify as
bucket items that allow them to implement the BucketItem interface.
SetID() is a simple setter method that allows the updating of the
bucket item's ID. The function's only use should be in the Save()
function; otherwise, the ID should not be modified.
*/
func (x *Company) SetID(id int) {
	x.Index = id
}

func (x *Asset) SetID(id int) {
	x.Index = id
}

func (x *City) SetID(id int) {
	x.Index = id
}

func (x *Country) SetID(id int) {
	x.Index = id
}

func (x *Oversight) SetID(id int) {
	x.Index = id
}

func (x *Pledge) SetID(id int) {
	x.ID = id
}

func (x *Region) SetID(id int) {
	x.Index = id
}

func (x *ConnectRequest) SetID(id int) {
	x.Index = id
}

func (x *State) SetID(id int) {
	x.Index = id
}

func (x *User) SetID(id int) {
	x.Index = id
}

/* 	BucketItem interface method:

A getter method for structs that implement the BucketItem interface.
The method allows you to retrieve the ID of structs that implement the
BucketItem interface methods, even if you don't know specifically which
struct you are workign with.
*/
func (x *Company) GetID() int {
	return x.Index
}

func (x *Asset) GetID() int {
	return x.Index
}

func (x *City) GetID() int {
	return x.Index
}

func (x *Country) GetID() int {
	return x.Index
}

func (x *Oversight) GetID() int {
	return x.Index
}

func (x *Pledge) GetID() int {
	return x.ID
}

func (x *Region) GetID() int {
	return x.Index
}

func (x *ConnectRequest) GetID() int {
	return x.Index
}

func (x *State) GetID() int {
	return x.Index
}

func (x *User) GetID() int {
	return x.Index
}

/*	Actor Interface method:

	Allows for the updating of the chosen reporting methodology
	for any of the climate actor types that implement the
	Actor interface.
*/
func (c *Company) UpdateMRV(MRV string) {
	c.MRV = MRV
	c.Save()
}

func (c *City) UpdateMRV(MRV string) {
	c.MRV = MRV
	c.Save()
}

func (c *Country) UpdateMRV(MRV string) {
	c.MRV = MRV
	c.Save()
}

func (c *State) UpdateMRV(MRV string) {
	c.MRV = MRV
	c.Save()
}

func (c *Region) UpdateMRV(MRV string) {
	c.MRV = MRV
	c.Save()
}

func (c *Oversight) UpdateMRV(MRV string) {
	c.MRV = MRV
	c.Save()
}

func SearchState(name string) ([]State, error) {
	var arr []State
	all, err := RetrieveAllStates()
	if err != nil {
		return arr, errors.Wrap(err, "Error while retrieving all cities from database")
	}

	for _, val := range all {
		if val.Name == name {
			arr = append(arr, val)
		}
	}

	return arr, nil
}

func SearchCity(name string) ([]City, error) {
	var arr []City
	all, err := RetrieveAllCities()
	if err != nil {
		return arr, errors.Wrap(err, "Error while retrieving all cities from database")
	}

	for _, val := range all {
		if val.Name == name {
			arr = append(arr, val)
		}
	}

	return arr, nil
}

func SearchRegion(name string) ([]Region, error) {
	var arr []Region
	all, err := RetrieveAllRegions()
	if err != nil {
		return arr, errors.Wrap(err, "Error while retrieving all regions from database")
	}

	for _, val := range all {
		if val.Name == name {
			arr = append(arr, val)
		}
	}

	return arr, nil
}

func SearchCompany(name string) ([]Company, error) {
	var arr []Company
	all, err := RetrieveAllCompanies()
	if err != nil {
		return arr, errors.Wrap(err, "Error while retrieving all companies from database")
	}

	for _, val := range all {
		if val.Name == name {
			arr = append(arr, val)
		}
	}

	return arr, nil
}
