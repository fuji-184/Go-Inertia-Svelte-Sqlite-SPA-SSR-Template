<script>
export let tes;
console.log(tes)
</script>

<div>home</div>

{#each tes as t}
  <div>{t.id}. {t.name}</div>
{/each}

<a href="/about">about</a>