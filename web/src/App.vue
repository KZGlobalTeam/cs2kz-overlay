<template>
  <div class="m-4 bg-black rounded-md p-2 text-2xl">
    <div class="flex gap-6">
      <div class="flex items-center gap-1">
        <div v-if="mapName" class="flex items-center gap-1">
          <span class="text-gray-400">Map:</span>
          <span class="text-white font-medium">{{ mapName }}</span>
        </div>
        <span
          v-if="mapStatus === 'NON-GLOBAL'"
          class="px-1 text-sm rounded-sm text-yellow-300 bg-yellow-900/90"
        >
          NON-GLOBAL
        </span>
        <span
          v-if="mapStatus === 'GLOBAL'"
          class="px-1 border text-sm rounded-sm text-green-300 bg-green-900/90"
        >
          GLOBAL
        </span>
      </div>

      <div class="flex gap-1">
        <div v-if="courseName" class="flex items-center gap-1">
          <span class="text-gray-400">Course:</span>
          <span class="text-gray-100 font-medium">{{ courseName }}</span>
        </div>
        <span class="text-gray-400">/</span>
        <span
          v-if="courseTier"
          :style="{
            color: tierColorMap[courseTier],
          }"
        >
          T{{ tierNumberMap[courseTier] }}
        </span>
      </div>

      <div v-if="mode" class="flex items-center gap-1">
        <span class="text-gray-400">Mode:</span>
        <span class="text-gray-100 font-medium">{{ modeName }}</span>
      </div>
    </div>

    <RecordRow
      v-if="mapStatus === 'GLOBAL'"
      type="overall"
      :wr="overallWr"
      :pb="overallPb"
      :gain="overallGain"
      :player-profile="playerProfile"
      :wr-holder-profile="overallWrHolderProfile"
    />

    <RecordRow
      v-if="mapStatus === 'GLOBAL'"
      type="pro"
      :wr="proWr"
      :pb="proPb"
      :gain="proGain"
      :player-profile="playerProfile"
      :wr-holder-profile="proWrHolderProfile"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { GameState, Map, PlayerProfile, Record, RecordResponse } from './types'
import { api, tierNumberMap, tierColorMap } from './utils'
import RecordRow from './components/RecordRow.vue'

const mapName = ref<string | null>(null)
const map = ref<Map | null>(null)
const courseName = ref<string | null>(null)
const modeName = ref<string | null>(null)
const playerId = ref<string | null>(null)
const overallPb = ref<Record | null>(null)
const proPb = ref<Record | null>(null)
const overallWr = ref<Record | null>(null)
const proWr = ref<Record | null>(null)
const playerProfile = ref<PlayerProfile | null>(null)
const overallWrHolderProfile = ref<PlayerProfile | null>(null)
const proWrHolderProfile = ref<PlayerProfile | null>(null)

const mode = computed(() => {
  if (modeName.value) {
    return modeName.value === 'CKZ' ? 'classic' : 'vanilla'
  } else {
    return null
  }
})

const mapStatus = computed(() => {
  if (map.value) {
    return map.value.state === 'approved' ? 'GLOBAL' : 'NON-GLOBAL'
  } else {
    return 'NON-GLOBAL'
  }
})
const courseTier = computed(() => {
  if (map.value && mode.value && courseName.value) {
    const courseIndex = map.value.courses.findIndex((course) => course.name === courseName.value)!
    return map.value.courses[courseIndex]!.filters[mode.value].nub_tier
  } else {
    return null
  }
})

const overallGain = computed(() => {
  if (overallWr.value && overallPb.value) {
    return overallPb.value.time - overallWr.value.time
  } else {
    return null
  }
})

const proGain = computed(() => {
  if (proWr.value && proPb.value) {
    return proPb.value.time - proWr.value.time
  } else {
    return null
  }
})

watch(playerId, async (playerId) => {
  if (playerId) {
    playerProfile.value = await getProfile(playerId)
  }
})

watch(mapName, (mapName) => {
  if (mapName) {
    getMapInfo()
  }
})

watch(
  [() => mapName.value, () => courseName.value, () => mode.value],
  ([mapName, courseName, mode]) => {
    if (mapName && courseName && mode) {
      getWr()
    }
  },
)

watch(
  [() => mapName.value, () => courseName.value, () => mode.value, () => playerId.value],
  ([mapName, courseName, mode, playerId]) => {
    if (mapName && courseName && mode && playerId) {
      getPb()
    }
  },
)

watch(overallWr, async (overallWr) => {
  if (overallWr) {
    overallWrHolderProfile.value = await getProfile(overallWr.player.id)
  }
})

watch(proWr, async (proWr) => {
  if (proWr) {
    proWrHolderProfile.value = await getProfile(proWr.player.id)
  }
})

initListener()

async function getProfile(steamId: string) {
  try {
    const { data } = await api.get<PlayerProfile | undefined>(`/players/${steamId}/steam-profile`)
    return data || null
  } catch (error) {
    console.log('[fetch error]', error)
    return null
  }
}

async function getMapInfo() {
  try {
    const { data } = await api.get<Map | undefined>(`/maps/${mapName.value}`)
    if (data) {
      map.value = data
    } else {
      map.value = null
    }
  } catch (error) {
    console.log('[fetch error]', error)
    map.value = null
  }
}

async function getPb() {
  try {
    const [overallRes, proRes] = await Promise.all([
      api.get<RecordResponse | undefined>('/records', {
        params: {
          top: true,
          map: mapName.value,
          course: courseName.value,
          mode: mode.value,
          player: playerId.value,
        },
      }),
      api.get<RecordResponse | undefined>('/records', {
        params: {
          top: true,
          map: mapName.value,
          course: courseName.value,
          mode: mode.value,
          player: playerId.value,
          has_teleports: false,
        },
      }),
    ])

    if (overallRes.data && overallRes.data.values.length > 0) {
      overallPb.value = overallRes.data.values[0]!
    } else {
      overallPb.value = null
    }
    if (proRes.data && proRes.data.values.length > 0) {
      proPb.value = proRes.data.values[0]!
    } else {
      proPb.value = null
    }
  } catch (error) {
    console.log('[fetch error]', error)
    overallPb.value = null
    proPb.value = null
  }
}

async function getWr() {
  try {
    const [overallRes, proRes] = await Promise.all([
      api.get<RecordResponse | undefined>('/records', {
        params: {
          top: true,
          map: mapName.value,
          course: courseName.value,
          mode: mode.value,
          max_rank: 1,
        },
      }),
      api.get<RecordResponse | undefined>('/records', {
        params: {
          top: true,
          map: mapName.value,
          course: courseName.value,
          mode: mode.value,
          has_teleports: false,
          max_rank: 1,
        },
      }),
    ])

    if (overallRes.data && overallRes.data.values.length > 0) {
      overallWr.value = overallRes.data.values[0]!
    } else {
      overallWr.value = null
    }
    if (proRes.data && proRes.data.values.length > 0) {
      proWr.value = proRes.data.values[0]!
    } else {
      proWr.value = null
    }
  } catch (error) {
    console.log('[fetch error]', error)
    overallWr.value = null
    proWr.value = null
  }
}

function initListener() {
  mapName.value = 'kz_antimony'
  courseName.value = 'Main'
  modeName.value = 'CKZ'
  playerId.value = '76561199067702427'
  // const eventSource = new EventSource('http://127.0.0.1:4433/events')

  // eventSource.onmessage = (event: MessageEvent<GameState>) => {
  //   mapName.value = event.data.map
  //   courseName.value = event.data.course
  //   mode.value = event.data.mode
  //   playerId.value = event.data.steamId
  // }
}
</script>

<style scoped></style>
