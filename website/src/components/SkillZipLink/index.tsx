import type {ReactNode} from 'react';
import useBaseUrl from '@docusaurus/useBaseUrl';

// Plain anchor. A markdown link to the zip is rewritten into a hashed asset
// and receives a trailing slash, which does not download. useBaseUrl keeps
// the published path at /docs/files/ssc-mcp-skill.zip.
export default function SkillZipLink(): ReactNode {
  const href = useBaseUrl('/files/ssc-mcp-skill.zip');
  return (
    <strong>
      <a href={href}>下载 SSC MCP Skill</a>
    </strong>
  );
}
