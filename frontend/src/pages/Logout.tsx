import { Component, onMount } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { useAuth } from "../contexts/AuthContext";

const Logout: Component = () => {
  const auth = useAuth();
  const navigate = useNavigate();

  onMount(() => {
    auth?.logout();
    navigate("/login");
  });

  return (
    <main>
      <p>Logging out...</p>
    </main>
  );
};

export default Logout;
