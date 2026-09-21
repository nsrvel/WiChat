# Security Policy

## Supported versions

| Version | Supported |
|---------|-----------|
| main (pre-release) | Best effort |

Tagged releases will list supported versions here once shipping starts.

## Reporting a vulnerability

**Do not open public GitHub issues for exploitable security bugs.**

Email the maintainers with:

- Description and impact
- Steps to reproduce
- Affected component (core, client, plugin, deployment)
- Proof-of-concept if available

We aim to acknowledge within **72 hours** and provide a remediation timeline when confirmed.

## Scope

In scope: WiChat core, official clients, default deployment (`docker compose`), and documented plugin contracts.

Out of scope for the core team: misconfiguration of self-hosted instances (weak passwords, exposed admin ports, missing TLS) unless caused by default insecure settings in our compose templates.

## Safe harbor

Good-faith research that avoids privacy violations, data destruction, and service disruption is appreciated.
