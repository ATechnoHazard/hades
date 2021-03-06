package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/api/views"
	u "github.com/ATechnoHazard/hades-2/internal/utils"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
	"github.com/ATechnoHazard/hades-2/pkg/event"
	"github.com/ATechnoHazard/janus"
	"github.com/julienschmidt/httprouter"
)

func saveEvent(eSvc event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		jtk := ctx.Value("janus_context").(*janus.Account)

		if jtk.Role != "admin" {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

		e := &entities.Event{}
		if err := json.NewDecoder(r.Body).Decode(e); err != nil {
			views.Wrap(err, w)
			return
		}

		if e.OrganizationID == 0 {
			e.OrganizationID = tk.OrgID
		}

		if tk.OrgID != e.OrganizationID {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

		e, err := eSvc.SaveEvent(e)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		msg := u.Message(http.StatusOK, "Event successfully saved")
		msg["event"] = e

		u.Respond(w, msg)
		return
	}
}

func getEvent(eSvc event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		e := &entities.Event{}
		if err := json.NewDecoder(r.Body).Decode(e); err != nil {
			views.Wrap(err, w)
			return
		}

		evnt, err := eSvc.ReadEvent(e.ID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		if evnt.OrganizationID != tk.OrgID {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from accessing this resource"))
			return
		}

		msg := u.Message(http.StatusOK, "Event successfully retrieved")
		msg["event"] = evnt

		u.Respond(w, msg)
		return
	}
}

func deleteEvent(eSvc event.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		jtk := ctx.Value("janus_context").(*janus.Account)

		if jtk.Role != "admin" {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

		e := &entities.Event{}
		if err := json.NewDecoder(r.Body).Decode(e); err != nil {
			views.Wrap(err, w)
			return
		}

		evnt, err := eSvc.ReadEvent(e.ID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		if evnt.OrganizationID != tk.OrgID {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

		if err := eSvc.DeleteEvent(evnt.ID); err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Event successfully deleted"))
		return
	}
}

func MakeEventHandler(r *httprouter.Router, eSvc event.Service, j *janus.Janus) {
	r.HandlerFunc("POST", "/api/v2/event/save", middleware.JwtAuthentication(j.GetHandler(saveEvent(eSvc))))
	r.HandlerFunc("POST", "/api/v2/event/read", middleware.JwtAuthentication(getEvent(eSvc)))
	r.HandlerFunc("DELETE", "/api/v2/event/delete", middleware.JwtAuthentication(j.GetHandler(deleteEvent(eSvc))))
}
