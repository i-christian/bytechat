import { Route, Router } from "@solidjs/router";
import { Component, lazy, Show } from "solid-js";
import { useAuth } from "./contexts/AuthContext";
import Login from "./pages/Login";
import Logout from "./pages/Logout";
import Register from "./pages/Register";

const Home = lazy(() => import("./pages/Home"));
const Forbidden = lazy(() => import("./pages/403"));
const WrongPage = lazy(() => import("./pages/404"));

const ProtectedRoute: Component<{ component: Component }> = (props) => {
  const auth = useAuth();
  return (
    <Show when={auth?.isAuthenticated} fallback={<Login />}>
      {props.component}
    </Show>
  );
};

const App: Component = () => {
  return (
    <Router>
      <Route path="/" component={() => <ProtectedRoute component={Home} />} />
      <Route path="/register" component={Register} />
      <Route path="/login" component={Login} />
      <Route path="/logout" component={Logout} />
      <Route path="/403" component={Forbidden} />
      <Route path="*" component={WrongPage} />
    </Router>
  );
};

export default App;
