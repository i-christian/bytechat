import { useNavigate } from "@solidjs/router";
import { Component } from "solid-js";

const Forbidden: Component = () => {
  const navigate = useNavigate();
  return (
    <section class="bg-inherit">
      <div class="flex sm:flex-wrap items-center min-h-screen px-6 py-12 mx-auto">
        <div>
          <p class="text-sm font-medium text-red-500 dark:text-red-400">403 Forbidden</p>
          <h1 class="mt-3 text-2xl font-semibold text-gray-800 dark:text-white md:text-3xl">Access Denied</h1>
          <p class="mt-4 text-gray-500 dark:text-gray-400">
            Sorry, you do not have the necessary permissions to view this page.
          </p>

          <div class="flex items-center mt-6 gap-x-3">
            <button
              class="flex items-center justify-center w-1/2 px-5 py-2 text-sm text-gray-700 transition-colors duration-200 bg-white border rounded-lg gap-x-2 sm:w-auto dark:hover:bg-gray-800 dark:bg-gray-900 hover:bg-gray-100 dark:text-gray-200 dark:border-gray-700"
              onClick={() => navigate("/")}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="w-5 h-5 rtl:rotate-180"
              >
                <path stroke-linecap="round" stroke-linejoin="round" d="M6.75 15.75L3 12m0 0l3.75-3.75M3 12h18" />
              </svg>

              <span>Go home</span>
            </button>
            <button
              class="w-1/2 px-5 py-2 text-sm text-gray-700 transition-colors duration-200 bg-white border rounded-lg sm:w-auto dark:bg-gray-900 hover:bg-gray-100 dark:text-gray-200 dark:border-gray-700"
              onClick={() => navigate(-1)}
            >
              Go Back
            </button>
          </div>
        </div>
      </div>
    </section>
  );
};

export default Forbidden;
