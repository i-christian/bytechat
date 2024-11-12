import { createSignal } from "solid-js";
import { useNavigate } from "@solidjs/router";

type InputEvent = Event & { target: HTMLInputElement };

export default function Register() {
  const [name, setName] = createSignal<string>("");
  const [email, setEmail] = createSignal<string>("");
  const [pw, setPw] = createSignal<string>("");
  const [pwConfirm, setPwConfirm] = createSignal<string>("");
  const [pwVis, setPwVis] = createSignal<boolean>(false);
  const [error, setError] = createSignal<string>("");

  const navigate = useNavigate();

  const handleSubmit = async (e: Event) => {
    e.preventDefault();

    const url = `//${window.location.host}/api/auth/register`;

    try {
      const response = await fetch(url, {
        method: "POST",
        mode: "cors",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          name: name(),
          email: email(),
          password: pw(),
        }),
      });

      if (response.status === 400) {
        setError("Email already registered")
        return;
      }

      if (response.status === 201) {
        navigate("/login");
      }
    } catch (e) {
      console.error(`Error: ${e}`);
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
        <h1 class="lg:text-2xl text-xl text-center">Create an Account</h1>
        {error() && (
          <p class="mt-4 text-center text-red-600 text-sm animate-bounce">
            {error()}
          </p>
        )}

        <fieldset class="mt-10">
          <label for="name" class="text-xs tracking-wide text-gray-600">
            Name:
          </label>
          <div class="relative mb-4">
            <input
              type="text"
              name="name"
              class="text-sm placeholder-gray-500 pl-10 pr-4 rounded-md border border-gray-400 w-full py-2 focus:outline-none focus:border-black"
              required
              value={name()}
              onInput={(e: InputEvent) => setName(e.target.value)}
              placeholder="Enter your name"
            />
          </div>

          <label for="email" class="text-xs tracking-wide text-gray-600">
            Email:
          </label>
          <div class="relative mb-4">
            <input
              type="email"
              name="email"
              class="text-sm placeholder-gray-500 pl-10 pr-4 rounded-md border border-gray-400 w-full py-2 focus:outline-none focus:border-black"
              required
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
              required
              value={pw()}
              onInput={(e: InputEvent) => setPw(e.target.value)}
            />
            <button
              class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-500 text-sm"
              onClick={togglePassword}
              type="button"
            >
              {pwVis() ? "Hide" : "Show"}
            </button>
          </div>

          <label for="confirm" class="text-xs tracking-wide text-gray-600">
            Confirm Password:
          </label>
          <div class="relative mb-4">
            <input
              type={pwVis() ? "text" : "password"}
              name="confirm"
              class="text-sm placeholder-gray-500 pl-10 pr-4 rounded-md border border-gray-400 w-full py-2 focus:outline-none focus:border-black"
              required
              value={pwConfirm()}
              onInput={(e: InputEvent) => setPwConfirm(e.target.value)}
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
              class="flex mt-2 items-center justify-center focus:outline-none text-white text-sm sm:text-base bg-slate-800 hover:bg-slate-950 rounded-md py-2 w-full transition duration-150 ease-in"
            >
              <span class="mr-2 uppercase">Sign Up &rarr;</span>
            </button>
          </div>
        </fieldset>
      </div>

      <div class="flex justify-center items-center mt-6">
        <span class="inline-flex items-center text-gray-700 font-medium text-xs text-center">
          You have an account?
          <a href="/login" class="text-xs ml-2 text-black font-semibold">Login here</a>
        </span>
      </div>
    </form>
  );
}
