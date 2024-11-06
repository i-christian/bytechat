import { Component, createSignal } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { setIsLoggedIn } from "../index";

type InputEvent = Event & { target: HTMLInputElement };

const Login: Component = () => {
  const [email, setEmail] = createSignal("");
  const [pw, setPw] = createSignal("");
  const [pwVis, setPwVis] = createSignal(false);
  const navigate = useNavigate();

  const checkAuthStatus = async () => {
    try {
      const response = await fetch(`//${window.location.host}/api/check`, {
        method: "GET",
        credentials: "include"
      });

      if (response.ok) {
        setIsLoggedIn(true);
        navigate("/");
      }
    } catch (error) {
      console.error("Auth check error:", error);
    }
  };

  if (localStorage.getItem("isLoggedIn") === "true") {
    checkAuthStatus();
  }


  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    const url = `//${window.location.host}/api/auth/login`;

    try {
      const response = await fetch(url, {
        method: "POST",
        mode: "cors",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: email(),
          password: pw(),
        }),
      });

      if (response.ok) {
        setIsLoggedIn(true);
        localStorage.setItem("isLoggedIn", "true");
        navigate("/");
      } else {
        alert("Login failed");
      }
    } catch (error) {
      console.error("Login error:", error);
      alert("An error occurred while logging in.");
    }
  };

  const togglePassword = (e: Event) => {
    e.preventDefault();
    setPwVis(!pwVis());
  };

  return (
    <form
      class="px-5 min-h-screen flex flex-col items-center justify-center bg-gray-100"
      onSubmit={handleSubmit}
    >
      <div class="flex flex-col bg-white shadow-md px-4 sm:px-6 md:px-8 lg:px-10 py-8 rounded-md w-50 max-w-md">
        <h1 class="lg:text-2xl text-xl text-center">Log In</h1>

        <fieldset class="mt-10">
          <label for="email" class="text-xs tracking-wide text-gray-600">
            Email Address:
          </label>
          <div class="relative mb-4">
            <input
              type="email"
              name="email"
              class="text-sm placeholder-gray-500 pl-10 pr-4 rounded-md border border-gray-400 w-full py-2 focus:outline-none focus:border-black"
              value={email()}
              onInput={(e: InputEvent) => setEmail(e.target.value)}
              placeholder="Enter your email"
            />
          </div>

          <label for="password" class="text-xs tracking-wide text-gray-600">
            Password:
          </label>
          <div class="relative mb-4">
            <input
              type={pwVis() ? "text" : "password"}
              name="password"
              class="text-sm placeholder-gray-500 pl-10 pr-4 rounded-md border border-gray-400 w-full py-2 focus:outline-none focus:border-black"
              value={pw()}
              onInput={(e: InputEvent) => setPw(e.target.value)}
              placeholder="Enter your password"
            />
            <button
              class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-500 text-sm"
              onClick={togglePassword}
              type="button"
            >
              {pwVis() ? "Hide" : "Show"}
            </button>
          </div>

          <div class="flex w-full">
            <button
              type="submit"
              class="flex mt-2 items-center justify-center focus:outline-none text-white text-sm sm:text-base bg-black hover:bg-slate-950 rounded-md py-2 w-full transition duration-150 ease-in"
            >
              <span class="mr-2 uppercase">Log In &rarr;</span>
            </button>
          </div>
        </fieldset>
      </div>

      <div class="flex justify-center items-center mt-6">
        <span class="inline-flex items-center text-gray-700 font-medium text-xs text-center">
          Don't have an account?
          <a href="/register" class="text-xs ml-2 text-black font-semibold">Register here</a>
        </span>
      </div>
    </form>
  );
};

export default Login;
