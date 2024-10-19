// src/routes/user/+page.server.js

export async function load({ cookies }) {
    const token = cookies.get('at');
    
    if (!token) {
      return {
        status: 401,
        error: new Error('Unauthorized: No token found'),
      };
    }
  
    const apiUrl = import.meta.env.VITE_API_URL + '/me';
    const apiUrlRepo = import.meta.env.VITE_API_URL + '/repo';
  
    try {
      // Make the request to the external API with the token
      const response = await fetch(apiUrl, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`, // Send JWT in Authorization header
          'Content-Type': 'application/json',
        },
      });
  
      // Check if the response is OK
      if (!response.ok) {
        throw new Error('Failed to fetch user information');
      }
  
      const data = await response.json();
  
      let repos = [];
      if (data.hasAccess) {
        const repoResponse = await fetch(apiUrlRepo, {
            method: 'GET',
            headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
            },
        });

        if (!repoResponse.ok) {
            throw new Error('Failed to fetch repositories');
        }

        const repoData = await repoResponse.json();
        repos = repoData.repos;
      }

      data.repos = repos;

      return data;

    } catch (error) {
      console.error('Error fetching user data:', error);
  
      // Handle any errors and return an error status
      return {
        status: 500,
        error: new Error('Failed to fetch user data'),
      };
    }
  }
  