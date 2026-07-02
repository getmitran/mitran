# Mitran — Vision & Strategy

## One-Liner
"Describe your team → Mitran sets up everything in 30 minutes."

## The Problem
Every growing team needs internal infrastructure: CI/CD pipelines, documentation, ticketing, monitoring, wikis, HR portals. Today this requires 5-15 engineers spending months building, maintaining, and patching these systems. Most startups can't afford this. Most enterprises waste millions on it.

## The Solution
Mitran is an open-source platform where AI agents build and maintain your entire internal infrastructure. You describe what your team does — agents set up your dev pipeline, docs site, ticketing system, monitoring, wiki, and HR portal. Humans control priority, approve deployments, and steer direction. Agents do all the work.

## Core Principles
1. **Human decides, agents execute** — Sync execution, human checkpoints, priority queue
2. **Opinionated by default, customizable by need** — Ships with all agents, extensible via marketplace
3. **Open-source core, enterprise tier for scale** — AGPL licensed, self-hosted first
4. **Integrate where possible, build where necessary** — GitHub, Slack (integrate). Ticketing, Wiki, Observability (own). Grafana, Prometheus (install & integrate, license-compliant).

## Integration Strategy

| Category | Approach | Tool |
|----------|----------|------|
| Source Control | Integrate | GitHub |
| Communication | Integrate | Slack |
| Ticketing | **Build own** | Mitran Tickets (first-party app) |
| CI/CD | Integrate + orchestrate | GitHub Actions (controlled by Mitran) |
| Documentation | Build own | Mitran Docs (first-party app) |
| Wiki | Build own | Mitran Wiki (first-party app) |
| Observability | Install & integrate | Grafana (Apache 2.0) + Prometheus (Apache 2.0) |
| HR Portal | Build own | Mitran HR (first-party app) |
| Dashboards | Install & integrate | Grafana (Apache 2.0) |
| Alerting | Install & integrate | Alertmanager (Apache 2.0) |

### License Compliance
- Grafana: Apache 2.0 → free to install, bundle, modify
- Prometheus: Apache 2.0 → free to install, bundle, modify
- Alertmanager: Apache 2.0 → free to install, bundle, modify
- We do NOT fork or modify AGPL/SSPL tools (MongoDB, Elasticsearch post-2021)
- We do NOT violate any copyleft by linking — we run them as standalone processes and integrate via APIs

## Business Model

| Tier | Price | Includes |
|------|-------|----------|
| **Community (AGPL)** | Free | Full platform, all agents, self-hosted, local apps |
| **Enterprise** | $5K-50K/year | Managed hosting, SSO/SAML, audit logs, SLA support, custom agent training |
| **Consulting** | $10K-100K/project | "We set up Mitran for your team" — onboarding + customization |

## Team (Phase 1: 5 people)
1. Founder/CEO — Product, strategy, sales, fundraising
2. Senior Go Engineer — Core engine, DAG scheduler, CLI
3. Senior Python Engineer — Agent SDK, all agent implementations
4. Full-stack Engineer — Dashboard (TypeScript/React), Ticketing/Wiki UI
5. DevRel/Designer — Docs, community, demo, branding

Scale to 10 after: 3K GitHub stars + 3 enterprise design partners.

## Timeline

| Phase | Duration | Deliverable |
|-------|----------|-------------|
| Phase 1 | Months 0-6 | Platform + ALL agents (Dev, Docs, CI/CD, Tickets, Wiki, Ops, HR, Dashboards) |
| Phase 2 | Months 6-12 | Marketplace, enterprise tier GA, managed cloud |
| Phase 3 | Months 12-18 | Custom agent training, advanced workflows, Series A |

## Killer Demo (GTM Hook)
User runs: `mitran init`
- "What does your team do?" → "We build a mobile fitness app"
- "How many engineers?" → "8"
- "What languages?" → "TypeScript, React Native, Go backend"

Mitran then:
1. Creates GitHub repos with proper branch protection
2. Sets up CI/CD pipelines (lint, test, deploy)
3. Generates initial documentation structure
4. Creates ticketing board with sprint templates
5. Installs Prometheus + Grafana for monitoring
6. Sets up alerting rules
7. Creates wiki with onboarding guide
8. Sets up HR portal with PTO tracking

All in 30 minutes. Human approves each major step.
