// server/http.go — just the skeleton, no domains wired yet
func New(cfg *config.Config, db *gorm.DB) *Server {
	return &Server{
		config: cfg,
		db:     db,
	}
}
func (s *Server) Start() error { return s.engine.Start(...) }
func (s *Server) Shutdown(ctx context.Context) error { return s.engine.Shutdown(ctx) }