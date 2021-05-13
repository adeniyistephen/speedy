package app

import (
	"context"
	"net/http"
	"github.com/pkg/errors"
)


type speedyGroup struct {
	speedy Speedy
}

func (sg speedyGroup) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	speedtest, err := sg.speedy.Fasttest()
	if err != nil {
		return errors.Wrap(err, "unable to query for users")
	}

	return Respond(w, speedtest, http.StatusOK)
}
