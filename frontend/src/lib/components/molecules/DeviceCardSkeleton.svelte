<script lang="ts">
  import { Skeleton } from "$lib/components/ui/skeleton";
</script>

<div class="relative bg-card p-5 rounded-sm border border-transparent">
  <div class="flex justify-between items-start mb-6">
    <div class="flex items-center gap-2">
        <Skeleton class="w-1.5 h-1.5 rounded-full" />
        <Skeleton class="h-3 w-16" />
    </div>
    <Skeleton class="h-6 w-6 rounded-md" />
  </div>

  <div class="space-y-3 mb-6">
    <div>
        <Skeleton class="h-5 w-32 mb-1" />
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
