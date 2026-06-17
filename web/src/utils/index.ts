import type { Tier } from '@/types'
import axios from 'axios'

const api = axios.create({
  baseURL: 'https://api.cs2kz.org',
})

const tierNumberMap = {
  'very-easy': 1,
  easy: 2,
  medium: 3,
  advanced: 4,
  hard: 5,
  'very-hard': 6,
  extreme: 7,
  death: 8,
  unfeasible: 9,
  impossible: 10,
}

const tierColorMap = {
  'very-easy': '#6bc96f',
  easy: '#33bd3a',
  medium: '#d8e302',
  advanced: '#FFC107',
  hard: '#e37910',
  'very-hard': '#e34202',
  extreme: '#e31c02',
  death: '#bb02db',
  unfeasible: '#e800e1',
  impossible: '#d1d1d1',
}

function formatTime(seconds: number) {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const remainingSeconds = (seconds % 60).toFixed(3)

  const timeParts = []

  if (hours > 0) {
    timeParts.push(hours.toString().padStart(2, '0'))
  }

  if (minutes > 0 || hours > 0) {
    timeParts.push(minutes.toString().padStart(2, '0'))
  }

  timeParts.push(remainingSeconds.padStart(6, '0'))

  return timeParts.join(':')
}

export { api, tierNumberMap, tierColorMap, formatTime }
