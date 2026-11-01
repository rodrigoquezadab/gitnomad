package github

import "testing"

func TestParse(t *testing.T) {
	data := []byte(`[
	 {"name":"zeta","nameWithOwner":"me/zeta","url":"https://github.com/me/zeta","isPrivate":true,"isFork":false,"isArchived":false,"isEmpty":false,"defaultBranchRef":{"name":"main"}},
	 {"name":"Alpha","nameWithOwner":"me/Alpha","url":"https://github.com/me/Alpha","isPrivate":false,"isFork":true,"isArchived":false,"isEmpty":true,"defaultBranchRef":null}
	]`)
	repos, err := parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(repos) != 2 || repos[0].Name != "Alpha" || repos[1].DefaultBranch != "main" || repos[0].DefaultBranch != "" {
		t.Fatalf("parse: %+v", repos)
	}
	if got := repos[1].CloneURL(); got != "https://github.com/me/zeta.git" {
		t.Errorf("CloneURL: %s", got)
	}
	if _, err := parse([]byte("no json")); err == nil {
		t.Error("se esperaba error")
	}
}
