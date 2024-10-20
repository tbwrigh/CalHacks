<script>
    export let repo = "";
    export let owner = "";
    export let token = "";
    
    let issues = [];

    let loading = true;

    import { onMount } from 'svelte';

    const getToken = () => {
        try {          
          const cookie = document.cookie.split('; ').find(row => row.startsWith('at='));
          return cookie ? cookie.split('=')[1] : "";
        }catch(e){
          return "";
        }
    };

    let loadIssues = () => {
        console.log('Loading issues...');
        const baseUrl = import.meta.env.VITE_CLIENT_API_URL;
        const token = getToken();

        fetch(`${baseUrl}/issues/list`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify({
                owner: owner,
                repo: repo,
            }),
        }).then((res) => res.json())
            .then((data) => {
                if (!data.error) {
                    issues = data;
                    console.log(data);  
                    loading = false;                    
                }else{
                    setTimeout(loadIssues, 1500);
                }
            })
            .catch((err) => {
                console.log(err);
            });
    }

    onMount(() => {
        loadIssues();
    });

    const fixIssue = (issue) => {
        const baseUrl = import.meta.env.VITE_CLIENT_API_URL;
        const token = getToken();

        fetch(`${baseUrl}/issues/resolve`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify({
                id: issue,
            }),
        }).then((res) => res.json())
            .then((data) => {
                if (!data.error) {
                    console.log(data);
                    loadIssues();
                }else{
                    console.log(data.error);
                }
            })
            .catch((err) => {
                console.log(err);
            });
    }

  
</script>
  
{#if loading}
    <div class="w-full max-w-md">
        <p class="text-xl font-semibold text-gray-700 mb-4">Loading...</p>
        <div class="w-full bg-gray-300 rounded-full h-4 mb-4">
            <div class="bg-blue-600 h-4 rounded-full animate-pulse" style="width: 75%"></div>
        </div>
        <p class="text-sm text-gray-600">Please wait while we fetch the security issues.</p>
    </div>
{:else}
  <!-- Security Issues List -->
  <div class="grid grid-cols-1 gap-6">
    {#if issues.length > 0}
      {#each issues as issue}
      <div class="bg-white shadow-md rounded-lg p-6 w-full mb-4">
        <div class="flex justify-between items-start">
          <div class="text-left">
            <h3 class="text-lg font-bold text-gray-900 mb-2">Issue #{issue.GithubNumber}</h3>
            <p class="text-sm text-gray-600">Path: {issue.Path || 'N/A'}</p>
            <p class="text-sm text-gray-600">Lines: {issue.StartLine} - {issue.EndLine}</p>
            <p class="text-sm text-gray-600">Description: {issue.FullDescription || 'No description provided'}</p>
            <p class="text-sm text-gray-600">Repository: {issue.Repository.Name} ({issue.Repository.Owner})</p>
          </div>
          {#if !issue.FixSuggested}
          <button 
            class="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 transition-colors" 
            on:click={() => {fixIssue(issue.ID)}}>
            Fix ✨
          </button>
          {:else}
          <h3 class="text-lg font-bold text-gray-900 mb-2">PR Opened ✨</h3>
          {/if}
        </div>
      </div>         
      {/each}
    {:else}
      <p class="text-gray-600">No security issues found.</p>
    {/if}
  </div>
{/if}