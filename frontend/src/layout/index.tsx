import type { ReactNode } from "react";
import { Layout as RALayout, CheckForApplicationUpdate } from "react-admin";
import { SscAppBar } from "./AppBar";
import Menu from "./Menu";

export const Layout = ({ children }: { children: ReactNode }) => (
  <RALayout appBar={SscAppBar} menu={Menu}>
    {children}
    <CheckForApplicationUpdate />
  </RALayout>
);
