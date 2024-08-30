import { A } from "@solidjs/router";
import { type Component } from "solid-js";

const Register: Component = () => {
  return (
    <section class="bg-inherit flex justify-center items-center min-h-screen mx-auto">
      <form class="px-10 py-5">
        <h1 class="text-xl mx-auto w-fit font-bold">Create an account </h1>
        <label for="fullname"></label>
        <input
          placeholder="full name"
          type="text"
          id="fullname"
          class="p-2 rounded-md my-5 block"
        />

        <label for="username"></label>
        <input
          placeholder="phone number"
          type="text"
          id="username"
          class="p-2 rounded-md my-5 block"
        />
        <label for="password"></label>
        <input
          placeholder="new password"
          type="password"
          id="password"
          class="p-2 rounded-md mb-5 block"
        />
        <label for="confirm"></label>
        <input
          placeholder="confirm password"
          type="password"
          id="password"
          class="p-2 rounded-md mb-5 block"
        />

        <button class="px-2 py-1 mx-auto btn w-fit">Sign Up</button>
        <p class="text-sm pt-10">
          Already have an account,
          <A href="/login" class="px-1 text-normal underline text-blue-800 hover:no-underline active:text-green-500"> Sign In</A>
        </p>
      </form>
    </section>
  )
}

export default Register;
