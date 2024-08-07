package auth

import (
	"net/http"
)

func (s *Server) IntrospectionEndpoint(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	introspectionResponder, err := s.provider.NewIntrospectionRequest(ctx, req, s.newSession("", nil))
	if err != nil {
		s.provider.WriteIntrospectionError(ctx, rw, err)
		return
	}
	s.provider.WriteIntrospectionResponse(ctx, rw, introspectionResponder)
}
