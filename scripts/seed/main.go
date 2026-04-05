package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// ─── Seed data — edit these to change what gets created ───────────────────────

type seedSuperuser struct{ Email, Password string }
type seedUser struct{ Email, Password, Name string }
type seedPost struct{ Title, Body, AuthorEmail string }

var superusers = []seedSuperuser{
	{Email: "root@dev.com", Password: "root1234"},
}

var users = []seedUser{
	{Email: "alice@dev.com", Password: "password1234", Name: "Alice"},
	{Email: "bob@dev.com", Password: "password1234", Name: "Bob"},
}

var posts = []seedPost{
	{Title: "Hello World", Body: "First test post.", AuthorEmail: "alice@dev.com"},
}

// ─── Main ─────────────────────────────────────────────────────────────────────

func main() {
	port := os.Getenv("PUBLIC_PB_PORT")
	if port == "" {
		port = "8090"
	}
	base := "http://localhost:" + port
	if override := os.Getenv("PB_URL"); override != "" {
		base = override
	}

	s := &seeder{base: base, client: &http.Client{}}

	if err := s.checkHealth(); err != nil {
		fatalf("PocketBase is not reachable (%v).\n       Start the server with `task dev` first.", err)
	}

	if err := s.bootstrap(superusers[0]); err != nil {
		fatalf("%v", err)
	}
	for _, su := range superusers[1:] {
		s.ensureSuperuser(su)
	}
	for _, u := range users {
		s.ensureUser(u)
	}
	for _, p := range posts {
		s.ensurePost(p)
	}

	fmt.Println("done.")
}

// ─── Seeder ──────────────────────────────────────────────────────────────────

type seeder struct {
	base   string
	token  string
	client *http.Client
}

// bootstrap authenticates as the first seed superuser.
// On fresh pb_data with no superusers yet, it creates one via PocketBase's
// first-run API (the same path the setup wizard uses), then authenticates.
func (s *seeder) bootstrap(su seedSuperuser) error {
	token, err := s.authSuperuser(su.Email, su.Password)
	if err == nil {
		s.token = token
		fmt.Printf("superuser %s: authenticated\n", su.Email)
		return nil
	}

	// Auth failed — try first-run creation (only works when _superusers is empty).
	resp, postErr := s.post("/api/collections/_superusers/records", map[string]string{
		"email":           su.Email,
		"password":        su.Password,
		"passwordConfirm": su.Password,
	}, "")
	if postErr != nil || resp.StatusCode != http.StatusOK {
		var detail string
		if resp != nil {
			b, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if len(b) > 0 {
				detail = ": " + string(b)
			}
		}
		return fmt.Errorf(
			"authentication failed and superuser creation failed%s\n"+
				"       Hint: run `go run ./cmd/server superuser create %s %s`",
			detail, su.Email, su.Password,
		)
	}
	resp.Body.Close()
	fmt.Printf("superuser %s: created\n", su.Email)

	token, err = s.authSuperuser(su.Email, su.Password)
	if err != nil {
		return fmt.Errorf("superuser created but auth still failed: %w", err)
	}
	s.token = token
	return nil
}

func (s *seeder) authSuperuser(email, password string) (string, error) {
	resp, err := s.post("/api/collections/_superusers/auth-with-password",
		map[string]string{"identity": email, "password": password}, "")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}
	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Token, nil
}

func (s *seeder) ensureSuperuser(su seedSuperuser) {
	if s.findRecord("_superusers", "email", su.Email) != "" {
		fmt.Printf("superuser %s: skipping (exists)\n", su.Email)
		return
	}
	resp, err := s.post("/api/collections/_superusers/records", map[string]string{
		"email":           su.Email,
		"password":        su.Password,
		"passwordConfirm": su.Password,
	}, s.token)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			resp.Body.Close()
		}
		fmt.Fprintf(os.Stderr, "superuser %s: error creating\n", su.Email)
		return
	}
	resp.Body.Close()
	fmt.Printf("superuser %s: created\n", su.Email)
}

func (s *seeder) ensureUser(u seedUser) {
	if s.findRecord("users", "email", u.Email) != "" {
		fmt.Printf("user %s: skipping (exists)\n", u.Email)
		return
	}
	if err := s.createRecord("users", map[string]string{
		"email":           u.Email,
		"password":        u.Password,
		"passwordConfirm": u.Password,
		"name":            u.Name,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "user %s: %v\n", u.Email, err)
		return
	}
	fmt.Printf("user %s: created\n", u.Email)
}

func (s *seeder) ensurePost(p seedPost) {
	if s.findRecord("posts", "title", p.Title) != "" {
		fmt.Printf("post %q: skipping (exists)\n", p.Title)
		return
	}
	authorID := s.findRecord("users", "email", p.AuthorEmail)
	if authorID == "" {
		fmt.Fprintf(os.Stderr, "post %q: author %s not found\n", p.Title, p.AuthorEmail)
		return
	}
	if err := s.createRecord("posts", map[string]string{
		"title":  p.Title,
		"body":   p.Body,
		"author": authorID,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "post %q: %v\n", p.Title, err)
		return
	}
	fmt.Printf("post %q: created\n", p.Title)
}

// ─── HTTP helpers ─────────────────────────────────────────────────────────────

func (s *seeder) checkHealth() error {
	resp, err := s.client.Get(s.base + "/api/health")
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return nil
}

// findRecord returns the ID of the first record matching field=value, or "".
func (s *seeder) findRecord(collection, field, value string) string {
	filter := fmt.Sprintf("(%s='%s')", field, value)
	u := fmt.Sprintf("%s/api/collections/%s/records?filter=%s",
		s.base, collection, url.QueryEscape(filter))
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	if s.token != "" {
		req.Header.Set("Authorization", "Bearer "+s.token)
	}
	resp, err := s.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			resp.Body.Close()
		}
		return ""
	}
	defer resp.Body.Close()
	var result struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || len(result.Items) == 0 {
		return ""
	}
	return result.Items[0].ID
}

func (s *seeder) createRecord(collection string, body map[string]string) error {
	resp, err := s.post("/api/collections/"+collection+"/records", body, s.token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, b)
	}
	return nil
}

func (s *seeder) post(path string, body map[string]string, token string) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, s.base+path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return s.client.Do(req)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
