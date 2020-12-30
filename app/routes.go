package app

func (s *Server) routes() {
	s.Router.HandleFunc("/api/v1/sets", s.handleGetSets()).Methods("GET")
	s.Router.HandleFunc("/api/v1/sets", s.handleCreateSet()).Methods("POST")
	s.Router.HandleFunc("/api/v1/sets/{id:[0-9]+}", s.handleGetSet()).Methods("GET")
	s.Router.HandleFunc("/api/v1/sets/{id:[0-9]+}", s.handleUpdateSet()).Methods("PUT")
	s.Router.HandleFunc("/api/v1/sets/{id:[0-9]+}", s.handleDeleteSet()).Methods("DELETE")
	s.Router.HandleFunc("/api/login", Login).Methods("POST")
	s.Router.HandleFunc("/api/welcome", Welcome).Methods("GET")
	s.Router.HandleFunc("/api/refresh", Refresh).Methods("POST")
}