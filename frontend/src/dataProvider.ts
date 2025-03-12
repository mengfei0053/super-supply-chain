import localforage from "localforage";
import simpleRestProvider from "ra-data-simple-rest";
import { CreateParams, DataProvider, fetchUtils } from "react-admin";
import { User } from "./authProvider";
import { CreateTemplateParams } from "./pages/excels/CreateTemplate";
import qs from "qs";

export const httpClient = async (
  url: string,
  options: fetchUtils.Options = {},
) => {
  const user = (await localforage.getItem("user")) as User;
  const headers = new Headers({
    ...options.headers,
    Authorization: `Bearer ${user?.token || ""}`,
  });

  return fetchUtils.fetchJson(url, {
    ...options,
    headers: headers,
  });
};

const baseDataProvider = simpleRestProvider(
  import.meta.env.VITE_JSON_SERVER_URL,
  httpClient,
);

const createPostFormData = (params: CreateParams<CreateTemplateParams>) => {
  const formData = new FormData();
  params.data.file?.rawFile &&
    formData.append("file", params.data.file.rawFile);
  params.data.alias && formData.append("alias", params.data.alias);
  return formData;
};

export const dataProvider: DataProvider = {
  ...baseDataProvider,
  getList: async (resource, params) => {
    console.log(params, "params");
    const { filter, sort, pagination, ...rest } = params;
    const current = pagination?.page || 1;
    const perPage = pagination?.perPage || 10;
    const start = (current - 1) * perPage;
    const end = current * perPage;

    const query = qs.stringify(
      {
        filter: JSON.stringify(filter),
        sort: JSON.stringify(sort),
        range: [start, end],
        ...rest,
      },
      {
        arrayFormat: "repeat",
      },
    );
    return httpClient(
      `${import.meta.env.VITE_JSON_SERVER_URL}/${resource}?${query}`,
      {
        method: "GET",
      },
    ).then(({ json, headers }) => {
      console.log(json, "json");
      const total = headers.get("content-range")
        ? parseInt(headers.get("content-range") as string)
        : 0;

      return {
        data: json,
        total: total,
      };
    });
  },
  create: async (resource, params) => {
    if (resource.match("excel-export-rule/template/")) {
      const formData = createPostFormData(params);
      return httpClient(`${import.meta.env.VITE_JSON_SERVER_URL}/${resource}`, {
        method: "POST",
        body: formData,
      }).then(({ json }) => ({ data: json }));
    }

    return baseDataProvider.create(resource, params);
  },
};
