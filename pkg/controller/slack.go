package controller

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// https://api.slack.com/methods/users.list
// > We recommend no more than 200 results at a time.
const usersOpitonLimit = 190

func (ctrl *Controller) getSlackUser(ctx context.Context, name string, logE *logrus.Entry) (slack.User, error) {
	userPagination := ctrl.SlackBot.GetUsersPaginated(slack.GetUsersOptionLimit(usersOpitonLimit))
	for {
		for _, user := range userPagination.Users {
			// https://api.slack.com/types/user
			user := user
			if getSlackUserName(&user) == name {
				if user.Deleted {
					logE.WithFields(logrus.Fields{
						"user_id":      user.ID,
						"display_name": name,
					}).Warn("user is found, but the user is deleted")
					continue
				}
				return user, nil
			}
		}
		var err error
		userPagination, err = userPagination.Next(ctx)
		if userPagination.Done(err) {
			return slack.User{}, errors.New("user isn't found: " + name)
		}
	}
}

func getSlackUserName(user *slack.User) string {
	for _, name := range []string{user.Profile.DisplayName, user.Profile.RealNameNormalized} {
		if name != "" {
			return name
		}
	}
	return ""
}
