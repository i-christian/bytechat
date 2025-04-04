package server

import (
	"net/http"

	"bytechat/cmd/web/chat"
)

// showChatPage renders a chat page
func (s *Server) showChatPage(w http.ResponseWriter, r *http.Request) {
	s.renderComponent(w, r, chat.ChatPage())
}
