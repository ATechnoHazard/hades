package handler

import (
	"encoding/json"
	"github.com/ATechnoHazard/hades-2/api/middleware"
	"github.com/ATechnoHazard/hades-2/api/views"
	u "github.com/ATechnoHazard/hades-2/internal/utils"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
	"github.com/ATechnoHazard/hades-2/pkg/organization"
	"github.com/ATechnoHazard/janus"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
)

func acceptJoinRequest(oSvc organization.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		jtk := ctx.Value("janus_context").(*janus.Account)

		if jtk.Role != "admin" {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

		j := &entities.JoinRequest{}
		if err := json.NewDecoder(r.Body).Decode(j); err != nil {
			views.Wrap(err, w)
			return
		}

		if j.OrganizationID != tk.OrgID {
			u.Respond(w, u.Message(http.StatusForbidden, "You are forbidden from modifying this resource"))
			return
		}

		if err := oSvc.AcceptJoinReq(j.OrganizationID, j.Email); err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Join request accepted"))
	}
}

func sendJoinRequest(oSvc organization.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		j := &entities.JoinRequest{}
		if err := json.NewDecoder(r.Body).Decode(j); err != nil {
			views.Wrap(err, w)
			return
		}
		if err := oSvc.SendJoinRequest(j.OrganizationID, tk.Email); err != nil {
			views.Wrap(err, w)
			return
		}

		u.Respond(w, u.Message(http.StatusOK, "Join request created successfully"))
	}
}

func loginOrg(oSvc organization.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		j := &entities.JoinRequest{}
		ctx := r.Context()
		tkn := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		if err := json.NewDecoder(r.Body).Decode(j); err != nil {
			views.Wrap(err, w)
			return
		}

		if tkn.Email != j.Email {
			u.Respond(w, u.Message(http.StatusUnauthorized, "You are not authorized to use this resource"))
			return
		}

		tk, err := oSvc.LoginOrg(j.OrganizationID, j.Email)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		tkString, err := tk.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
		if err != nil {
			views.Wrap(err, w)
			return
		}

		msg := u.Message(http.StatusOK, "Logged in to organization")
		msg["token"] = tkString
		u.Respond(w, msg)
		return
	}
}

func getOrgEvents(oSvc organization.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		events, err := oSvc.GetOrgEvents(tk.OrgID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		msg := u.Message(http.StatusOK, "Successfully retrieved events")
		msg["events"] = events
		u.Respond(w, msg)
		return
	}
}

func viewJoinRequests(oSvc organization.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		reqs, err := oSvc.GetOrgJoinReqs(tk.OrgID)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		msg := u.Message(http.StatusOK, "Successfully retrieved join requests")
		msg["join_reqs"] = reqs
		u.Respond(w, msg)
		return
	}
}

func createOrg(oSvc organization.Service, j *janus.Janus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		org := &entities.Organization{}
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)
		if err := json.NewDecoder(r.Body).Decode(org); err != nil {
			views.Wrap(err, w)
			return
		}

		org, err := oSvc.SaveOrg(org)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		err = oSvc.AddUserToOrg(org.ID, tk.Email)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		tkn, err := oSvc.LoginOrg(org.ID, tk.Email)
		if err != nil {
			views.Wrap(err, w)
			return
		}

		err = j.SetRights(&janus.Account{
			OrganizationID: org.ID,
			Key:            tk.Email,
			Role:           "admin",
		})
		if err != nil {
			views.Wrap(err, w)
			return
		}

		tkString, err := tkn.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
		if err != nil {
			views.Wrap(err, w)
			return
		}

		msg := u.Message(http.StatusOK, "Organization created successfully")
		msg["token"] = tkString
		msg["org"] = org

		u.Respond(w, msg)
		return
	}
}

func MakeOrgHandler(r *httprouter.Router, oSvc organization.Service, j *janus.Janus) {
	r.HandlerFunc("POST", "/api/v2/org/accept", middleware.JwtAuthentication(j.GetHandler(acceptJoinRequest(oSvc))))
	r.HandlerFunc("POST", "/api/v2/org/join", middleware.JwtAuthentication(sendJoinRequest(oSvc)))
	r.HandlerFunc("POST", "/api/v2/org/login-org", middleware.JwtAuthentication(loginOrg(oSvc)))
	r.HandlerFunc("GET", "/api/v2/org/events", middleware.JwtAuthentication(getOrgEvents(oSvc)))
	r.HandlerFunc("GET", "/api/v2/org/requests", middleware.JwtAuthentication(viewJoinRequests(oSvc)))
	r.HandlerFunc("POST", "/api/v2/org/create", middleware.JwtAuthentication(createOrg(oSvc, j)))
}
