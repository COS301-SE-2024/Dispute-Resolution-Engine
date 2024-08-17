<script lang="ts">
	/** The ID of the HTML element to bind the label to */
	export let target: string;

	/** The text content of the form field's label */
	export let label: string;

	/** The errors (if any) related to the field*/
	export let error: string[] | undefined;
</script>

<div>
	<label for={target} class="label {error ? 'text-error-500 font-semibold' : ''}">{label}</label>
	<slot />
	{#if error}
		<p role="alert" class="text-error-500">
			{error[0]}
		</p>
	{/if}
</div>
