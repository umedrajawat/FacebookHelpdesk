import axios from "axios";
const baseURL = process.env.REACT_APP_API_URL || window.API_URL;
const graphURL = "https://graph.facebook.com/v19.0";

console.log('REACT_APP_GRAPH_URL', graphURL, process.env.REACT_APP_API_URL, baseURL);

let GraphApi = axios.create({
  baseURL: graphURL,
  withCredentials: false,
});

let Api = axios.create({
  baseURL: "http://localhost:8081",
  withCredentials: false,
  headers: {
    Authorization: `Bearer ${localStorage.getItem("AUTH_TOKEN")}`,
    'x-user': "kumar Umed Rajawat",
    "Access-Control-Allow-Origin": "*",
    "Access-Control-Allow-Headers": "*",
    "Access-Control-Allow-Methods": "*",
    Accept: "*/*",
  },
});

let resetApiHeaders = (token) => {
  if (!token || token === "") {
    axios.defaults.headers.common["Authorization"] = null;
    Api = axios.create({
      baseURL: baseURL,
    });
  }
  Api = axios.create({
    baseURL: baseURL,
    headers: {
      Authorization: `Bearer ${token}`,
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Headers": "*",
      "Access-Control-Allow-Methods": "*",
      Accept: "*/*",
    },
  });
};

export { Api, resetApiHeaders, GraphApi };
