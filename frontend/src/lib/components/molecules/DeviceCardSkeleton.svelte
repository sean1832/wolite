<script lang="ts">
	import { Skeleton } from '$lib/components/ui/skeleton';
</script>

<div class="relative rounded-sm border border-transparent bg-card p-5">
	<div class="mb-6 flex items-start justify-between">
		<div class="flex items-center gap-2">
			<Skeleton class="h-1.5 w-1.5 rounded-full" />
			<Skeleton class="h-3 w-16" />
		</div>
		<Skeleton class="h-6 w-6 rounded-md" />
	</div>

	<div class="mb-6 space-y-3">
		<div>
			<Skeleton class="mb-1 h-5 w-32" />
			<Skeleton class="h-3 w-48" />
		</div>
		<div class="flex flex-col gap-0.5">
			<Skeleton class="h-3 w-24" />
		</div>
	</div>

	<div class="flex items-center justify-end">
		<Skeleton class="h-8 w-16" />
	</div>
</div>
