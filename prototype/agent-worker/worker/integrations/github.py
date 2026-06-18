"""GitHub REST API v3 integration for Mitran."""

import os
import base64
from typing import Optional

import requests


class GitHubAPIError(Exception):
    def __init__(self, status_code: int, message: str):
        self.status_code = status_code
        self.message = message
        super().__init__(f"GitHub API error {status_code}: {message}")


def _base_url() -> str:
    return os.environ.get("GITHUB_API_URL", "https://api.github.com").rstrip("/")


def _headers() -> dict:
    token = os.environ.get("GITHUB_TOKEN")
    if not token:
        raise GitHubAPIError(401, "GITHUB_TOKEN environment variable not set")
    return {
        "Authorization": f"token {token}",
        "Accept": "application/vnd.github.v3+json",
    }


def _request(method: str, path: str, **kwargs) -> dict | list:
    url = f"{_base_url()}{path}"
    resp = requests.request(method, url, headers=_headers(), **kwargs)
    if resp.status_code >= 400:
        msg = resp.json().get("message", resp.text) if resp.text else str(resp.status_code)
        raise GitHubAPIError(resp.status_code, msg)
    if resp.status_code == 204:
        return {}
    return resp.json()


def create_repo(org: str, name: str, description: str, private: bool = False) -> dict:
    return _request("POST", f"/orgs/{org}/repos", json={
        "name": name,
        "description": description,
        "private": private,
        "auto_init": True,
    })


def get_repo(owner: str, repo: str) -> dict:
    return _request("GET", f"/repos/{owner}/{repo}")


def list_repos(org: str) -> list:
    return _request("GET", f"/orgs/{org}/repos", params={"per_page": 100})


def create_branch(owner: str, repo: str, branch_name: str, from_branch: str = "main") -> dict:
    ref_data = _request("GET", f"/repos/{owner}/{repo}/git/ref/heads/{from_branch}")
    sha = ref_data["object"]["sha"]
    return _request("POST", f"/repos/{owner}/{repo}/git/refs", json={
        "ref": f"refs/heads/{branch_name}",
        "sha": sha,
    })


def create_pull_request(owner: str, repo: str, title: str, body: str, head: str, base: str = "main") -> dict:
    return _request("POST", f"/repos/{owner}/{repo}/pulls", json={
        "title": title,
        "body": body,
        "head": head,
        "base": base,
    })


def set_branch_protection(owner: str, repo: str, branch: str) -> dict:
    return _request("PUT", f"/repos/{owner}/{repo}/branches/{branch}/protection", json={
        "required_status_checks": {"strict": True, "contexts": []},
        "enforce_admins": True,
        "required_pull_request_reviews": {"required_approving_review_count": 1},
        "restrictions": None,
    })


def create_workflow_file(owner: str, repo: str, workflow_yaml: str, path: str = ".github/workflows/ci.yml") -> dict:
    content_b64 = base64.b64encode(workflow_yaml.encode()).decode()
    return _request("PUT", f"/repos/{owner}/{repo}/contents/{path}", json={
        "message": f"Add workflow: {path}",
        "content": content_b64,
    })


def trigger_workflow(owner: str, repo: str, workflow_id: str) -> dict:
    return _request("POST", f"/repos/{owner}/{repo}/actions/workflows/{workflow_id}/dispatches", json={
        "ref": "main",
    })
