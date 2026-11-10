package util

import (
	"net/http"

	"github.com/ockendenjo/ale-trail/pkg/apighandler"
	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/ockendenjo/handler"
)

func GetStashForDate(logger *handler.Logger, fileWithEtag *data.FileWithETag, date, key string) (*data.Stash, error) {
	var foundStash *data.Stash
	for _, stash := range fileWithEtag.File.Stashes {
		if stash.Date == date {
			foundStash = stash
			break
		}
	}

	if foundStash == nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusNotFound, Message: "no stash found", CWLogMsg: "stash not found"}
	}
	logger.Info("Found stash")

	switch key {
	case foundStash.StashKey:
		logger.Info("Authorized using stash key")
	case fileWithEtag.File.AdminKey:
		logger.Info("Authorized using admin key")
	default:
		return nil, apighandler.HttpError{StatusCode: http.StatusUnauthorized, Message: "incorrect key", CWLogMsg: "incorrect stash key"}
	}

	return foundStash, nil
}

