import { FILE_EXTENSIONS } from '../constants/kanban'

/**
 * Detects URLs in text and prepends https:// if missing
 */
export function prependHttpsToUrls(text: string): string {
  // Match URLs without protocol: domain.tld or subdomain.domain.tld
  const urlPattern = /\b([a-zA-Z0-9][-a-zA-Z0-9]*\.)+[a-zA-Z]{2,}\b(\/[^\s]*)?/g

  return text.replace(urlPattern, (match) => {
    // Check if already has a protocol
    const hasProtocol = /^(https?|ftp|ftps):\/\//i.test(match)
    if (hasProtocol) {
      return match
    }

    // Check if this looks like a file extension
    const parts = match.split('/')
    const domain = parts[0]
    if (!domain) {
      return match
    }

    const lastDot = domain.lastIndexOf('.')
    if (lastDot !== -1) {
      const tld = domain.substring(lastDot + 1).toLowerCase()
      if ((FILE_EXTENSIONS as readonly string[]).includes(tld) && parts.length === 1) {
        return match
      }
    }

    return `https://${match}`
  })
}

/**
 * Highlights syntax in todo text
 */
export function highlightTodoSyntax(text: string): string {
  // Escape HTML
  let highlighted = text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

  // Highlight creation date (YYYY-MM-DD at start, not prefixed with due:)
  highlighted = highlighted.replace(
    /^(\d{4}-\d{2}-\d{2})\b/g,
    '<span class="highlight-date">$1</span>',
  )

  // Highlight @context tags
  highlighted = highlighted.replace(/(@\w+)/g, '<span class="highlight-context">$1</span>')

  // Highlight +project tags
  highlighted = highlighted.replace(/(\+\w+)/g, '<span class="highlight-project">$1</span>')

  // Highlight =status tags
  highlighted = highlighted.replace(/(=[\w-]+)/g, '<span class="highlight-status">$1</span>')

  // Highlight due: tags
  highlighted = highlighted.replace(
    /(due:\d{4}-\d{2}-\d{2})/g,
    '<span class="highlight-due">$1</span>',
  )

  return highlighted
}

/**
 * Reconstructs full todo text from Todo object
 */
export function reconstructTodoText(todo: {
  created_date?: string
  project?: string
  todo: string
  context?: string
  status?: string
}): string {
  const parts = []
  if (todo.created_date) parts.push(todo.created_date)
  if (todo.project) parts.push(todo.project)
  parts.push(todo.todo)
  if (todo.context) parts.push(todo.context)
  // Don't append =backlog since it's the default state
  if (todo.status && todo.status !== 'backlog') parts.push(`=${todo.status}`)
  return parts.join(' ')
}

/**
 * Cleans todo text by removing default status and extra spaces
 */
export function cleanTodoText(text: string): string {
  // Remove =backlog since it's the default state
  let cleaned = text.replace(/=backlog\b/g, '').trim()
  // Clean up multiple spaces
  cleaned = cleaned.replace(/\s+/g, ' ')
  return cleaned
}
