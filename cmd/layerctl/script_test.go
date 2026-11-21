package main

import (
	"os"
	"testing"
	"time"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"layerctl": main,
	})
}

func TestScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/script",
		Setup: func(env *testscript.Env) error {
			env.Setenv("GIT_CONFIG_NOSYSTEM", "1")
			env.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
			env.Setenv("GIT_AUTHOR_NAME", "Layerctl Test")
			env.Setenv("GIT_AUTHOR_EMAIL", "layerctl@example.test")
			env.Setenv("GIT_COMMITTER_NAME", "Layerctl Test")
			env.Setenv("GIT_COMMITTER_EMAIL", "layerctl@example.test")
			setGitDate(env.Setenv, "2000-01-01T00:00:00Z")
			return nil
		},
		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"git": scriptGit,
			"at":  scriptAt,
		},
	})
}

func scriptGit(ts *testscript.TestScript, neg bool, args []string) {
	err := ts.Exec("git", args...)
	if neg {
		if err == nil {
			ts.Fatalf("unexpected git command success")
		}
		return
	}
	ts.Check(err)
}

func scriptAt(ts *testscript.TestScript, neg bool, args []string) {
	if neg || len(args) != 1 {
		ts.Fatalf("usage: at <RFC3339 timestamp>")
	}
	date, err := time.Parse(time.RFC3339, args[0])
	ts.Check(err)
	setGitDate(ts.Setenv, date.Format(time.RFC3339))
}

func setGitDate(setenv func(key, value string), date string) {
	setenv("GIT_AUTHOR_DATE", date)
	setenv("GIT_COMMITTER_DATE", date)
}
