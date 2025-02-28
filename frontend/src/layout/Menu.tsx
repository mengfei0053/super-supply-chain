import {
  MenuProps,
  useSidebarState,
  DashboardMenuItem,
  useTranslate,
  MenuItemLink,
  useStore,
} from "react-admin";
import { Box } from "@mui/material";
import ProductIcon from "@mui/icons-material/Collections";
import OrderIcon from "@mui/icons-material/AttachMoney";
import * as React from "react";
import SubMenu from "./SubMenu";
import { httpClient } from "../dataProvider";

type IMenuProps = MenuProps;

type MenuName = "menuCatalog" | "menuSales" | "menuCustomers";

const Menu: React.FunctionComponent<IMenuProps> = ({ dense = false }) => {
  const [state, setState] = React.useState({
    menuCatalog: true,
    menuSales: true,
    menuCustomers: true,
  });
  const translate = useTranslate();
  const [open] = useSidebarState();
  const [menus, setMenus] = useStore<
    { id: number; dynamicTableName: string; menuName: string }[]
  >("global.menu", []);

  const getMenus = async () => {
    const res = await httpClient(
      import.meta.env.VITE_JSON_SERVER_URL + "/menus",
    );
    setMenus(res.json);
  };

  const handleToggle = (menu: MenuName) => {
    setState((state) => ({ ...state, [menu]: !state[menu] }));
  };

  React.useEffect(() => {
    getMenus();
  }, []);

  return (
    <Box
      sx={{
        width: open ? 200 : 50,
        marginTop: 1,
        marginBottom: 1,
        transition: (theme) =>
          theme.transitions.create("width", {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
          }),
      }}
    >
      <DashboardMenuItem></DashboardMenuItem>

      <SubMenu
        handleToggle={() => handleToggle("menuCatalog")}
        isOpen={state.menuCatalog}
        name="翊帆"
        icon={<ProductIcon />}
        dense={dense}
      >
        <MenuItemLink
          to={`/yifan/cost-calculation`}
          state={{ _scrollToTop: true }}
          primaryText={translate("费用结算", {
            smart_count: 2,
          })}
          leftIcon={<OrderIcon />}
          dense={dense}
        />
      </SubMenu>
      <SubMenu
        handleToggle={() => handleToggle("menuCatalog")}
        isOpen={state.menuCatalog}
        name="Excel处理"
        icon={<ProductIcon />}
        dense={dense}
      >
        {menus.map((i) => {
          return (
            <MenuItemLink
              key={i.id}
              to={`/excel/${i.dynamicTableName}`}
              state={{ _scrollToTop: true }}
              primaryText={translate(i.menuName, {
                smart_count: 2,
              })}
              leftIcon={<OrderIcon />}
              dense={dense}
            />
          );
        })}
      </SubMenu>
      <SubMenu
        handleToggle={() => handleToggle("menuCatalog")}
        isOpen={state.menuCatalog}
        name="工具"
        icon={<ProductIcon />}
        dense={dense}
      >
        <MenuItemLink
          to="/dict-manage"
          state={{ _scrollToTop: true }}
          primaryText={translate(`字典管理`, {
            smart_count: 2,
          })}
          leftIcon={<OrderIcon />}
          dense={dense}
        />
        <MenuItemLink
          to="/excel-read-rules"
          state={{ _scrollToTop: true }}
          primaryText={translate(`Excel读取规则`, {
            smart_count: 2,
          })}
          leftIcon={<OrderIcon />}
          dense={dense}
        />

        <MenuItemLink
          to="/excel-mapping-rule"
          state={{ _scrollToTop: true }}
          primaryText={translate(`Excel映射规则`, {
            smart_count: 2,
          })}
          leftIcon={<OrderIcon />}
          dense={dense}
        />
      </SubMenu>
    </Box>
  );
};

export default Menu;
