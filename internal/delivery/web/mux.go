package web

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/slim-crown/issue-1-website/internal/services/session"
	issue1 "github.com/slim-crown/issue-1-website/pkg/issue1.REST.client/http.issue1"

	"github.com/julienschmidt/httprouter"
)

// Setup is used to inject dependencies and other required data used by the handlers.
type Setup struct {
	Config
	Dependencies
	templates *template.Template
}

// Dependencies contains dependencies used by the handlers.
type Dependencies struct {
	Iss1C          *issue1.Client
	Logger         *log.Logger
	sessionValues  sessionValues
	SessionService session.Service
}

// Config contains the different settings used to set up the handlers
type Config struct {
	TemplatesStoragePath, AssetStoragePath, AssetServingRoute string
	HostAddress, Port                                         string
	CookieName                                                string
	CSRFTokenLifetime                                         time.Duration
	SessionIdleLifetime, SessionHardLifetime                  time.Duration
	TokenSigningSecret                                        []byte
	HTTPS                                                     bool
}

// NewMux returns a fully configured issue1 website server.
func NewMux(s *Setup) *httprouter.Router {
	mainRouter := httprouter.New()

	err := s.ParseTemplates()
	if err != nil {
		s.Logger.Printf("error: initial template parsing failed because: %v", err)
		//s.Logger.Fatalf("fatal: server start-up aborted.")
	}
	s.sessionValues.restRefreshToken = "restRefreshToken"
	s.sessionValues.csrf = "CSRF"
	s.sessionValues.username = "username"

	fs := http.FileServer(http.Dir(s.AssetStoragePath))
	mainRouter.Handler("GET", s.AssetServingRoute+"*filepath", http.StripPrefix(s.AssetServingRoute, fs))

	mainRouter.HandlerFunc("GET", "/", getFront(s))
	mainRouter.HandlerFunc("POST", "/login", postLogin(s))
	mainRouter.HandlerFunc("POST", "/signup", postSignUp(s))
	mainRouter.HandlerFunc("GET", "/home", getHome(s))
	mainRouter.HandlerFunc("POST", "/home-feed-posts", postFeedPosts(s))
	mainRouter.HandlerFunc("GET", "/error", getError(s))
	mainRouter.HandlerFunc("GET", "/404", get404(s))
	mainRouter.HandlerFunc("GET", "/p/:postID", getPostView(s))
	mainRouter.HandlerFunc("POST", "/p/:postID/comment-board", postPostComments(s))
	mainRouter.HandlerFunc("POST", "/p/:postID/add-comment", postComment(s))

	return mainRouter
}
