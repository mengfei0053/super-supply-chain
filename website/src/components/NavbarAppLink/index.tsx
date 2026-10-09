import type {ReactNode} from 'react';
import clsx from 'clsx';

type Props = {
  className?: string;
  label?: string;
  href: string;
  mobile?: boolean;
  onClick?: () => void;
};

// Plain anchor so Docusaurus does not prefix baseUrl (/docs/) onto a host-root path.
export default function NavbarAppLink({
  className,
  label,
  href,
  mobile = false,
  onClick,
}: Props): ReactNode {
  const link = (
    <a
      className={clsx(
        mobile ? 'menu__link' : 'navbar__item navbar__link',
        className,
      )}
      href={href}
      onClick={onClick}>
      {label}
    </a>
  );
  if (mobile) {
    return <li className="menu__list-item">{link}</li>;
  }
  return link;
}
