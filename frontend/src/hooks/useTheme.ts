import { createSignal, onMount } from 'solid-js';

export function useTheme() {
  const [isDark, setIsDark] = createSignal(false);

  const toggleTheme = () => {
    const newTheme = !isDark() ? 'dark' : 'light';
    setIsDark(!isDark());
    document.documentElement.classList.toggle('dark', newTheme === 'dark');
    localStorage.setItem('theme', newTheme);
  };

  onMount(() => {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'dark') {
      setIsDark(true);
      document.documentElement.classList.add('dark');
    } else if (savedTheme === 'light') {
      setIsDark(false);
      document.documentElement.classList.remove('dark');
    }
  });

  return { isDark, toggleTheme };
}
