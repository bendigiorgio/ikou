import { wrapSSRComponent } from "@/serverEntry";

const AltPage = wrapSSRComponent(() => {
  return <div>AltPage</div>;
});

export default AltPage;
