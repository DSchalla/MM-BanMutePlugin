package banmuteplugin

import (
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/mattermost/mattermost-server/model"
)

type Config struct {
	Matching               string
	Mode string
}

type Server struct {
	config Config
	api  plugin.API
}

func NewServer(api plugin.API, config Config) (*Server, error) {
	s := Server{}
	s.api = api
	s.config = config
	return &s, nil
}


func (s *Server) HandleMessage(post *model.Post, intercept bool) (*model.Post, string) {
	return post, ""
}

func (s *Server) ReloadConfig(config Config) {
	s.config = config
}
