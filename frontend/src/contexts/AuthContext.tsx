import { createContext, useContext, createSignal, Component, JSX } from "solid-js";

type AuthContextType = {
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType>();

export const AuthProvider: Component<{ children: JSX.Element }> = (props) => {
  const [isAuthenticated, setIsAuthenticated] = createSignal(false);

  const login = async (email: string, password: string) => {
    try {
      const response = await fetch(`//${window.location.host}/api/auth/login`, {
        method: "POST",
        mode: "cors",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      if (response.ok) {
        setIsAuthenticated(true);
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
