import React from "react";
import ReactDOM from "react-dom/client";
import Root from "./root";

const root = ReactDOM.hydrateRoot(
  document.getElementById("app")!,
  //@ts-ignore
  <Root {...(window.APP_PROPS || {})} />
);
