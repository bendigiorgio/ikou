import { wrapSSRComponent } from "@/serverEntry";

const TestPage = wrapSSRComponent(() => {
  return <main>TestPage</main>;
});

export default TestPage;
