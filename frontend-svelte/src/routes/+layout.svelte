<script lang="ts">
	import '../app.postcss';
	import { AppBar } from '@skeletonlabs/skeleton';
</script>

<!-- App Shell -->
<AppBar gap="gap-4">
	<svelte:fragment slot="lead">
		<a href="/" class="flex items-center gap-2">
			<figure class="h-16">
				<section class="img-bg" />
				<img src="ralph.png" alt="DRE Mascot" class="aspect-square h-16" />
			</figure>
			<strong title="Dispute Resolution Engine" class="text-xl uppercase">DRE</strong>
		</a>
	</svelte:fragment>
	<nav>
		<a class="btn variant-soft-surface" href="/archive"> Archive </a>
		<a
			class="btn variant-soft-surface"
			href="https://cos301-se-2024.github.io/Dispute-Resolution-Engine/"
		>
			Help
		</a>
	</nav>
	<svelte:fragment slot="trail">
		<a class="btn btn-sm variant-ghost-surface" href="/profile" target="_blank" rel="noreferrer">
			Profile
		</a>
	</svelte:fragment>
</AppBar>
<slot />

<style lang="postcss">
	figure {
		@apply relative flex items-center justify-center;
	}
	figure img {
		@apply z-10;
	}
	.img-bg {
		@apply absolute bg-secondary-400 blur-[15px];
		@apply w-8 h-8;
	}
</style>
