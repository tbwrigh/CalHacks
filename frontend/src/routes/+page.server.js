// src/routes/+page.server.js

export async function load({ cookies }) {
  const token = cookies.get('at');
  console.log('Access Token:', token); // Log the access token for debugging

  let loggedIn = false;

  if (token) {
      loggedIn = true; // Token is valid
  }

  return { "logged_in": loggedIn };
}
