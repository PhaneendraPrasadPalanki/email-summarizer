package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type EmailSummary struct {
	Subject string
	Summary string
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open("token.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to save token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Request a token from the web, then return the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser:\n\n%v\n\n", authURL)

	fmt.Print("Enter the authorization code: ")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read auth code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func getClient(config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile("token.json")
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken("token.json", tok)
	}
	return config.Client(context.Background(), tok)
}

func fetchEmails() []EmailSummary {
	ctx := context.Background()

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read credentials.json: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file: %v", err)
	}

	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	r, err := srv.Users.Messages.List(user).MaxResults(10).Q("category:primary").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	var emails []EmailSummary

	for _, m := range r.Messages {
		msg, _ := srv.Users.Messages.Get(user, m.Id).Format("full").Do()

		var body string
		for _, part := range msg.Payload.Parts {
			if part.MimeType == "text/plain" {
				data, _ := base64.URLEncoding.DecodeString(part.Body.Data)
				body = string(data)
				break
			}
		}
		if body == "" {
			continue
		}

		subject := "Unknown"
		for _, header := range msg.Payload.Headers {
			if header.Name == "Subject" {
				subject = header.Value
				break
			}
		}

		emails = append(emails, EmailSummary{
			Subject: subject,
			Summary: "📄 (Mock Summary) This email's summary will appear here once OpenAI is connected.",
		})
	}

	return emails
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	emails := fetchEmails()

	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, "Error loading template", 500)
		return
	}

	tmpl.Execute(w, emails)
}

func main() {
	http.HandleFunc("/", dashboardHandler)
	fmt.Println("📊 Dashboard running at: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
