export function isNavItemActive(pathname: string, href: string): boolean {
  if (pathname === href) {
    return true;
  }

  if (href === "/") {
    return pathname === "/";
  }

  return pathname.startsWith(`${href}/`);
}

