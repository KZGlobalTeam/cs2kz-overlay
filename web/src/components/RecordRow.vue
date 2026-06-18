<template>
  <div class="flex mt-2 gap-2">
    <span :class="type === 'all' ? 'text-yellow-600' : 'text-blue-600'" class="w-12 text-right">{{
      type.toUpperCase()
    }}</span>
    <span class="text-gray-400">|</span>

    <div v-if="wr" class="flex gap-3">
      <div class="flex items-start gap-1">
        <span class="text-gray-100">{{ formatTime(wr.time) }}</span>
        <div
          class="flex justify-center items-center text-gray-100 text-sm rounded-sm px-1"
          :class="wr.teleports > 0 ? 'bg-yellow-600' : 'bg-blue-600'"
        >
          {{ wr.teleports > 0 ? 'TP' : 'PRO' }}
        </div>
      </div>
      <span class="text-gray-300">by</span>

      <div class="flex gap-1.5">
        <img
          v-if="wrHolderProfile"
          class="rounded-full w-8 h-auto ring-1 ring-gray-400"
          :src="wrHolderProfile.avatar_url"
          alt=""
        />
        <span :class="gain === 0 ? 'text-green-500' : 'text-cyan-500'">{{ wr.player.name }}</span>
      </div>
      <span class="text-amber-600" v-if="gain && gain > 0">{{ `(+${formatTime(gain)})` }}</span>
    </div>
    <div v-else class="text-gray-400">No records yet</div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { PlayerProfile, Record } from '@/types'
import { formatTime } from '@/utils'

const props = defineProps<{
  type: 'all' | 'pro'
  wr: Record | null
  pb: Record | null
  wrHolderProfile: PlayerProfile | null
  playerProfile: PlayerProfile | null
}>()

const gain = computed(() => {
  if (props.wr && props.pb) {
    return props.pb.time - props.wr.time
  } else {
    return null
  }
})
</script>
