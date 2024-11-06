import { createContext, useContext, createSignal, Component, JSX } from "solid-js";
import { useNavigate } from "@solidjs/router";

type AuthContextType = {
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType>();

export const AuthProvider: Component<{ children: JSX.Element }> = (props) => {
  const navigate = useNavigate();
  const [isAuthenticated, setIsAuthenticated] = createSignal(false);

  const login = async (email: string, password: string) => {
    try {
      const response = await fetch(`//${window.location.host}/api/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      if (response.ok) {
        setIsAuthenticated(true);
        navigate("/");
      } else {
        alert("Login failed");
      }
    } catch (error) {
      console.error("Login error:", error);
    }
  };

  const logout = () => {
    fetch(`//${window.location.host}/api/auth/logout`, { method: "GET", credentials: "include" })
      .then(() => {
        setIsAuthenticated(false);
        navigate("/login");
      })
      .catch((error) => console.error("Logout error:", error));
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated: isAuthenticated(), login, logout }}>
      {props.children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
