# Privacy Policy (template for self-hosted deployments)

> **Not legal advice.** Replace bracketed placeholders before publishing on your WiChat instance.

**Effective date:** [DATE]

**Data controller:** [YOUR ORGANIZATION NAME]  
**Contact:** [PRIVACY EMAIL]

## Overview

[YOUR ORGANIZATION NAME] operates a self-hosted WiChat instance at [INSTANCE URL]. WiChat is open-source communication software; **your organization** is responsible for personal data processed on your deployment.

## Data we process

Depending on how you configure WiChat, we may process:

- Account data (email, display name, profile image)
- Authentication data (hashed passwords, OAuth identifiers, session tokens)
- Messages, files, and metadata (channels, reactions, read state)
- Technical logs (IP address, user agent, timestamps) for security and abuse prevention

## Purposes and legal bases

| Purpose | Examples | Typical basis (adjust for your jurisdiction) |
|---------|----------|-----------------------------------------------|
| Provide the service | Chat, calls, notifications | Contract / legitimate interest |
| Security & abuse prevention | Rate limits, audit logs | Legitimate interest |
| Compliance | Legal requests you respond to | Legal obligation |

## Retention

Configure retention in workspace settings (see WiChat admin documentation). Defaults are described in the system design (data retention vs trash grace period).

## Subprocessors

Self-hosted WiChat may integrate with services **you** configure, for example:

- Google OAuth (if enabled)
- SMTP provider (password reset email)
- Firebase Cloud Messaging (Android push, if used)

List your choices here.

## Your rights

Depending on applicable law (e.g. GDPR), users may request access, correction, deletion, or restriction. **Account deletion** in WiChat anonymizes the user while preserving message context; describe your process for erasure requests.

## International transfers

[DESCRIBE IF YOU HOST OUTSIDE USER REGION]

## Changes

We will post updates at [POLICY URL] with a revised effective date.

---

Generated from WiChat `PRIVACY_TEMPLATE.md`. Customize using a generator such as Termly or TermsFeed if needed.
