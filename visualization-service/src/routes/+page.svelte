<script>
  import { ChartTheme, LineChart, ScaleTypes } from "@carbon/charts-svelte";
  export let data;
  export let users = Object.keys(data.results);
  let currentUser = users[0];

  $: currentDataSet = data.results[currentUser]
  let selectedImage = -1;
</script>

<div class="container mx-auto py-4">
  <div class="flex flex-row justify-between">
    <h1 class="text-6xl">Visualization Service Mock</h1>

    <details class="dropdown">
      {#if currentUser === "None"}
        <summary class="m-1 btn">Select User</summary>
      {:else}
        <summary class="m-1 btn">{currentUser}</summary>
      {/if}
      <ul
        class="p-2 shadow menu dropdown-content z-[1] bg-base-100 rounded-box w-52"
      >
        {#each users as user}
          <li><a href="#" on:click={() => (currentUser = user)}>{user}</a></li>
        {/each}
      </ul>
    </details>
  </div>
  {#if currentUser !== "None"}
    <hr class="mt-2 mb-4" />
    <div>
      <h2 class="text-4xl">Select Image</h2>
      <div class="flex flex-row gap-8 justify-center p-4">
        {#each currentDataSet as result, idx}
          <div>
            <img
              class={selectedImage === idx
                ? "border-4 border-primary hover:border-4"
                : "hover:border-4"}
              on:click={() => (selectedImage = idx)}
              src={`${data.s3Url}${result.media.id}`}
            />
          </div>
        {/each}
      </div>
    </div>
  {/if}

  {#if selectedImage >= 0}
    <div>
      <h2 class="text-4xl">Image Insights</h2>
      <LineChart
        data={data.results[currentUser][selectedImage].insights}
        options={{
          title: "Line (dense time series)",
          axes: {
            bottom: {
              title: "Timestamp",
              mapsTo: "timestamp",
              scaleType: ScaleTypes.TIME,
            },
            left: {
              mapsTo: "value",
              title: "Value",
              scaleType: ScaleTypes.LINEAR,
            },
          },
          curve: "curveMonotoneX",
          height: "400px",
          theme: ChartTheme.G100,
        }}
      />
      <div>
        {data.results[currentUser][selectedImage].media.caption}
      </div>
    </div>
  {/if} 
</div>
