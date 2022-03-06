<script>
	import App from "./app";
	import Canvas from "./Canvas.svelte";

	// parse parameters
	// parse requested file
	const urlSearchParams = new URLSearchParams(window.location.search);
	const params = Object.fromEntries(urlSearchParams.entries());

	// fire init events
	if (params.i != null) App.ImageUrl(params.i);
	if (params.hash != null) App.ImageHash(params.hash);
</script>

<main>
	<Canvas />
</main>

<style>
	main {
		text-align: center;
		padding: 0;
		margin: 0 auto;
	}
</style>
