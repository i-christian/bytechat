import { Route, Router } from "@solidjs/router";
import { Component, lazy } from "solid-js";
import Login from "./pages/Login";
import Logout from "./pages/Logout";

const Home: Component = lazy(() => import("./pages/Home"));
const Forbidden: Component = lazy(() => import("./pages/403"));
const WrongPage: Component = lazy(() => import("./pages/404"));


const App: Component = () => {
  return (
    <Router>
      <Route path="/" component={Home} />
      <Route path="/login" component={Login} />
      <Route path="/logout" component={Logout} />
      <Route path="/403" component={Forbidden} />
      <Route path="*404" component={WrongPage} />
    </Router>
  );
};

export default App;
