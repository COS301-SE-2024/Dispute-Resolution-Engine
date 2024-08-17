<script lang="ts">
	import FormField from '$lib/components/FormField.svelte';
	import { signupSchema } from '$lib/schema/auth.js';
	import { RadioGroup, RadioItem } from '@skeletonlabs/skeleton';
	import { superForm } from 'sveltekit-superforms';
	import { zod } from 'sveltekit-superforms/adapters';

	const steps: {
		id: string;
		name: string;
		fields: string[];
	}[] = [
		{ id: 'Step 1', name: 'User Type', fields: ['userType'] },
		{ id: 'Step 2', name: 'The Basics', fields: ['email', 'password', 'passwordConfirm'] },
		{
			id: 'Step 3',
			name: 'Personal Details',
			fields: ['firstName', 'lastName', 'gender', 'nationality', 'preferredLanguage', 'dateOfBirth']
		}
	];

	export let data;
	const { form, errors, enhance } = superForm(data.form, {
		validators: zod(signupSchema),
		validationMethod: 'onblur'
	});

	let currentStep = 0;

	$: lastStep = currentStep == steps.length - 1;
</script>

<div class="space-y-6">
	<header>
		<ol class="grid grid-cols-3 mb-6">
			{#each steps as step, i}
				<li class="step" data-state={i <= currentStep ? 'active' : 'inactive'}>
					<button
						class="btn text-left flex-col items-start"
						on:click={() => {
							currentStep = i;
						}}
					>
						<strong>{step.id}</strong>
						<span>{step.name}</span>
					</button>
				</li>
			{/each}
		</ol>
		<h2 class="h3">{steps[currentStep].id}</h2>
		<p>{steps[currentStep].name}</p>
	</header>

	<form class="flex flex-col gap-4 grow" id="form" method="post" use:enhance>
		{#if currentStep == 0}
			<RadioGroup flexDirection="flex-col" active="variant-outline-primary">
				<RadioItem bind:group={$form.userType} name="user_type" value="user">User</RadioItem>
				<RadioItem bind:group={$form.userType} name="user_type" value="expert">Expert</RadioItem>
			</RadioGroup>
		{:else if currentStep == 1}
			<FormField target="first_name" label="First Name" error={$errors.firstName}>
				<input
					id="first_name"
					type="text"
					class="input {$errors.firstName ? 'input-error' : ''}"
					placeholder="First Name"
					bind:value={$form.firstName}
					aria-invalid={$errors.firstName ? 'true' : undefined}
				/>
			</FormField>

			<FormField target="last_name" label="Last Name" error={$errors.lastName}>
				<input
					id="last_name"
					type="text"
					class="input"
					placeholder="Last Name"
					bind:value={$form.lastName}
					aria-invalid={$errors.lastName ? 'true' : undefined}
				/>
			</FormField>
			<FormField target="email" label="Email" error={$errors.email}>
				<input id="email" type="email" class="input" placeholder="Email" bind:value={$form.email} />
			</FormField>
			<FormField target="password" label="Password" error={$errors.password}>
				<input
					id="password"
					type="password"
					class="input"
					placeholder="Password"
					bind:value={$form.password}
				/>
			</FormField>
			<FormField target="password_confirm" label="Confirm Password" error={$errors.passwordConfirm}>
				<input
					id="password_confirm"
					type="password"
					class="input"
					placeholder="Confirm Password"
					bind:value={$form.passwordConfirm}
				/>
			</FormField>
		{:else if currentStep == 2}
			<FormField target="gender" label="Gender" error={$errors.gender}>
				<select name="gender" id="gender" class="select" bind:value={$form.gender}>
					<option value="Male">Male</option>
					<option value="Female">Female</option>
					<option value="Non-binary">Non-binary</option>
					<option value="Prefer not to say">Prefer not to say</option>
					<option value="Other">Other</option>
				</select>
			</FormField>
			<FormField target="lang" label="Preferred Language" error={$errors.preferredLanguage}>
				<select name="lang" id="lang" class="select" bind:value={$form.preferredLanguage}>
					<option value="Male">Male</option>
					<option value="Female">Female</option>
					<option value="Non-binary">Non-binary</option>
					<option value="Prefer not to say">Prefer not to say</option>
					<option value="Other">Other</option>
				</select>
			</FormField>
			<FormField target="nationality" label="Nationality" error={$errors.nationality}>
				<select name="nationality" id="nationality" class="select" bind:value={$form.nationality}>
					<option value="Male">Male</option>
					<option value="Female">Female</option>
					<option value="Non-binary">Non-binary</option>
					<option value="Prefer not to say">Prefer not to say</option>
					<option value="Other">Other</option>
				</select>
			</FormField>
			<FormField target="date_of_birth" label="Date of Birth" error={$errors.passwordConfirm}>
				<input id="date_of_birth" type="date" class="input" bind:value={$form.dateOfBirth} />
			</FormField>
		{/if}
	</form>

	<footer class="flex">
		<p class="grow">
			Already have an account?
			<a href="/login" class="anchor"> Login </a>
		</p>
		{#if lastStep}
			<button class="btn variant-filled-primary" form="form" type="submit">Sign Up</button>
		{:else}
			<button class="btn variant-soft-primary" type="button" on:click={() => currentStep++}>
				Next
			</button>
		{/if}
	</footer>
</div>

<style lang="postcss">
	.step[data-state='active'] {
		@apply border-t-4 border-secondary-500;
	}
</style>
