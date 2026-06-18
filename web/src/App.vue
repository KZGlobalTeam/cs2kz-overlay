<template>
  <div class="m-4 rounded-md text-2xl">
    <div v-if="mapName !== '<empty>' && mapName !== ''" class="flex gap-6">
      <div class="flex items-center gap-1.5">
        <div class="flex items-center gap-1">
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
        <span v-if="courseTiers" class="text-gray-600">/</span>
        <template v-if="courseTiers && courseTiers.nub !== courseTiers.pro">
          <div class="flex items-start gap-1">
            <span
              :style="{
                color: tierColorMap[courseTiers.nub],
              }"
            >
              T{{ tierNumberMap[courseTiers.nub] }}
            </span>
            <div
              class="flex justify-center items-center bg-yellow-600 text-gray-100 text-sm rounded-sm px-1"
            >
              TP
            </div>
          </div>
          <span class="text-gray-600">-</span>
          <div class="flex items-start gap-1">
            <span
              :style="{
                color: tierColorMap[courseTiers.pro],
              }"
            >
              T{{ tierNumberMap[courseTiers.pro] }}
            </span>
            <div
              class="flex justify-center items-center bg-blue-600 text-gray-100 text-sm rounded-sm px-1"
            >
              PRO
            </div>
          </div>
        </template>
        <span
          v-else-if="courseTiers"
          :style="{
            color: tierColorMap[courseTiers.nub],
          }"
        >
          T{{ tierNumberMap[courseTiers.nub] }}
        </span>
      </div>

      <div v-if="mode" class="flex items-center gap-1">
        <span class="text-gray-400">Mode:</span>
        <span class="text-gray-100 font-medium">{{ modeName }}</span>
      </div>
    </div>

    <RecordRow
      v-if="mapStatus === 'GLOBAL' && courseName"
      type="all"
      :wr="overallWr"
      :pb="overallPb"
      :player-profile="playerProfile"
      :wr-holder-profile="overallWrHolderProfile"
    />

    <RecordRow
      v-if="mapStatus === 'GLOBAL' && courseName"
      type="pro"
      :wr="proWr"
      :pb="proPb"
      :player-profile="playerProfile"
      :wr-holder-profile="proWrHolderProfile"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { Map, PlayerProfile, Record, RecordResponse } from './types'
import { api, tierNumberMap, tierColorMap } from './utils'
import RecordRow from './components/RecordRow.vue'

const mapName = ref<string>('')
const map = ref<Map | null>(null)
const courseName = ref<string>('')
const modeName = ref<string>('')
const playerId = ref<string>('')
const timerStatus = ref<string>('')
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
const courseTiers = computed(() => {
  if (map.value && mode.value && courseName.value) {
    const courseIndex = map.value.courses.findIndex((course) => course.name === courseName.value)!
    const filter = map.value.courses[courseIndex]!.filters[mode.value]
    return { nub: filter.nub_tier, pro: filter.pro_tier }
    return
  } else {
    return null
  }
})

watch(
  () => playerId.value,
  async (playerId) => {
    playerProfile.value = null
    if (playerId) {
      playerProfile.value = await getProfile(playerId)
    }
  },
  { immediate: true },
)

watch(
  () => timerStatus.value,
  (timerStatus) => {
    if (timerStatus === 'finished') {
      overallWr.value = null
      proWr.value = null
      overallPb.value = null
      proPb.value = null
      getWr()
      getPb()
    }
  },
)

watch(
  () => mapName.value,
  (mapName) => {
    map.value = null
    if (mapName !== '<empty>' && mapName !== '') {
      getMapInfo()
    }
  },
  { immediate: true },
)

watch(
  [() => mapName.value, () => courseName.value, () => mode.value],
  ([mapName, courseName, mode]) => {
    overallWr.value = null
    proWr.value = null
    if (mapName !== '<empty>' && mapName !== '' && courseName && mode) {
      getWr()
    }
  },
  { immediate: true },
)

watch(
  [() => mapName.value, () => courseName.value, () => mode.value, () => playerId.value],
  ([mapName, courseName, mode, playerId]) => {
    overallPb.value = null
    proPb.value = null
    if (mapName !== '<empty>' && mapName !== '' && courseName && mode && playerId) {
      getPb()
    }
  },
  { immediate: true },
)

watch(
  () => overallWr.value,
  async (overallWr) => {
    if (overallWr) {
      overallWrHolderProfile.value = await getProfile(overallWr.player.id)
    }
  },
)

watch(
  () => proWr.value,
  async (proWr) => {
    if (proWr) {
      proWrHolderProfile.value = await getProfile(proWr.player.id)
    }
  },
)

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
  // mapName.value = 'kz_antimony'
  // courseName.value = 'Main'
  // modeName.value = 'CKZ'
  // playerId.value = '76561199067702427'
  const eventSource = new EventSource('http://127.0.0.1:4433/events')

  eventSource.onmessage = (event: MessageEvent<string>) => {
    const data = JSON.parse(event.data)
    mapName.value = data.map
    courseName.value = data.course
    modeName.value = data.mode
    playerId.value = data.steamID
    timerStatus.value = data.timer
  }
}
</script>

<style scoped></style>
