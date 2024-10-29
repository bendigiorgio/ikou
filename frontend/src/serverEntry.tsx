import * as React from "react";
import { renderToString } from "react-dom/server";
import Root from "./root";
import "./styles/base.css";

globalThis.React = React;

// @ts-ignore
globalThis.renderApp = (PageComponent: React.FC<any>, props: any) => {
  const renderedHTML = renderToString(
    <Root>
      <PageComponent {...props} />
    </Root>
  );
  return renderedHTML;
};
