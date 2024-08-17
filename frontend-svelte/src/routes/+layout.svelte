<script lang="ts">
	import '../app.postcss';
	import { LightSwitch } from '@skeletonlabs/skeleton';
</script>

<!-- App Shell -->

<header
	class="grid grid-cols-[1fr_auto] md:grid-cols-[auto_1fr_auto] gap-6 sticky top-0 z-10 bg-surface-100 dark:bg-surface-800 p-4"
>
	<div>
		<a href="/" class="flex items-center gap-2">
			<figure class="h-16">
				<section class="img-bg" />
				<img src="ralph.png" alt="DRE Mascot" class="aspect-square h-16" />
			</figure>
			<strong title="Dispute Resolution Engine" class="text-xl uppercase">DRE</strong>
		</a>
	</div>
	<nav class="row-start-2 md:row-start-auto flex items-center gap-2">
		<a class="btn variant-soft-surface" href="/archive"> Archive </a>
		<a
			class="btn variant-soft-surface"
			href="https://cos301-se-2024.github.io/Dispute-Resolution-Engine/"
		>
			Help
		</a>
	</nav>
	<div class="flex items-center gap-2">
		<LightSwitch />
		<a class="btn btn-sm variant-ghost-surface" href="/profile" target="_blank" rel="noreferrer">
			Profile
		</a>
	</div>
</header>

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
