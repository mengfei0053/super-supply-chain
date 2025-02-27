import type { ReactNode } from "react";
import { Layout as RALayout, CheckForApplicationUpdate } from "react-admin";
import Menu from "./Menu";

export const Layout = ({ children }: { children: ReactNode }) => (
  <RALayout menu={Menu}>
    {children}
    <CheckForApplicationUpdate />
  </RALayout>
);
