package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bndr/gotabulate"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

var (
	clientID       string
	clientSecret   string
	userID         string
	port           string
	env            string
	authServerURL  string
	tokenServerURL string
	sourceSeverURL string
	expectedState  string
	codeVerifier   string
)

func init() {

	flag.StringVar(&clientID, "i", "", "client id")

	flag.StringVar(&clientSecret, "s", "", "client secret")

	flag.StringVar(&userID, "u", "", "user id")

	flag.StringVar(&port, "p", "9094", "client serve port")

	flag.StringVar(&env, "e", "stage", "BITGIN environment mode e.g. stage, prod")

	flag.Usage = usage

	flag.Parse()

	if len(clientID) == 0 {
		panic("client id is empty")
	}

	if len(clientSecret) == 0 {
		panic("client secret is empty")
	}

	if len(userID) == 0 {
		panic("user id is empty")
	}

	switch env {
	case "stage":
		authServerURL = "https://stage.bitgin.app"
		tokenServerURL = "https://oauth.bitgin.app"
	case "prod":
		authServerURL = "https://bitgin.net"
		tokenServerURL = "https://oauth.bitgin.net"
	default:
		panic("unknown environment")
	}

	sourceSeverURL = tokenServerURL
	expectedState = "xyz"
	codeVerifier = "s256example"
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: oa2cli [options] \n")
	fmt.Fprintf(os.Stderr, "  Currently, the following flags can be used\n")
	flag.PrintDefaults()
}

var (
	globalToken *oauth2.Token // Non-concurrent security
)

func main() {

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   authServerURL + "/v1/oauth/authorize",
			TokenURL:  tokenServerURL + "/v1/oauth/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		u := config.AuthCodeURL(expectedState,
			oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256(codeVerifier)),
			oauth2.SetAuthURLParam("user_id", userID),
		)
		http.Redirect(c.Response().Writer, c.Request(), u, http.StatusFound)

		return nil
	})

	e.GET("/oauth", func(c echo.Context) error {

		if err := c.Request().ParseForm(); err != nil {
			http.Error(c.Response().Writer, "Parse Form Failed", http.StatusInternalServerError)
			return nil
		}

		state := c.Request().Form.Get("state")
		if state != expectedState {
			response(false, c.Response().Writer, http.StatusBadRequest, "invalid state")
			return nil
		}
		code := c.Request().Form.Get("code")
		if code == "" {
			response(false, c.Response().Writer, http.StatusBadRequest, "invalid code")
			return nil
		}

		token, err := config.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
		if err != nil {
			response(false, c.Response().Writer, http.StatusInternalServerError, "exchange token failed")
			return nil
		}

		globalToken = token

		prettyPrintToken()

		response(true, c.Response().Writer, http.StatusOK)
		return nil
	})

	e.GET("/refresh", func(c echo.Context) error {

		w := c.Response().Writer
		r := c.Request()

		if globalToken == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return nil
		}

		globalToken.Expiry = time.Now()
		token, err := config.TokenSource(context.Background(), globalToken).Token()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}

		globalToken = token

		prettyPrintToken()

		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)

		return nil
	})

	e.GET("/account", func(c echo.Context) error {
		w := c.Response().Writer
		r := c.Request()

		if globalToken == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return nil
		}
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/oauth/exchange/account", sourceSeverURL), nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}

		req.Header.Set("Authorization", "Bearer "+globalToken.AccessToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil
		}
		defer resp.Body.Close()

		io.Copy(w, resp.Body)
		return nil
	})

	log.Printf("Client is running at %s port.", port)
	log.Fatal(e.Start(":" + port))
}

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}

func prettyPrintToken() {

	// Some Strings
	string_1 := []string{globalToken.AccessToken, globalToken.RefreshToken, globalToken.Expiry.String(), globalToken.TokenType}

	// Create Object
	tabulate := gotabulate.Create([][]string{string_1})

	// Set Headers
	tabulate.SetHeaders([]string{"Access Token", "Refresh Token", "Expiry", "Token Type"})

	// Render
	fmt.Println(tabulate.Render("simple"))
}

type tokenResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func response(success bool, w http.ResponseWriter, code int, message ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	res := tokenResponse{
		Success: success,
	}

	if len(message) > 0 {
		res.Message = message[0]
	}

	data, _ := json.Marshal(res)
	w.Write(data)
}
