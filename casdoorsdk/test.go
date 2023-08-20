package casdoorsdk

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/casdoor/casdoor-go-sdk/conf"
	"github.com/joho/godotenv"
)

//go:embed token_jwt_key.pem
var JwtPublicKey string
var localClient *Client

func InitConfigTest() {
	err := godotenv.Load("../conf/app.conf")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	casdoorEndpoint := strings.TrimRight(conf.GetConfigString("casdoorEndpoint"), "/")
	clientId := conf.GetConfigString("clientId")
	clientSecret := conf.GetConfigString("clientSecret")
	casdoorOrganization := conf.GetConfigString("casdoorOrganization")
	casdoorApplication := conf.GetConfigString("casdoorApplication")

	localClient = NewClient(casdoorEndpoint, clientId, clientSecret, JwtPublicKey, casdoorOrganization, casdoorApplication)
}
