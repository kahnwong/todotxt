export const TODO_STATUS = {
  BACKLOG: 'backlog',
  IN_PROGRESS: 'in-progress',
  STUCK: 'stuck',
  DONE: 'done',
} as const

export const LANE_IDS = {
  TODAY: 'today',
} as const

export const CONTEXT = {
  TINKERING: '@tinkering',
  WORK: '@work',
} as const

export const UPDATE_INTERVAL_MS = 10000

export const FILE_EXTENSIONS = [
  'txt',
  'pdf',
  'doc',
  'docx',
  'xls',
  'xlsx',
  'jpg',
  'png',
  'gif',
  'zip',
  'tar',
  'gz',
] as const

export const COLUMN_CONFIG = [
  {
    id: TODO_STATUS.BACKLOG,
    title: 'Backlog',
    icon: 'mdi:inbox-multiple',
    iconClass: 'column-icon-backlog',
  },
  {
    id: TODO_STATUS.IN_PROGRESS,
    title: 'In Progress',
    icon: 'mdi:progress-clock',
    iconClass: 'column-icon-progress',
  },
  {
    id: TODO_STATUS.STUCK,
    title: 'Stuck',
    icon: 'mdi:alert-circle',
    iconClass: 'column-icon-stuck',
  },
  {
    id: TODO_STATUS.DONE,
    title: 'Done',
    icon: 'mdi:check-circle',
    iconClass: 'column-icon-done',
  },
] as const
