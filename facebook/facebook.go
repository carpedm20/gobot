package facebook

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	u "net/url"
	"strings"
)

const (
	page_access_token_url = "https://graph.facebook.com/me/accounts?access_token="
	api_url               = "https://graph.facebook.com/v2.2"
)

type Facebook struct {
	Client   http.Client
	Accounts []Account
}

type Account struct {
	Name         string
	Client       http.Client
	Category     string
	Id           string
	Access_token string
	Perms        []string
}

type Post struct {
	Id string
}

type Data struct {
	Data []Account
}

type Error struct {
	Error struct {
		Message          string
		Type             string
		Code             int
		Error_subcode    int
		Is_transient     bool
		Error_user_title string
		Error_user_msg   string
	}
}

func New() *Facebook {
	client := http.Client{}
	accounts := []Account{}
	return &Facebook{
		Client:   client,
		Accounts: accounts,
	}
}

func (g *Facebook) SetProxy(url string) {
	proxyURL, err := u.Parse(url)
	if err != nil {
		panic(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	if proxyURL.Scheme == "https" {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	g.Client.Transport = transport
}

func (g *Facebook) Login(accessToken string) {
	resp, err := g.Client.Get(page_access_token_url + accessToken)
	if err != nil {
		panic(err)
	}
	data := &Data{}
	GetStructFromResponse(resp, data)
	for _, account := range data.Data {
		account.Client = g.Client
		g.Accounts = append(g.Accounts, account)
	}
}

func (g *Facebook) GetAccessByName(name string) (target Account) {
	for _, account := range g.Accounts {
		if account.Name == name {
			target = account
		}
	}
	return
}

func (a *Account) Post(text string) {
	data := strings.NewReader(
		u.Values{
			"message":      {text},
			"access_token": {a.Access_token},
		}.Encode(),
	)
	resp, err := a.Client.Post(api_url+"/feed", "application/json", data)
	if err != nil {
		panic(err)
	}
	post := &Post{}
	GetStructFromResponse(resp, post)
}

func GetStructFromResponse(resp *http.Response, data interface{}) {
	bodyIo, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	body := (string(bodyIo))
	if strings.Index(body, "\"error\":") != -1 {
		error := &Error{}
		err := json.Unmarshal([]byte(body), &error)
		if err != nil {
			panic(err)
		}
		log.Printf("[!] Err: %+v\n", error)
	} else {
		err := json.Unmarshal([]byte(body), &data)
		if err != nil {
			panic(err)
		}
		log.Printf("[*] Inst: %+v\n", data)
	}
}
