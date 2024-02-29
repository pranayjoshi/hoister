package upload

import (
	"os"

	"github.com/go-git/go-git/v5"
)

func CloneRepo() {

	_, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
		URL:      "https://github.com/go-git/go-git",
		Progress: os.Stdout,
	})
	if err != nil {
		panic(err)
	}
}
