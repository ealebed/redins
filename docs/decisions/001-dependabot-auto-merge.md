# ADR-001: Auto-merge Dependabot minor/patch PRs with a GitHub App

## Status
Accepted

## Date
2026-09-11

## Context
Dependabot opens daily PRs for Go modules, Docker base images, and GitHub Actions. Manual review of minor and patch bumps is slow and low-value. We still want humans to review major (and non-semver) updates, and we want CI to stay a merge gate.

Constraints:

- `GITHUB_TOKEN` reviews are attributed to `github-actions[bot]`. That identity is not a CODEOWNER.
- Official [CODEOWNERS](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners) syntax is users and teams only. GitHub App reviews land as `app-name[bot]` and do not count as code-owner approvals.
- This repository is personal (`ealebed/redins`), so organization ruleset bypass lists are not available.
- Dependabot-triggered `pull_request` workflows receive **Dependabot secrets only**, and `GITHUB_TOKEN` is read-only by default. [Source](https://docs.github.com/en/code-security/reference/supply-chain-security/dependabot-on-actions)

## Decision
Use a dedicated GitHub App (`automerger`) from a GitHub Actions workflow to approve and squash-auto-merge **semver-minor** and **semver-patch** Dependabot PRs.

Keep `.github/CODEOWNERS` so humans still get review requests. Do **not** enable “Require review from Code Owners”. Require **one** approving review plus required status checks. The workflow runs only when the PR author is `dependabot[bot]`.

`gh pr merge --auto --squash` queues the merge. GitHub performs the squash only after required checks pass. Required checks must exist before this workflow is enabled, or a PR can merge with no CI.

Same pattern as `ealebed/token-injector` and `ealebed/dha`.

## Alternatives Considered

### Fine-grained PAT of `@ealebed` (CODEOWNER)
- Pros: Approval would satisfy “Require review from Code Owners”.
- Cons: Long-lived credential tied to a person; revocation or expiry silently stops automation.
- Rejected: The App is the intended identity, and we accepted dropping the code-owner merge gate.

### `GITHUB_TOKEN` / `github-actions[bot]`
- Pros: No extra secrets.
- Cons: Does not satisfy code-owner reviews; still needs “Allow GitHub Actions to create and approve pull requests”; weaker attribution.
- Rejected: We want a dedicated App identity for approve/merge.

### Empty CODEOWNERS entries for lockfiles
- Pros: Could keep “Require review from Code Owners” for application code.
- Cons: Easy to miss a path Dependabot touches (workflow YAML, `go.mod` / `go.sum`, `Dockerfile`).
- Rejected: Simpler to keep one global owner for review requests and gate merge with 1 review + CI.

## Consequences
- Human PRs still request `@ealebed`; they are not auto-approved.
- Major and `semver-unknown` updates stay open for manual review (typical for some Docker tags).
- `APP_CLIENT_ID` and `APP_PRIVATE_KEY` must exist in **both** Actions and Dependabot secret stores under identical names.
- Do not store an OAuth client secret; installation tokens need the App private key PEM.
- Client ID is read from secrets (not Actions variables) so Dependabot-triggered jobs can see it.
