package auth

import (
	"net/http"
)

func (s *Server) RevokeEndpoint(rw http.ResponseWriter, req *http.Request) {
	// This context will be passed to all methods.
	ctx := req.Context()

	// This will accept the token revocation request and validate various parameters.
	err := s.provider.NewRevocationRequest(ctx, req)

	// All done, send the response.
	s.provider.WriteRevocationResponse(ctx, rw, err)
}
