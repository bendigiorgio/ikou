import React from "react";
import ReactDOM from "react-dom/client";
import Root from "./root";

const renderClientSide = (PageComponent: React.FC<any>, props: any) =>
  ReactDOM.hydrateRoot(
    document.getElementById("app")!,
    //@ts-ignore
    <Root {...(window.APP_PROPS || {})}>
      <PageComponent {...props} />
    </Root>
  );

// @ts-ignore
globalThis.renderClientSide = renderClientSide;
