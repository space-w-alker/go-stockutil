package pathutil

import (
	"fmt"
	"os/user"
	"strings"
)

func ExpandUser(path string) (string, error) {
	if u, err := user.Current(); err == nil {
		fullTilde := fmt.Sprintf("~%s", u.Name)

		if strings.HasPrefix(path, `~/`) || path == `~` {
			return strings.Replace(path, `~`, u.HomeDir, 1), nil
		}

		if strings.HasPrefix(path, fullTilde+`/`) || path == fullTilde {
			return strings.Replace(path, fullTilde, u.HomeDir, 1), nil
		}

		return path, nil
	} else {
		return path, err
	}
}
