package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.schibsted.io/smmx/golin/boltdb"
	"github.schibsted.io/smmx/golin/config"
	"github.schibsted.io/smmx/golin/login"
	"github.schibsted.io/smmx/golin/tokens"
	"log"
	"regexp"
	"time"
)

// BoltDB was used because of it's speed and architecture goes
// very well accordingly into what we want to achieve, which is
// a secure and fast service to consume tokens.
var (
	boltWrite *bolt.DB
	sqlDB     *sql.DB
	cfg       *config.Configuration
)

type User struct {
	AccountId string `json:"account_id,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}

type Token struct {
	Token string
}

type Account struct {
	Email    string
	Password string
}

type Claim struct {
	Id  string
	exp int64
}

func init() {
	var err error
	cfg, err = config.ReadConfig("config/config.conf")
	fmt.Println("Reading conf:", cfg.AccountsDB, cfg.Schema, cfg.Keys)
	if err != nil {
		panic("Failed to init due to configurations")
	}

	boltWrite, err = boltdb.OpenBoltDB("tokens")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	dbconn := fmt.Sprintf("host=%s port=%d dbname=%s sslmode=disable", cfg.AccountsDB.Host, cfg.AccountsDB.Port, cfg.AccountsDB.DBName)
	sqlDB, err = sql.Open("postgres", dbconn)
	if err != nil {
		fmt.Println("Error opening DB: ", err)
		fmt.Println(dbconn)
		log.Fatal(err)
	}
	fmt.Println("Init finished")
}

// Use config file to decide which DB you'll be using for storing the tokens
func main() {
	r := gin.Default()
	publics := r.Group("api/v1/public")
	privates := r.Group("api/v1/private")
	privates.GET("/users/:id", GetUser)
	publics.POST("/users", LoginUser)
	r.Run()
}

// These are the endpoints required to do a login and verifying that tokens are
// alive
func LoginUser(c *gin.Context) {
	var user User
	var loginer login.Loginer
	var tokenStr string

	var token Token
	var tokener tokens.Tokener

	c.BindJSON(&user)

	loginer = user
	tokener = token

	SignatureStr, err := loginer.Login(user.Email, user.Password)
	tokenStr, err = tokener.GenerateToken(SignatureStr)

	if err != nil {
		c.JSON(404, gin.H{"error generating token": err})
	} else {
		// Here the token is sent to BoltDB
		data := structs.Map(user)
		err = boltdb.UpdateBucket(boltWrite, tokenStr, data)
		if err != nil {
			c.JSON(404, gin.H{"error updating bucket": err})
		} else {
			c.JSON(201, gin.H{"token": tokenStr, "email": user.Email})
		}
	}
}

// Here we use the function to validate if a token exists and if the user should
// alllowed to enter and go through
func GetUser(c *gin.Context) {
	_ = c.Params.ByName("id")
	var currentToken Token

	if len(c.Request.Header["Authorization"]) > 0 {
		currentToken.Token = string(c.Request.Header["Authorization"][0])
	}

	newToken, err := currentToken.ValidateToken(currentToken.Token)

	if err != nil {
		c.JSON(403, gin.H{"error": err})
	} else {
		c.JSON(200, gin.H{"message": "chingon perron", "token": newToken})
	}
}

func (t Token) ValidateToken(encriptedToken string) (string, error) {
	// TODO: Blacklist mechanism of logout
	tokData := regexp.MustCompile(`\s*$`).ReplaceAll([]byte(encriptedToken), []byte{})

	currentToken, err := jwt.Parse(string(tokData), func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.Keys.Secret), nil
	})

	// Print an error if we can't parse for some reason
	if err != nil {
		fmt.Println("Couldn't parse token: ", err)
		return "", err
	}

	fmt.Println("currentToken ", currentToken)
	// Is token invalid?
	if !currentToken.Valid {
		return "", fmt.Errorf("Token is invalid")
	}

	//Validate hasn't expired
	fmt.Println(currentToken)
	email := currentToken.Claims["Id"].(string)
	if Expired(currentToken.Raw, email) {
		return "", fmt.Errorf("Token has expired")
	}

	// Print the token details
	_, err = json.MarshalIndent(currentToken.Claims, "", "    ")

	return string(tokData), err
}

func Expired(token, email string) bool {
	some, err := boltdb.GetEmailValue(boltWrite, token)
	if err != nil {
		fmt.Println(err)
		return true
	}
	return !(some == email)
}

// This is where config file should be used to read and compare users
// since this is the MVP of this microservice, this works for achieving
// what we want.
func (u User) Login(user, password string) (string, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE %s = '%s' `, cfg.Schema.Table, cfg.Schema.Email, user)
	fmt.Println("Query: ", query)
	rows, err := sqlDB.Query(query)
	if err != nil {
		rows.Close()
		return "", err
	}

	var account Account
	for rows.Next() {
		rows.Scan(&account.Email, &account.Password)
		fmt.Println("Query Result:", account.Email, account.Password)
	}
	rows.Close()
	return account.Email, err
}

// JWT way to create and generate tokens
func (t Token) GenerateToken(SignatureStr string) (string, error) {
	var claim Claim

	claim.Id = SignatureStr
	claim.exp = time.Now().Add(time.Hour * 1).Unix()

	alg := jwt.GetSigningMethod("HS256")
	token := jwt.New(alg)
	token.Claims = structs.Map(claim)

	if tokenStr, err := token.SignedString([]byte(cfg.Keys.Secret)); err == nil {
		return tokenStr, nil
	} else {
		return "", err
	}
}
