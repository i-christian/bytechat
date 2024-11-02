import { render } from "solid-js/web";
import App from "./App.tsx";

// const BASE_URL = import.meta.env.VITE_API_URL

render(
  () => (
    <App />
  ),
  document.getElementById("root") as HTMLElement
);
