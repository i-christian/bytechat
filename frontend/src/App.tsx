import { Route, Router } from "@solidjs/router";
import { lazy, type Component } from "solid-js";

const ChatPage: Component = lazy(() => import("./routes/ChatPage"));
const Messages: Component = lazy(() => import("./routes/Messages"));
const Login: Component = lazy(() => import("./routes/Login"));
const Register: Component = lazy(() => import("./routes/Register"));


const App: Component = () => {
  return (
    <Router>
      <Route path="/" component={ChatPage}>
        <Route path="/" component={Messages} />
      </Route>
      <Route path="/login" component={Login} />
      <Route path="/register" component={Register} />
    </Router>
  )
}

export default App;
