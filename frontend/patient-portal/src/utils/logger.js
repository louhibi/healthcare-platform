import chalk from 'chalk'
import { useDebugStore } from '../stores/debug'

// Lightweight logging wrapper that respects debug store flags.
// Adds caller file:line information by parsing a generated stack trace.
// Designed to work in modern browsers & Vite dev environment.

function getStoreSafe() {
  try { return useDebugStore() } catch { return null }
}

function shouldLog() {
  const store = getStoreSafe()
  return !!store && store.logsEnabled === true
}

function ts() { return new Date().toISOString() }

function getCallerFrame() {
  // Generate an Error to capture stack frames
  const err = new Error()
  if (!err.stack) return null
  const lines = err.stack.split('\n').map(l => l.trim())
  // Typical formats:
  // Chrome/V8: at functionName (http://localhost:3000/src/components/Foo.vue:10:15)
  // Firefox: functionName@http://localhost:3000/src/components/Foo.vue:10:15
  // We skip frames belonging to this logger file
  const skipPattern = /\/src\/utils\/logger\.js/
  for (const line of lines.slice(2, 12)) { // skip first lines (Error + current fn) limit search depth
    if (!line) continue
    if (skipPattern.test(line)) continue
    // Extract URL or path with line:col
    const match = line.match(/(\(?)(https?:\/\/[^)]+|\/[^)]+):(\d+):(\d+)/)
    if (match) {
      const fullPath = match[2]
      const lineNo = match[3]
      // Derive file name (strip query params)
      const pathPart = fullPath.split('?')[0]
      const fileName = pathPart.split('/').pop() || pathPart
      return { file: fileName, line: lineNo }
    }
  }
  return null
}

function formatPrefix(level, colorFn) {
  const lvl = colorFn(`[${level.toUpperCase()}]`)
  const time = chalk.gray(ts())
  const caller = getCallerFrame()
  const loc = caller ? chalk.magenta(`[${caller.file}:${caller.line}]`) : chalk.magenta('[unknown]')
  return `${lvl}${time ? ' ' + time : ''} ${loc}`
}

// Core log emitter
function emit(level, colorFn, args) {
  if (!shouldLog()) return
  const prefix = formatPrefix(level, colorFn)
  // eslint-disable-next-line no-console
  console.log(prefix, ...args)
}

export const logger = {
  debug: (...a) => emit('debug', chalk.cyan, a),
  info: (...a) => emit('info', chalk.green, a),
  warn: (...a) => emit('warn', chalk.yellow, a),
  error: (...a) => emit('error', chalk.red, a),
  group(label) {
    if (!shouldLog()) return { end: () => {} }
    const header = formatPrefix('group', chalk.magenta) + ' ' + label
    // eslint-disable-next-line no-console
    console.log(header)
    return {
      end: () => {
        const footer = formatPrefix('group-end', chalk.magenta) + ' ' + label
        // eslint-disable-next-line no-console
        console.log(footer)
      }
    }
  }
}

export default logger
