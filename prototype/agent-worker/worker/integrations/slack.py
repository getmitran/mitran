"""Slack integration for Mitran agent worker."""

import os
import requests

SLACK_BOT_TOKEN = os.environ.get("SLACK_BOT_TOKEN", "")
SLACK_API = "https://slack.com/api"


def _headers():
    return {
        "Authorization": f"Bearer {SLACK_BOT_TOKEN}",
        "Content-Type": "application/json",
    }


def send_message(channel: str, text: str, thread_ts: str = None) -> dict:
    """Send a message to a Slack channel."""
    payload = {"channel": channel, "text": text, "mrkdwn": True}
    if thread_ts:
        payload["thread_ts"] = thread_ts
    resp = requests.post(f"{SLACK_API}/chat.postMessage", json=payload, headers=_headers())
    return resp.json()


def send_checkpoint_notification(channel: str, checkpoint: dict) -> dict:
    """Post a formatted checkpoint notification for human approval."""
    agent = checkpoint.get("agent", "unknown")
    task = checkpoint.get("task", "untitled")
    status = checkpoint.get("status", "pending")
    checkpoint_id = checkpoint.get("id", "")

    text = (
        f":large_blue_circle: *Checkpoint Requires Approval*\n"
        f">*Agent:* `{agent}`\n"
        f">*Task:* {task}\n"
        f">*Status:* {status}\n"
        f">*ID:* `{checkpoint_id}`\n\n"
        f"Reply with `approve {checkpoint_id}` or `reject {checkpoint_id}`"
    )
    return send_message(channel, text)


def post_task_update(channel: str, task: dict) -> dict:
    """Post a task status change notification."""
    status_emoji = {
        "pending": ":hourglass_flowing_sand:",
        "running": ":runner:",
        "completed": ":white_check_mark:",
        "failed": ":x:",
        "blocked": ":no_entry_sign:",
    }
    name = task.get("name", "untitled")
    agent = task.get("agent", "unknown")
    status = task.get("status", "pending")
    emoji = status_emoji.get(status, ":grey_question:")

    text = f"{emoji} *{name}*\n>*Agent:* `{agent}` | *Status:* `{status}`"
    return send_message(channel, text)


def get_channel_id(channel_name: str) -> str | None:
    """Resolve a #channel-name to its Slack channel ID."""
    name = channel_name.lstrip("#")
    resp = requests.get(
        f"{SLACK_API}/conversations.list",
        headers=_headers(),
        params={"types": "public_channel,private_channel", "limit": 200},
    )
    data = resp.json()
    if not data.get("ok"):
        return None
    for ch in data.get("channels", []):
        if ch["name"] == name:
            return ch["id"]
    return None
