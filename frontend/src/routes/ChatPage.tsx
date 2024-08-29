import { type ParentComponent } from "solid-js";
import ThemeToggle from "../component/theme/ThemeToggle";

const ChatPage: ParentComponent<{}> = (props) => {
  return (
    <main>
      <header>
        < ThemeToggle />
      </header>
      <aside>side bar </aside>
      {props.children}
      <footer>Footer</footer>
    </main>
  )
}

export default ChatPage;
