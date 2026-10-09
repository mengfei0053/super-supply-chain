import { AppBar, TitlePortal } from "react-admin";
import { Button } from "@mui/material";

// Docs are hosted at the site root (/docs/), outside the SPA base
// /super-supply-chain/. A plain href performs a full navigation.
export const SscAppBar = () => (
  <AppBar>
    <TitlePortal />
    <Button
      color="inherit"
      href="/docs/"
      sx={{ fontWeight: 600, whiteSpace: "nowrap" }}
    >
      文档
    </Button>
  </AppBar>
);
