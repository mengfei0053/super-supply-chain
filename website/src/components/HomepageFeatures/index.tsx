import type {ReactNode} from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  description: ReactNode;
};

const FeatureList: FeatureItem[] = [
  {
    title: '管理后台',
    description: (
      <>
        React Admin 界面位于 <code>/super-supply-chain/</code>
        ，用来处理结算单、动态 Excel、字典和读取规则。
      </>
    ),
  },
  {
    title: 'HTTP API',
    description: (
      <>
        Go / Gin 提供 <code>/api/login</code> 和受保护的 <code>/api/admin</code>
        。上传、查询、导出发票、删除，以及公司查询、新增和更新见 <Link to="/common-apis">常用接口</Link>。
      </>
    ),
  },
  {
    title: 'MCP 与 CLI',
    description: (
      <>
        <code>ssc-mcp/</code> 把同一套 API 暴露给编辑器。接入见{' '}
        <Link to="/mcp">MCP 接入与使用</Link>。<code>cli/</code>{' '}
        负责登录和状态检查。
      </>
    ),
  },
];

function Feature({title, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): ReactNode {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props) => (
            <Feature key={props.title} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
