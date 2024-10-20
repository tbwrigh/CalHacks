// src/routes/repo/[owner]/[name]/+page.server.js

// Helper to delay polling (to avoid tight loops)
function delay(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }
  
  export async function load({ params, cookies }) {
    const token = cookies.get('at');
    const { owner, name } = params; // Extract owner and repo name from URL

    console.log('owner:', owner);
    console.log('name:', name);
  
    if (!token) {
      return {
        status: 401,
        error: 'Unauthorized: No token found',
      };
    }
  
    return {
      owner: owner,
      name: name,
    }
  }
