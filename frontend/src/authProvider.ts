import { AuthProvider, HttpError, UserIdentity, fetchUtils } from "react-admin";
import localforage from "localforage";

export interface User {
  avatar: string;
  email: string;
  fullName: string;
  id: number;
  token: string;
  username: string;
}

/**
 * This authProvider is only for test purposes. Don't use it in production.
 */
export const authProvider: AuthProvider = {
  login: async ({ username, password }) => {
    const res = await fetchUtils.fetchJson(
      `${import.meta.env.VITE_LOGIN_URL}`,
      {
        method: "POST",
        body: JSON.stringify({ username, password }),
      },
    );
    if (res.json && res.json.token) {
      localforage.setItem("user", res.json);
      return Promise.resolve();
    }

    return Promise.reject(
      new HttpError("Unauthorized", 401, {
        message: "Invalid username or password",
      }),
    );
  },
  logout: () => {
    localforage.removeItem("user");
    return Promise.resolve();
  },
  checkError: () => Promise.resolve(),
  checkAuth: async () => {
    const user = await localforage.getItem("user");
    if (user) {
      return Promise.resolve();
    }

    window.location.href = "/super-supply-chain/#/login";

    return Promise.reject();
  },
  getPermissions: () => {
    return Promise.resolve(undefined);
  },
  getIdentity: async () => {
    const persistedUser = (await localforage.getItem("user")) as UserIdentity;

    return Promise.resolve(persistedUser);
  },
};

export default authProvider;
