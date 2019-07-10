package server

import (
	ipfs "github.com/Varunram/essentials/ipfs"
	erpc "github.com/Varunram/essentials/rpc"
	"github.com/YaleOpenLab/openclimate/database"
	"log"
	"math/big"
	"net/http"
)

func setupDBHandlers() {
	newUser()
	retrieveUser()
	retrieveAllUsers()
	deleteUser()
	updateUser()
	getIpfsHash()
	getAllCompanies()
	getCompany()
	getAllRegions()
	getRegion()
}

/*****************/
/* USER HANDLERS */
/*****************/

// setupPingHandler is a ping route for remote callers to check if the platform is up
func newUser() {
	http.HandleFunc("/user/new", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if r.URL.Query()["username"] == nil || r.URL.Query()["pwhash"] == nil || r.URL.Query()["email"] == nil {
			log.Println("required params - username, pwhash, email missing")
			log.Println(r.URL.Query()["username"])
			log.Println(r.URL.Query()["pwhash"])
			log.Println(r.URL.Query()["email"])
			responseHandler(w, StatusBadRequest)
			return
		}

		username := r.URL.Query()["username"][0]
		pwhash := r.URL.Query()["pwhash"][0]
		email := r.URL.Query()["email"][0]

		user, err := database.NewUser(username, pwhash, email)
		if err != nil {
			log.Println("couldn't create new user", err)
			responseHandler(w, StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, user)
	})
}

func authorizeUser(r *http.Request) (database.User, error) {
	username := r.URL.Query()["username"][0]
	pwhash := r.URL.Query()["pwhash"][0]

	return database.ValidateUser(username, pwhash)
}

func retrieveUser() {
	http.HandleFunc("/user/retrieve", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		user, err := authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quittting")
			responseHandler(w, StatusBadRequest)
			return
		}

		erpc.MarshalSend(w, user)
	})
}

func retrieveAllUsers() {
	http.HandleFunc("/user/retrieve/all", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		_, err = authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quittting")
			responseHandler(w, StatusInternalServerError)
			return
		}

		users, err := database.RetrieveAllUsers()
		if err != nil {
			log.Println("could not retrieve user from the database, quittting")
			responseHandler(w, StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, users)
	})
}

func deleteUser() {
	http.HandleFunc("/user/delete", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		user, err := authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quittting")
			responseHandler(w, StatusInternalServerError)
			return
		}

		err = database.DeleteKeyFromBucket(user.Index, database.UserBucket)
		if err != nil {
			log.Println("could not delete user from database, quittting", err)
			responseHandler(w, StatusBadRequest)
			return
		}

		responseHandler(w, StatusOK)
	})
}

func updateUser() {
	http.HandleFunc("/user/update", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		user, err := authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quittting")
			responseHandler(w, StatusInternalServerError)
			return
		}

		if r.URL.Query()["email"] != nil {
			user.Email = r.URL.Query()["email"][0]
		} else if r.URL.Query()["newpwhash"] != nil {
			user.Pwhash = r.URL.Query()["newpwhash"][0]
		} else if r.URL.Query()["newusername"] != nil {
			user.Name = r.URL.Query()["newusername"][0]
		} else {
			responseHandler(w, StatusBadRequest)
			return
		}

		err = user.Save()
		if err != nil {
			log.Println("error while savingt user to database")
			responseHandler(w, StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, user)
	})
}

/*****************/
/* IPFS HANDLERS */
/*****************/

// getIpfsHash gets the ipfs hash of the passed string
func getIpfsHash() {
	http.HandleFunc("/ipfs/hash", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		_, err = authorizeUser(r)
		if err != nil {
			responseHandler(w, StatusUnauthorized)
			return
		}
		if r.URL.Query()["string"] == nil {
			responseHandler(w, StatusBadRequest)
			return
		}

		hashString := r.URL.Query()["string"][0]
		hash, err := ipfs.AddStringToIpfs(hashString)
		if err != nil {
			log.Println("did not add string to ipfs", err)
			responseHandler(w, StatusInternalServerError)
			return
		}

		hashCheck, err := ipfs.GetStringFromIpfs(hash)
		if err != nil || hashCheck != hashString {
			responseHandler(w, StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, hash)
	})
}

/*********************/
/* ETHEREUM HANDLERS */
/*********************/

func sendEth() {
	http.HandleFunc("/user/sendeth", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		user, err := authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quittting")
			responseHandler(w, StatusBadRequest)
			return
		}

		if r.URL.Query()["address"] == nil || r.URL.Query()["amount"] == nil {
			log.Println("address or amount missing, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}

		address := r.URL.Query()["address"][0]
		amountStr := r.URL.Query()["amount"][0] // convert this to bigint

		var amount big.Int
		_, boolErr := amount.SetString(amountStr, 10)
		if !boolErr {
			responseHandler(w, StatusInternalServerError)
			return
		}

		txhash, err := user.SendEthereumTx(address, amount)
		if err != nil {
			responseHandler(w, StatusInternalServerError)
			return
		}

		log.Println("user: ", user.Name, "has sent tx with txhash: ", txhash)
		responseHandler(w, StatusOK)
	})
}

/*******************/
/* REGION HANDLERS */
/*******************/

func getAllRegions() {
	http.HandleFunc("/region/all", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		_, err = authorizeUser(r)
		if err != nil {
			log.Println("Could not retrieve user from the database, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}

		regions, err := database.RetrieveAllRegions()
		if err != nil {
			log.Println("Error while retrieving all companies, quitting")
			responseHandler(w, StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, regions)
	})
}

func getRegion() {
	http.HandleFunc("/region", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		_, err = authorizeUser(r)
		if err != nil {
			log.Println("Could not retrieve user from the database, quitting")
			responseHandler(w, StatusBadRequest)
		}

		if r.URL.Query()["region_name"] == nil || r.URL.Query()["region_country"] == nil {
			log.Println("Region_name or region_country not passed, quitting")
			responseHandler(w, StatusBadRequest)
		}

		name := r.URL.Query()["region_name"][0]
		country := r.URL.Query()["region_country"][0]
		region, err := database.RetrieveRegionByName(name, country) //************ STOP ***********
		if err != nil {
			log.Println("Error while retrieving all companies, quitting")
			responseHandler(w, StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, region)
	})
}

/********************/
/* COMPANY HANDLERS */
/********************/

func getAllCompanies() {
	http.HandleFunc("/companies/all", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		_, err = authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quittting")
			responseHandler(w, StatusBadRequest)
			return
		}

		companies, err := database.RetrieveAllCompanies()
		if err != nil {
			log.Println("error while retrieving all companies, quitting")
			responseHandler(w, StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, companies)
	})
}

func getCompany() {
	http.HandleFunc("/company", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		_, err = authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}

		if r.URL.Query()["company"] == nil {
			log.Println("company name not passed, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}

		companyName := r.URL.Query()["company"][0]
		company, err := database.RetrieveCompanyByName(companyName)
		if err != nil {
			log.Println("error while retrieving all companies, quitting")
			responseHandler(w, StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, company)
	})
}
