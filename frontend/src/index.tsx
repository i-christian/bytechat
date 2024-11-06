import { render } from "solid-js/web";
import App from "./App.tsx";
import { AuthProvider } from "./contexts/AuthContext";

render(
  () => (
    <AuthProvider>
      <App />
    </AuthProvider>
  ),
  document.getElementById("root") as HTMLElement
);
