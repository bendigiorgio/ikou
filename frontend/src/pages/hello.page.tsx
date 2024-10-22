import { wrapSSRComponent } from "@/serverEntry";
import React from "react";

const HelloPage = wrapSSRComponent(() => {
  return <div>HelloPage</div>;
});

export default HelloPage;
