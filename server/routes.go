package server

func (s *Server) routes() {
	s.Router.HandleFunc("/", s.healthCheckHandle()).Methods("GET")
	s.Router.HandleFunc("/slack/hook", nil).Methods("POST")
	s.Router.HandleFunc("/telegram/hook", nil).Methods("POST")
}   
