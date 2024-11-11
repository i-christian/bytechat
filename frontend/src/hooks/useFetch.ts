export const getUser = async () => {
  const response = await fetch(`//${window.location.host}/api/auth/user`, {
    method: "GET",
    credentials: "include",
    mode: "cors",
  });

  if (!response.ok) throw new Error("failed to load user information");

  return response.json();
}


