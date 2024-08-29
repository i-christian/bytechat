import { Component } from 'solid-js';
import { useTheme } from '../../hooks/useTheme';

const ThemeToggle: Component = () => {
  const { isDark, toggleTheme } = useTheme();

  return (
    <div class="mt-4 flex items-center space-x-4">
      <div
        role="switch"
        aria-checked={isDark()}
        class="relative inline-flex items-center cursor-pointer"
        onClick={toggleTheme}
        aria-label="Toggle dark mode"
      >
        <input
          type="checkbox"
          checked={isDark()}
          onChange={toggleTheme}
          class="sr-only"
          aria-label="Toggle dark mode"
        />
        <div class="w-20 h-10 px-2 bg-gray-200 dark:bg-gray-600 rounded-full peer flex items-center relative">
          <div
            class={`absolute w-8 h-8 bg-white border border-gray-300 rounded-full shadow-md transition-transform duration-300 ${isDark() ? 'translate-x-10' : 'translate-x-0'
              }`}
          />
          <svg
            class={`absolute h-6 w-6 text-yellow-500 transition-opacity duration-300 ${isDark() ? 'opacity-0' : 'opacity-100'}`}
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
          >
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v2.25m6.364.386-1.591 1.591M21 12h-2.25m-.386 6.364-1.591-1.591M12 18.75V21m-4.773-4.227-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0Z" />
          </svg>
          <svg
            class={`absolute h-6 w-6 text-blue-500 transition-opacity duration-300 ${isDark() ? 'opacity-100' : 'opacity-0'}`}
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
          >
            <path stroke-linecap="round" stroke-linejoin="round" d="M21.752 15.002A9.72 9.72 0 0 1 18 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 0 0 3 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 0 0 9.002-5.998Z" />
          </svg>
        </div>
      </div>
    </div>
  );
};

export default ThemeToggle;
