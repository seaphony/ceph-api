package auth

import (
	"context"
	"net/http"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) AuthEndpoint(rw http.ResponseWriter, req *http.Request) {
	// This context will be passed to all methods.
	ctx := req.Context()
	log := zerolog.Ctx(req.Context())

	ar, err := s.provider.NewAuthorizeRequest(ctx, req)
	if err != nil {
		fositeError := fosite.ErrorToRFC6749Error(err)
		log.Error().Str("error_description", fositeError.GetDescription()).Err(err).Msg("Error occurred in NewAuthorizeRequest")
		s.provider.WriteAuthorizeError(ctx, rw, ar, err)
		return
	}

	prompt := ar.GetRequestForm().Get("prompt")
	if prompt == "none" {
		err := fosite.ErrLoginRequired.WithHint("Failed validate open id request cause prompt set to none but there is no valid login session for this user")
		s.provider.WriteAuthorizeError(ctx, rw, ar, err)
		return
	}

	username := ar.GetRequestForm().Get("username")
	password := ar.GetRequestForm().Get("password")

	log.Info().Str("username", username).Err(err).Msg("login via password")
	s.authorize(ctx, ar, rw, username, password)
}

func (s *Server) authorize(ctx context.Context, ar fosite.AuthorizeRequester, rw http.ResponseWriter, username, password string) {
	log := zerolog.Ctx(ctx)
	session := s.newSession(username, nil)

	// Set subject to the session
	usr, err := s.userSvc.GetUser(ctx, username)
	if err != nil {
		log.Error().Err(err).Msg("could not find account for subject")
		http.Error(rw, "can't find account", http.StatusForbidden)
		return
	}
	if !usr.Enabled {
		http.Error(rw, "access denied", http.StatusForbidden)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		log.Error().Err(err).Msg("could not validate password")
		http.Error(rw, "can't find account", http.StatusForbidden)
		return
	}

	// grant requested scopes
	for _, scope := range ar.GetRequestedScopes() {
		ar.GrantScope(scope)
	}
	jwtSession := session.(*oauth2.JWTSession)
	jwtSession.Subject = usr.Username

	response, err := s.provider.NewAuthorizeResponse(ctx, ar, session)

	if err != nil {
		fositeError := fosite.ErrorToRFC6749Error(err)
		log.Error().Str("error_description", fositeError.GetDescription()).Err(err).Msg("Error occurred in NewAuthorizeResponse")
		s.provider.WriteAuthorizeError(ctx, rw, ar, err)
		return
	}
	s.provider.WriteAuthorizeResponse(ctx, rw, ar, response)
}
