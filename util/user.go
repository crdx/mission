package util

import (
	"os/user"
	"strconv"
)

type UserInfo struct {
	UserId  int
	GroupId int
}

func GetUserInfo(userName string) (UserInfo, error) {
	user, err := user.Lookup(userName)
	if err != nil {
		return UserInfo{}, err
	}

	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		return UserInfo{}, err
	}

	gid, err := strconv.Atoi(user.Gid)
	if err != nil {
		return UserInfo{}, err
	}

	return UserInfo{UserId: uid, GroupId: gid}, nil
}
