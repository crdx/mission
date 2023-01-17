package util

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

func GetAbsoluteDir(relPath string, userName string) (string, error) {
	if relPath == "" {
		return "", fmt.Errorf("GetAbsoluteDir must not be called with an empty string")
	}

	if strings.HasPrefix(relPath, "~") {
		user, err := user.Lookup(userName)
		if err != nil {
			return "", err
		}

		if relPath == "~" {
			return user.HomeDir, nil
		} else {
			return filepath.Join(user.HomeDir, relPath[2:]), nil
		}
	} else {
		return filepath.Abs(relPath)
	}
}

func IsGitRepository(str string) bool {
	return IsDirectory(path.Join(str, ".git"))
}

func IsDirectory(str string) bool {
	stat, err := os.Stat(str)
	return err == nil && stat.IsDir()
}

func IsExecutable(str string) bool {
	stat, err := os.Stat(str)
	return err == nil && stat.Mode()&0111 != 0
}

func IsReadableFile(str string) bool {
	file, err := os.Open(str)
	file.Close()
	return err == nil
}

func ChownDirectory(dir string, userId int, groupId int) (int, error) {
	count := 0
	return count, filepath.WalkDir(dir, func(path string, _ fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory.
		if path == dir {
			return nil
		}

		if err := os.Chown(path, userId, groupId); err != nil {
			return err
		}

		count++
		return nil
	})
}
