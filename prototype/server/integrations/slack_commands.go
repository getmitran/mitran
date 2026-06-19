package integrations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"github.com/getmitran/mitran/server/apierr"
)

type SlackCommandResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

type SlackCommandHandler struct {
	// Dependencies injected at creation
	GetStatus  func() string
	ApproveTask func(id string) error
	RejectTask  func(id string) error
	ChatAgent   func(agent, msg string) (string, error)
}

func NewSlackCommandHandler() *SlackCommandHandler {
	return &SlackCommandHandler{
		GetStatus:   defaultGetStatus,
		ApproveTask: defaultApproveTask,
		RejectTask:  defaultRejectTask,
		ChatAgent:   defaultChatAgent,
	}
}

func (h *SlackCommandHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/slack/commands", h.HandleSlashCommand)
}

func (h *SlackCommandHandler) HandleSlashCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}

	if err := r.ParseForm(); err != nil {
		respondSlack(w, "ephemeral", "Failed to parse request")
		return
	}

	command := r.FormValue("command")
	text := strings.TrimSpace(r.FormValue("text"))

	if command != "/mitran" {
		respondSlack(w, "ephemeral", fmt.Sprintf("Unknown command: %s", command))
		return
	}

	parts := strings.Fields(text)
	if len(parts) == 0 {
		respondSlack(w, "ephemeral", "Usage: /mitran <status|approve|reject|chat> [args]")
		return
	}

	subcommand := parts[0]

	switch subcommand {
	case "status":
		result := h.GetStatus()
		respondSlack(w, "in_channel", result)

	case "approve":
		if len(parts) < 2 {
			respondSlack(w, "ephemeral", "Usage: /mitran approve <task-id>")
			return
		}
		if err := h.ApproveTask(parts[1]); err != nil {
			respondSlack(w, "ephemeral", fmt.Sprintf("❌ Failed to approve: %s", err.Error()))
			return
		}
		respondSlack(w, "in_channel", fmt.Sprintf("✅ Task `%s` approved", parts[1]))

	case "reject":
		if len(parts) < 2 {
			respondSlack(w, "ephemeral", "Usage: /mitran reject <task-id>")
			return
		}
		if err := h.RejectTask(parts[1]); err != nil {
			respondSlack(w, "ephemeral", fmt.Sprintf("❌ Failed to reject: %s", err.Error()))
			return
		}
		respondSlack(w, "in_channel", fmt.Sprintf("🚫 Task `%s` rejected", parts[1]))

	case "chat":
		if len(parts) < 3 {
			respondSlack(w, "ephemeral", "Usage: /mitran chat <agent> <message>")
			return
		}
		agent := parts[1]
		msg := strings.Join(parts[2:], " ")
		reply, err := h.ChatAgent(agent, msg)
		if err != nil {
			respondSlack(w, "ephemeral", fmt.Sprintf("❌ Agent error: %s", err.Error()))
			return
		}
		respondSlack(w, "in_channel", fmt.Sprintf("🤖 *%s*: %s", agent, reply))

	default:
		respondSlack(w, "ephemeral", fmt.Sprintf("Unknown subcommand: %s\nAvailable: status, approve, reject, chat", subcommand))
	}
}

func respondSlack(w http.ResponseWriter, responseType, text string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SlackCommandResponse{
		ResponseType: responseType,
		Text:         text,
	})
}

// Default implementations — replaced via dependency injection in production
func defaultGetStatus() string {
	return "🟢 Mitran engine running\n• Agents: 8 registered\n• Tasks: 0 pending approval"
}

func defaultApproveTask(id string) error {
	return fmt.Errorf("task %s not found", id)
}

func defaultRejectTask(id string) error {
	return fmt.Errorf("task %s not found", id)
}

func defaultChatAgent(agent, msg string) (string, error) {
	return fmt.Sprintf("Echo from %s: %s", agent, msg), nil
}
