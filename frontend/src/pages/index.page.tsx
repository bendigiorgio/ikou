import { wrapSSRComponent } from "@/serverEntry";

const HomePage = wrapSSRComponent(() => {
  return <main className="text-xl text-red-500">HomePage</main>;
});

export default HomePage;
