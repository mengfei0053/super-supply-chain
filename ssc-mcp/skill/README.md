# SSC MCP skill source

`ssc-mcp-skill/` is the Cursor skill for the public Super Supply Chain MCP. The docs site serves a zip of that folder at `website/static/files/ssc-mcp-skill.zip` (published under the docs `baseUrl` as `/docs/files/ssc-mcp-skill.zip`).

After editing `ssc-mcp-skill/SKILL.md`, regenerate the zip from the repository root:

```bash
cd ssc-mcp/skill && rm -f ../../website/static/files/ssc-mcp-skill.zip && zip -r -X ../../website/static/files/ssc-mcp-skill.zip ssc-mcp-skill
```

The archive must contain `ssc-mcp-skill/SKILL.md`.

The docs page links through `website/src/components/SkillZipLink`, which points at `/docs/files/ssc-mcp-skill.zip`. A plain `/files/...` markdown link is rewritten to a hashed asset and gets a trailing slash, which does not download.
