import React from "react";

const AltPage = () => {
  const [count, setCount] = React.useState(0);
  return (
    <div>
      AltPage
      <button onClick={() => setCount(count + 1)}>Click me</button>
      <div>{count}</div>
    </div>
  );
};

export default AltPage;
