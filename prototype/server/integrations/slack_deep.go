package integrations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/getmitran/mitran/server/apierr"
)

// SlackDeepHandler handles advanced Slack integrations:
// /mitran slash commands (ask, status, approve) and interactive button payloads.
type SlackDeepHandler struct {
	GetStatus   func() string
	ApproveTask func(id string) error
	AskAgent    func(agent, question string) (string, error)
}

func NewSlackDeepHandler() *SlackDeepHandler {
	return &SlackDeepHandler{
		GetStatus:   defaultDeepGetStatus,
		ApproveTask: defaultDeepApproveTask,
		AskAgent:    defaultDeepAskAgent,
	}
}

func (h *SlackDeepHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/integrations/slack/command", h.HandleCommand)
	mux.HandleFunc("/api/v1/integrations/slack/interact", h.HandleInteraction)
}

// SlackBlock represents a Slack Block Kit element.
type SlackBlock struct {
	Type string      `json:"type"`
	Text interface{} `json:"text,omitempty"`
}

// SlackBlockResponse is a rich Slack response with blocks.
type SlackBlockResponse struct {
	ResponseType    string      `json:"response_type"`
	Text            string      `json:"text,omitempty"`
	Blocks          interface{} `json:"blocks,omitempty"`
	ReplaceOriginal bool        `json:"replace_original,omitempty"`
}

// HandleCommand handles /mitran slash commands: ask, status, approve.
func (h *SlackDeepHandler) HandleCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}

	if err := r.ParseForm(); err != nil {
		respondSlackDeep(w, "ephemeral", "Failed to parse request")
		return
	}

	text := strings.TrimSpace(r.FormValue("text"))
	parts := strings.Fields(text)

	if len(parts) == 0 {
		respondSlackDeep(w, "ephemeral", "Usage: /mitran <ask|status|approve> [args]")
		return
	}

	subcommand := strings.ToLower(parts[0])

	switch subcommand {
	case "ask":
		if len(parts) < 2 {
			respondSlackDeep(w, "ephemeral", "Usage: /mitran ask <question>")
			return
		}
		question := strings.Join(parts[1:], " ")
		reply, err := h.AskAgent("default", question)
		if err != nil {
			respondSlackDeep(w, "ephemeral", fmt.Sprintf("❌ Error: %s", err.Error()))
			return
		}
		respondSlackDeep(w, "in_channel", fmt.Sprintf("🤖 %s", reply))

	case "status":
		result := h.GetStatus()
		respondSlackDeep(w, "in_channel", result)

	case "approve":
		if len(parts) < 2 {
			respondSlackDeep(w, "ephemeral", "Usage: /mitran approve <task-id>")
			return
		}
		taskID := parts[1]
		if err := h.ApproveTask(taskID); err != nil {
			respondSlackDeep(w, "ephemeral", fmt.Sprintf("❌ Approval failed: %s", err.Error()))
			return
		}
		respondSlackDeep(w, "in_channel", fmt.Sprintf("✅ Task `%s` approved", taskID))

	default:
		respondSlackDeep(w, "ephemeral", fmt.Sprintf("Unknown subcommand: %s\nAvailable: ask, status, approve", subcommand))
	}
}

// SlackInteractionPayload represents the incoming Slack interaction callback.
type SlackInteractionPayload struct {
	Type    string `json:"type"`
	Actions []struct {
		ActionID string `json:"action_id"`
		Value    string `json:"value"`
		BlockID  string `json:"block_id"`
	} `json:"actions"`
	User struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"user"`
	ResponseURL string `json:"response_url"`
}

// HandleInteraction handles Slack interactive component payloads (button clicks).
func (h *SlackDeepHandler) HandleInteraction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}

	if err := r.ParseForm(); err != nil {
		apierr.WriteError(w, apierr.BadRequest("failed to parse form"))
		return
	}

	payloadStr := r.FormValue("payload")
	if payloadStr == "" {
		apierr.WriteError(w, apierr.BadRequest("missing payload"))
		return
	}

	var payload SlackInteractionPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		apierr.WriteError(w, apierr.BadRequest("invalid interaction payload"))
		return
	}

	if len(payload.Actions) == 0 {
		respondSlackDeep(w, "ephemeral", "No actions received")
		return
	}

	action := payload.Actions[0]

	switch {
	case strings.HasPrefix(action.ActionID, "approve_"):
		taskID := strings.TrimPrefix(action.ActionID, "approve_")
		if err := h.ApproveTask(taskID); err != nil {
			respondSlackBlock(w, SlackBlockResponse{
				ResponseType:    "ephemeral",
				Text:            fmt.Sprintf("❌ Failed to approve task %s: %s", taskID, err.Error()),
				ReplaceOriginal: false,
			})
			return
		}
		respondSlackBlock(w, SlackBlockResponse{
			ResponseType:    "in_channel",
			Text:            fmt.Sprintf("✅ Task `%s` approved by <@%s>", taskID, payload.User.ID),
			ReplaceOriginal: true,
		})

	case strings.HasPrefix(action.ActionID, "reject_"):
		taskID := strings.TrimPrefix(action.ActionID, "reject_")
		respondSlackBlock(w, SlackBlockResponse{
			ResponseType:    "in_channel",
			Text:            fmt.Sprintf("🚫 Task `%s` rejected by <@%s>", taskID, payload.User.ID),
			ReplaceOriginal: true,
		})

	default:
		respondSlackDeep(w, "ephemeral", fmt.Sprintf("Unknown action: %s", action.ActionID))
	}
}

func respondSlackDeep(w http.ResponseWriter, responseType, text string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SlackCommandResponse{
		ResponseType: responseType,
		Text:         text,
	})
}

func respondSlackBlock(w http.ResponseWriter, resp SlackBlockResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func defaultDeepGetStatus() string {
	return "🟢 Mitran operational\n• Engine: running\n• Agents: ready\n• Tasks: 0 pending"
}

func defaultDeepApproveTask(id string) error {
	return fmt.Errorf("task %s not found in queue", id)
}

func defaultDeepAskAgent(agent, question string) (string, error) {
	return fmt.Sprintf("[%s] I received your question: %s", agent, question), nil
}
