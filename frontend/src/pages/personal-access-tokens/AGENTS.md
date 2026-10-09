# Personal Access Tokens Page Guide

Scope: `frontend/src/pages/personal-access-tokens/`.

This resource creates, lists, and revokes the signed-in user's personal access tokens.

- The plaintext token is shown only on the create page after a successful response. Do not put it in the URL, local storage, or the React Admin cache.
- List rows come from `GET /api/admin/personal-access-tokens` and must not expect a `token` field.
- Revoke uses the React Admin delete action against `DELETE /api/admin/personal-access-tokens/:id`.
- Keep labels in Chinese.
