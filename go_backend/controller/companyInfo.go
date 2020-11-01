package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zijianguan0204/hunger_at_home/model"
)

//Company Signup
// @Summary Add a new company
// @Description get company info and parse them to db
// @Tags company
// @Accept  json
// @Produce  json
// @Param companyInfo body model.CompanySignupInput true "Add a company"
// @Success 200 {object} model.SignupOutput
// @Failure 500 {object} model.Message
// @Failure 404 {object} model.Message
// @Router /company/companySignup/ [post]
func (ct *Controller) QueryCompanySignUp(c *gin.Context) {
	//getting data from request
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var signup model.CompanySignupInput
	er := json.Unmarshal([]byte(value), &signup)
	if er != nil {
		fmt.Println(err.Error())
	}
	e := reflect.ValueOf(&signup).Elem()
	companyName := fmt.Sprint(e.Field(0).Interface())
	fedID := fmt.Sprint(e.Field(1).Interface())
	einID := fmt.Sprint(e.Field(2).Interface())

	// connect to the db
	db, err := connectDB(c)
	if err != nil {
		message := model.Message{
			Message: "DB has something wrong",
		}
		c.JSON(500, message)
		panic(err.Error())
	}

	//insert to db
	insert, err := db.Query("INSERT INTO company VALUES(DEFAULT,?,?,?,DEFAULT)", companyName, fedID, einID)

	if err != nil {
		fmt.Println("Sign up error")
		if strings.Contains(err.Error(), "Access denied") {
			message := model.Message{
				Message: "DB access error, username or password is wrong",
			}
			c.JSON(500, message)
			panic(err.Error())
		}
		message := model.Message{
			Message: "Company name, FED or EIN is/are already taken",
		}
		c.JSON(500, message)
		panic(err.Error())
	} else {
		getCompany(companyName, c, db)
	}
	defer db.Close()
	defer insert.Close()
}

//Address Signup
// @Summary Add a new address
// @Description get address info and parse them to db
// @Tags company
// @Accept  json
// @Produce  json
// @Param companyInfo body model.AddressSignupInput true "Add an address"
// @Success 201 {object} model.Message
// @Failure 500 {object} model.Message
// @Failure 404 {object} model.Message
// @Router /company/addressSignUp/ [post]
func (ct *Controller) QueryAddressSignUp(c *gin.Context) {
	//getting data from request
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var signup model.AddressSignupInput
	er := json.Unmarshal([]byte(value), &signup)
	if er != nil {
		fmt.Println(err.Error())
	}
	e := reflect.ValueOf(&signup).Elem()
	address := fmt.Sprint(e.Field(0).Interface())
	city := fmt.Sprint(e.Field(1).Interface())
	state := fmt.Sprint(e.Field(2).Interface())
	zipCode := fmt.Sprint(e.Field(3).Interface())

	// connect to the db
	db, err := connectDB(c)
	if err != nil {
		message := model.Message{
			Message: "DB has something wrong",
		}
		c.JSON(500, message)
		panic(err.Error())
	}

	defer db.Close()
	//insert to db
	insert, err := db.Query("INSERT INTO address VALUES(DEFAULT,?,?,?,?,DEFAULT)", address, city, state, zipCode)
	if err != nil {
		if strings.Contains(err.Error(), "Access denied") {
			message := model.Message{
				Message: "DB access error, username or password is wrong",
			}
			c.JSON(500, message)
			panic(err.Error())
		}
		//TODO:check the address is taken or not
		message := model.Message{
			Message: "This address is already taken",
		}
		c.JSON(500, message)
		panic(err.Error())
	} else {
		message := model.Message{
			Message: "This address has been assigned or created",
		}
		c.JSON(201, message)
	}
	defer insert.Close()
}

// Get address list
// @Summary Get company adderss list based on Company ID
// @Description CompanyID as input then list of address in JSON as output
// @Tags company
// @Accept  json
// @Produce  json
// @Param companyInfo body model.AddressListInput true "get an address list by its companyID"
// @Success 200 {array} model.AddressListOutput
// @Failure 500 {object} model.Message
// @Failure 404 {object} model.Message
// @Router /company/addressList/ [post]
func (ct *Controller) QueryGetAddressList(c *gin.Context) {
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var input model.AddressListInput
	er := json.Unmarshal([]byte(value), &input)
	if er != nil {
		fmt.Println(err.Error())
	}
	e := reflect.ValueOf(&input).Elem()
	companyID := fmt.Sprint(e.Field(0).Interface())
	fmt.Println(companyID)
	//connect to DB
	db, err := connectDB(c)
	if err != nil {
		fmt.Println("DB error")
		c.JSON(500, gin.H{
			"message": "DB connection problem",
		})
		panic(err.Error())
	}
	defer db.Close()
	//get company names
	var (
		address string
		city    string
		state   string
		zipCode string
	)

	addressList := []model.AddressListOutput{}

	rows, err := db.Query("SELECT address.address, address.city, address.state, address.zipCode FROM address INNER JOIN companyAddressAssociate ON address.id = companyAddressAssociate.addressID WHERE companyAddressAssociate.companyID = ?", companyID)
	if err != nil {
		if strings.Contains(err.Error(), "Access denied") {
			c.JSON(500, gin.H{
				"message": "DB access error, username or password is wrong",
			})
			panic(err.Error())
		}
		c.JSON(500, gin.H{
			"message": "Query statements has something wrong",
		})
		panic(err.Error())
	} else {
		for rows.Next() {
			err := rows.Scan(&address, &city, &state, &zipCode)
			if err != nil {
				log.Fatal(err)
			}
			addressRecord := model.AddressListOutput{
				Address: address,
				City:    city,
				State:   state,
				ZipCode: zipCode,
			}
			addressList = append(addressList, addressRecord)
		}
	}

	c.JSON(200, addressList)
	defer rows.Close()

}

// Get company list
// @Summary Get the whole company list with their IDs
// @Description Return a JSON formated array of company
// @Tags company
// @Accept  json
// @Produce  json
// @Success 200 {array} model.CompanyListOutput
// @Failure 500 {object} model.Message
// @Failure 404 {object} model.Message
// @Router /company/companyList/ [post]
func (ct *Controller) QueryGetCompanyList(c *gin.Context) {
	//connect to DB
	db, err := connectDB(c)
	if err != nil {
		fmt.Println("DB error")
		c.JSON(500, gin.H{
			"message": "DB connection problem",
		})
		panic(err.Error())
	}
	defer db.Close()
	//get company names
	var (
		id          string
		companyName string
	)
	companyList := []model.CompanyListOutput{}

	rows, err := db.Query("SELECT id, companyName from company")
	if err != nil {
		if strings.Contains(err.Error(), "Access denied") {
			message := model.Message{
				Message: "DB access error, username or password is wrong",
			}
			c.JSON(500, message)
			panic(err.Error())
		}
		message := model.Message{
			Message: "Query statements has something wrong",
		}
		c.JSON(500, message)
		panic(err.Error())
	} else {
		for rows.Next() {
			err := rows.Scan(&id, &companyName)
			if err != nil {
				log.Fatal(err)
			}
			company := model.CompanyListOutput{
				CompanyID:   id,
				CompanyName: companyName,
			}
			companyList = append(companyList, company)
		}
	}

	c.JSON(200, companyList)
	defer rows.Close()

}

// Add Company&Address associate record
// @Summary Add a new record to company address associate table
// @Description get companyAddress associate record info and parse them to db
// @Tags company
// @Accept  json
// @Produce  json
// @Param companyInfo body model.CompanyAddressRecordSignupInput true "Add a record"
// @Success 201 {object} model.Message
// @Failure 500 {object} model.Message
// @Failure 404 {object} model.Message
// @Router /company/addressCompanyAssociate/ [post]
func (ct *Controller) QueryAddCompanyAddressAssociate(c *gin.Context) {
	//getting data from request
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var companyAddress model.CompanyAddressRecordSignupInput
	er := json.Unmarshal([]byte(value), &companyAddress)
	if er != nil {
		fmt.Println(err.Error())
	}
	e := reflect.ValueOf(&companyAddress).Elem()
	addressID := fmt.Sprint(e.Field(0).Interface())
	companyID := fmt.Sprint(e.Field(1).Interface())

	fmt.Println(addressID)
	fmt.Println(companyID)

	// connect to the db
	db, err := connectDB(c)
	if err != nil {
		message := model.Message{
			Message: "DB has something wrong",
		}
		c.JSON(500, message)
		panic(err.Error())
	}

	defer db.Close()
	//insert to db
	insert, err := db.Query("INSERT INTO companyAddressAssociate VALUES(?,?,DEFAULT)", addressID, companyID)
	if err != nil {
		fmt.Println("This record can not be submitted")
		if strings.Contains(err.Error(), "Access denied") {
			message := model.Message{
				Message: "DB access error, username or password is wrong",
			}
			c.JSON(500, message)
			panic(err.Error())
		}
		message := model.Message{
			Message: "This record can not be submitted",
		}
		c.JSON(500, message)
		panic(err.Error())
	} else {
		message := model.Message{
			Message: "This record has been assigned or created",
		}
		c.JSON(201, message)
	}
	defer insert.Close()
}