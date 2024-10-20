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
  
    const apiUrl = import.meta.env.VITE_API_URL;
    let scanComplete = false;
    let successMessage = '';
  
    try {
      // Step 1: Try to fetch scan results
      const resultResponse = await fetch(`${apiUrl}/scan/results`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          owner: owner,
          repo: name,
        }),
      });
  
      if (resultResponse.ok) {
        // If scan results exist, we are done
        successMessage = 'Success: Scan already completed.';
        scanComplete = true;

        const resultData = await resultResponse.json();

        const installResponse = await fetch(`${apiUrl}/install/status`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              owner: owner,
              repo: name,
            }),
        });

        const installData = await installResponse.json();

        return { 
          scanComplete: scanComplete, 
          successMessage: successMessage, 
          resultData: resultData, 
          name: name, 
          owner: owner,
          installStarted: installData.installStarted,
          installComplete: installData.installComplete,
        };
      }
  
      // Step 2: If scan results don't exist, start a new scan
      const startResponse = await fetch(`${apiUrl}/scan/start`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            owner: owner,
          repo: name,
        }),
      });
  
      if (!startResponse.ok) {
        throw new Error('Failed to start scan');
      }
  
      // Step 3: Poll for scan status until "scanComplete" is true
      let pollInterval = 2000; // 5 seconds interval
      while (!scanComplete) {
        await delay(pollInterval);
  
        const statusResponse = await fetch(`${apiUrl}/scan/status`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            owner: owner,
            repo: name,
          }),
        });
  
        const statusData = await statusResponse.json();
        if (statusData.scanComplete) {
          scanComplete = true;
          successMessage = 'Success: Scan completed.';
        }
      }
  
      // After scan completes, return success
      return {  
        scanComplete: scanComplete, 
        successMessage: successMessage, 
        resultData: {}, 
        name: name, 
        owner: owner,
        installStarted: false,
        installComplete: false, 
      };
  
    } catch (error) {
      console.error('Error during scan process:', error);
      return {
        status: 500,
        error: 'Failed to complete scan process',
      };
    }
  }
