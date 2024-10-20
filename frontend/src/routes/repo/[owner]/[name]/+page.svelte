<script>
  import SecurityIssuesList from "../../../../components/SecurityIssuesList.svelte";
  import { onMount } from 'svelte';

    export let data;

    const baseUrl = import.meta.env.VITE_CLIENT_API_URL;

    let installPressed = false;

    let languages = [];
    let scanComplete = false;
    let installComplete = false;
    let installStarted = false;

    const getToken = () => {
        try {          
          const cookie = document.cookie.split('; ').find(row => row.startsWith('at='));
          return cookie ? cookie.split('=')[1] : "";
        }catch(e){
          return "";
        }
    };

    

    const getScanState = () => {
      const token = getToken();

      fetch(`${baseUrl}/scan/status`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          owner: data.owner,
          repo: data.name,
        }),
      }).then((res) => res.json())
        .then((resp) => {
          if (!resp.error) {
            scanComplete = true;
            fetch(`${baseUrl}/scan/results`, {
              method: 'POST',
              headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({
                owner: data.owner,
                repo: data.name,
              })
            }).then((res) => res.json())
              .then((rdata) => {
                languages = rdata.languages;
              })
              .catch((err) => {
                console.error(err);
              });
          }else {
            if (resp.error == 'Repository not found') {
              fetch(`${baseUrl}/scan/start`, {
                method: 'POST',
                headers: {
                  'Authorization': `Bearer ${token}`,
                  'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                  owner: data.owner,
                  repo: data.name,
                }),
              })
            }
            setTimeout(getScanState, 1500);
          }
        })
        .catch((err) => {
          console.error(err);
        });
    }

    const getInstallState = () => {
      const token = getToken();

      fetch(`${baseUrl}/install/status`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          owner: data.owner,
          repo: data.name,
        }),
      }).then((res) => res.json())
        .then((resp) => {
          if (!resp.error) {
            installComplete = resp.installComplete;
            installStarted = resp.installStarted;
            if (installStarted && !installComplete) {
              setTimeout(getInstallState, 1500);
            }
          }else {
            setTimeout(getInstallState, 1500);
          }
        })
        .catch((err) => {
          console.error(err);
        });
    }

    const handleInstall = () => {
        installPressed = true;

        // get auth token from cookies (at)
        const token = getToken();

        installStarted = true;

      // post to the server to install the repository
      fetch(`${baseUrl}/install`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          owner: data.owner,
          repo: data.name,
        }),
      })
        .then((res) => res.json())
        .catch((err) => {
          console.error(err);
        });

        setTimeout(getInstallState,1500);

    };

    onMount(() => {
      getScanState();
      getInstallState();
    });

    console.log(data);
    function logout() {
        window.location.href = '/api/auth/logout';
    }
</script>

<div class="min-h-screen bg-gray-100 flex flex-col">

    <!-- Top Navbar -->
    <nav class="bg-white border-gray-200 px-4 py-3 sm:px-6 lg:px-8 shadow-sm">
        <div class="flex justify-between items-center mx-auto max-w-7xl">
            <!-- Brand Name -->
            <a href="/user" class="text-2xl font-bold text-blue-600">AutoLock</a>

            <!-- Logout Button -->
            <div>
                <button
                    class="text-white bg-red-600 hover:bg-red-700 focus:ring-4 focus:ring-red-300 font-medium rounded-lg text-sm px-4 py-2 text-center"
                    on:click={logout}>
                    Log Out
                </button>
            </div>
        </div>
    </nav>

    <main class="flex-grow flex flex-col items-center justify-center p-8">
      {#if scanComplete}
        {#if !installComplete}
          <!-- Display scan result -->
          <div class="text-center">
            <h1 class="text-3xl font-bold text-green-600 mb-4">
              {#if !installStarted}
              Scan Complete - Ready To Install
              {:else}
              Install Started
              {/if}
            </h1>

            <!-- Display repository name -->
            <h2 class="text-2xl font-semibold text-gray-800 mb-6">Repository: {data.name}</h2>

            <!-- Detected technologies -->
            <div class="w-full max-w-lg text-left">
              <h3 class="text-xl font-semibold text-gray-700 mb-4">Technologies Detected:</h3>
              {#if languages.length > 0}
                <ul class="list-disc list-inside">
                  {#each languages as language}
                    <li class="text-gray-700">
                      {language.Name} <span class="text-sm text-gray-500">({language.Extension})</span>
                    </li>
                  {/each}
                </ul>
              {:else}
                <p class="text-gray-600">No technologies detected.</p>
              {/if}
            </div>

            <!-- Install button (does nothing for now) -->
            {#if !installStarted}
            <button
              class="mt-6 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:ring-4 focus:ring-blue-300"
              on:click={handleInstall}>
              Install
            </button>
            {:else}
            <div class="w-full max-w-md mt-6">
              <div class="w-full bg-gray-300 rounded-full h-4 mb-4">
                <div class="bg-blue-600 h-4 rounded-full animate-pulse" style="width: 75%"></div>
              </div>
              <p class="text-sm text-gray-600">Your install has begun. Merge the PR to continue.</p>
            </div>
            {/if}
          </div>
        {:else}
          <div class="text-center">
            <h2 class="text-3xl font-bold text-green-600 mb-4">Security For {data.name}</h2>

            <SecurityIssuesList owner={data.owner} repo={data.name} />
          </div>
        {/if}
      {:else}
        <!-- Display loading bar while waiting for scan completion -->
        <div class="w-full max-w-md">
          <p class="text-xl font-semibold text-gray-700 mb-4">Scanning in progress...</p>
          <div class="w-full bg-gray-300 rounded-full h-4 mb-4">
            <div class="bg-blue-600 h-4 rounded-full animate-pulse" style="width: 75%"></div>
          </div>
          <p class="text-sm text-gray-600">Please wait while we scan your repository.</p>
        </div>
      {/if}
    </main>

    <!-- Footer (optional) -->
    <footer class="bg-gray-800 text-gray-400 text-center py-4">
        © 2024 AutoLock. All rights reserved.
    </footer>
</div>
      