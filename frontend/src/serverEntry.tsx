import * as React from "react";
import { renderToString } from "react-dom/server";
import Root from "./root";
import "./styles/base.css";

globalThis.React = React;

export const wrapSSRComponent = (PageComponent: any) => {
  const renderApp = (props: any) => {
    const renderedHTML = renderToString(
      <Root>
        <PageComponent {...props} />
      </Root>
    );
    return renderedHTML;
  };
  // @ts-ignore
  globalThis.renderApp = renderApp;
};
