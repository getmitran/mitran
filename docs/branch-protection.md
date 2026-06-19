# Branch Protection Rules

Recommended GitHub branch protection configuration for `main`.

## Required Settings

| Rule | Value |
|------|-------|
| Require pull request reviews | ✅ (minimum 1 reviewer) |
| Dismiss stale reviews on new commits | ✅ |
| Require status checks to pass | ✅ |
| Required status checks | `ci`, `lint`, `test` |
| Require linear history | ✅ |
| Allow force pushes | ❌ |
| Allow deletions | ❌ |

## Setup via GitHub CLI

```bash
gh api repos/{owner}/{repo}/branches/main/protection -X PUT -f \
  required_pull_request_reviews[dismiss_stale_reviews]=true \
  required_pull_request_reviews[required_approving_review_count]=1 \
  required_status_checks[strict]=true \
  required_status_checks[contexts][]="ci" \
  required_status_checks[contexts][]="lint" \
  required_status_checks[contexts][]="test" \
  required_linear_history[enabled]=true \
  allow_force_pushes[enabled]=false \
  allow_deletions[enabled]=false
```

## Rationale

- **1 reviewer minimum** — catches bugs without bottlenecking small teams.
- **Status checks** — CI, lint, and test must pass before merge.
- **Linear history** — clean `git log`; no merge commits.
- **No force pushes** — protects shared history from accidental rewrites.
- **Dismiss stale reviews** — ensures reviewers re-approve after significant changes.
