import { render } from "solid-js/web";
import { Route, Router, Navigate } from "@solidjs/router";
import { Component, createSignal, onMount } from "solid-js";
import ChatRoom from "./pages/ChatRoom";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Forbidden from "./pages/403";
import WrongPage from "./pages/404";

export const [isLoggedIn, setIsLoggedIn] = createSignal(false);

const App: Component = () => {
  onMount(async () => {
    if (localStorage.getItem("isLoggedIn") === "true") {
      try {
        const response = await fetch(`//${window.location.host}/api/check`, {
          method: "GET",
          credentials: "include"
        });

        if (response.ok) {
          setIsLoggedIn(true);
          localStorage.setItem("isLoggedIn", "true");
        } else {
          setIsLoggedIn(false);
          localStorage.removeItem("isLoggedIn");
        }
      } catch (error) {
        console.error("Auth check error:", error);
        setIsLoggedIn(false);
      }
    }
    setIsLoggedIn(false);
  });

  return (
    <Router>
      <Route path="/" component={() => isLoggedIn() ? <ChatRoom /> : <Navigate href="/login" />} />
      <Route path="/register" component={Register} />
      <Route path="/login" component={Login} />
      <Route path="/403" component={Forbidden} />
      <Route path="*" component={WrongPage} />
    </Router>
  );
};

export default App;

render(() => <App />, document.getElementById("root") as HTMLElement);
