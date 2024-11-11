import { Component, Match, Suspense, Switch, createResource, createSignal, onMount } from "solid-js";
import { For } from "solid-js/web";
import { useNavigate } from "@solidjs/router";
import { setIsLoggedIn } from "../index";
import { getUser } from "../hooks/useFetch";

interface Room {
  id: string;
  name: string;
  description: string;
}

const ChatRoom: Component = () => {
  const [rooms, setRooms] = createSignal<Room[]>([]);
  const [selectedRoom, setSelectedRoom] = createSignal<Room | null>(null);
  const [error, setError] = createSignal("");
  const navigate = useNavigate();
  const [user] = createResource(getUser)

  const fetchRooms = async () => {
    try {
      const response = await fetch(`//${window.location.host}/api/rooms`, {
        method: "GET",
        credentials: "include",
      });

      if (response.status === 403) {
        setIsLoggedIn(false);
        navigate("/login");
        return;
      }

      if (!response.ok) throw new Error("Failed to load rooms");

      const data = await response.json();
      setRooms(data);
    } catch (err) {
      setError("Failed to load rooms. Please try again later.");
    }
  };

  const joinRoom = async (roomId: string) => {
    try {
      const response = await fetch(`//${window.location.host}/api/rooms`, {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ room_id: roomId }),
      });

      if (response.status === 403) {
        localStorage.removeItem("isLoggedin");
        setIsLoggedIn(false);
        navigate("/login");
        return;
      }

      if (!response.ok) throw new Error("Failed to join room");

      const joinedRoom = rooms().find((room) => room.id === roomId);
      setSelectedRoom(joinedRoom || null);
    } catch (err) {
      setError("Could not join the room. Try again.");
    }
  };

  const logout = async () => {
    try {
      const response = await fetch(`//${window.location.host}/api/auth/logout`, {
        method: "GET",
        credentials: "include",
      });

      if (response.status === 500) {
        console.log("Internal server error");
      } else {
        setIsLoggedIn(false);
        localStorage.removeItem("isLoggedIn");
        navigate("/login");

      }

    } catch (err) {
      console.error("Error logging out:", err);
    }
  };

  onMount(fetchRooms);

  return (
    <div class="flex h-screen">
      <aside class="w-1/4 bg-gray-200 p-4 border-r border-gray-300">
        <h2 class="text-xl text-center">
          <Suspense fallback={<div>Loading...</div>}>
            <Switch>
              <Match when={user.error}>
                <span>Error: {user.error.message}</span>
              </Match>
              <Match when={user()}>
                <p>{user().name}</p>
              </Match>
            </Switch>
          </Suspense>

        </h2>
        <h2 class="text-xl font-semibold mb-4">Chats</h2>
        {error() && <p class="text-red-500">{error()}</p>}
        <button onClick={logout}>
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6a2.25 2.25 0 0 0-2.25 2.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15M12 9l-3 3m0 0 3 3m-3-3h12.75" />
          </svg>

        </button>
        <ul class="space-y-2">
          <For each={rooms()}>
            {(room) => (
              <li
                class="p-2 bg-white rounded-md shadow cursor-pointer hover:bg-gray-100"
                onClick={() => joinRoom(room.id)}
              >
                <h3 class="text-lg font-medium">{room.name}</h3>
                <p class="text-sm text-gray-600">{room.description}</p>
              </li>
            )}
          </For>
        </ul>
      </aside>
      <main class="flex-1 flex flex-col p-4">
        {selectedRoom() ? (
          <>
            <h2 class="text-2xl font-bold">{selectedRoom()!.name}</h2>
            <p class="text-sm text-gray-500 mb-4">{selectedRoom()!.description}</p>
            <div class="flex-1 bg-gray-100 rounded-lg p-4 overflow-y-auto">
              <p>Chat messages for {selectedRoom()!.name} will appear here...</p>
            </div>
            <div class="mt-4">
              <input
                type="text"
                placeholder="Type your message here..."
                class="w-full p-2 border rounded-lg"
              />
              <button class="mt-2 p-2 bg-blue-600 text-white rounded-lg">
                Send
              </button>
            </div>
          </>
        ) : (
          <div class="flex items-center justify-center h-full">
            <p class="text-gray-500">Select a room to start chatting</p>
          </div>
        )}
      </main>
    </div>
  );
};

export default ChatRoom;
