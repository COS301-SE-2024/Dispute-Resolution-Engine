<script lang="ts">
	import { RadioGroup, RadioItem } from '@skeletonlabs/skeleton';

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

	let currentStep = 0;
	let type = '';
</script>

<ol>
	{#each steps as step, i}
		<li>
			<button
				class="btn"
				on:click={() => {
					currentStep = i;
				}}>{step.name}</button
			>
		</li>
	{/each}
</ol>

<header class="mb-4">
	<h2 class="h3">{steps[currentStep].id}</h2>
	<p>{steps[currentStep].name}</p>
</header>

<form class="flex flex-col gap-4">
	{#if currentStep == 0}
		<RadioGroup>
			<RadioItem bind:group={type} name="user_type" value="user">User</RadioItem>
			<RadioItem bind:group={type} name="user_type" value="expert">Expert</RadioItem>
		</RadioGroup>
	{:else if currentStep == 1}
		<div>
			<label for="first_name" class="label">First Name</label>
			<input id="first_name" type="text" class="input" placeholder="First Name" />
		</div>

		<div>
			<label for="last_name" class="label">Last Name</label>
			<input id="last_name" type="text" class="input" placeholder="Last Nmae" />
		</div>
		<div>
			<label for="email" class="label">Email</label>
			<input id="email" type="email" class="input" placeholder="Email" />
		</div>
		<div>
			<label for="password" class="label">Password</label>
			<input id="password" type="password" class="input" placeholder="Password" />
		</div>
		<div>
			<label for="password_confirm" class="label">Confirm Password</label>
			<input id="password_confirm" type="password" class="input" placeholder="Confirm Password" />
		</div>
	{:else if currentStep == 2}
		<div>
			<label for="gender" class="label">Gender</label>
			<select name="gender" id="gender" class="select">
				<option value="Male">Male</option>
				<option value="Female">Female</option>
				<option value="Non-binary">Non-binary</option>
				<option value="Prefer not to say">Prefer not to say</option>
				<option value="Other">Other</option>
			</select>
		</div>
		<div>
			<label for="lang" class="label">Preferred Language</label>
			<select name="lang" id="lang" class="select">
				<option value="Male">Male</option>
				<option value="Female">Female</option>
				<option value="Non-binary">Non-binary</option>
				<option value="Prefer not to say">Prefer not to say</option>
				<option value="Other">Other</option>
			</select>
		</div>
		<div>
			<label for="nationality" class="label">Nationality</label>
			<select name="nationality" id="nationality" class="select">
				<option value="Male">Male</option>
				<option value="Female">Female</option>
				<option value="Non-binary">Non-binary</option>
				<option value="Prefer not to say">Prefer not to say</option>
				<option value="Other">Other</option>
			</select>
		</div>
		<div>
			<label for="date_of_birth" class="label">Date of Birth</label>
			<input id="date_of_birth" type="date" class="input" />
		</div>
	{/if}
</form>
