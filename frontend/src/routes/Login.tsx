import { A } from "@solidjs/router";
import { type Component } from "solid-js";

const Login: Component = () => {
  return (
    <section class="bg-inherit flex justify-center items-center min-h-screen mx-auto">
      <form class="shadow-2xl rounded-lg px-10 py-5">
        <h1 class="text-xl mx-auto w-fit font-bold">Welcome back! </h1>
        <label for="username"></label>
        <input
          placeholder="phone number"
          type="text"
          id="username"
          class="p-2 rounded-md my-5 block"
        />
        <label for="password"></label>
        <input
          placeholder="password"
          type="password"
          id="password"
          class="p-2 rounded-md mb-5 block"
        />
        <button class="px-2 py-1 mx-auto btn w-fit">Sign In</button>
        <p class="text-sm pt-10">
          Don't have an account yet,
          <A href="/register" class="px-1 text-normal underline text-blue-800 hover:no-underline active:text-green-500"> Sign Up</A>
        </p>
      </form>
    </section>
  )
}

export default Login;
