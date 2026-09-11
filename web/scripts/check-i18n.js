#!/usr/bin/env node
import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const localesDir = path.resolve(__dirname, '../src/locales')

function extractKeys(filePath) {
  let content = fs.readFileSync(filePath, 'utf-8')
  // Strip single-quoted and double-quoted strings
  content = content.replace(/'[^']*'/g, "''").replace(/"[^"]*"/g, '""')
  
  const keys = []
  const regex = /([a-zA-Z0-9_]+)\s*:\s*['"]/g
  let match
  while ((match = regex.exec(content)) !== null) {
    keys.push(match[1])
  }
  return new Set(keys)
}

const ptKeys = extractKeys(path.join(localesDir, 'pt.ts'))
const enKeys = extractKeys(path.join(localesDir, 'en.ts'))
const esKeys = extractKeys(path.join(localesDir, 'es.ts'))

let hasError = false

function checkMissing(sourceName, targetName, sourceKeys, targetKeys) {
  const missing = [...sourceKeys].filter(k => !targetKeys.has(k))
  if (missing.length > 0) {
    hasError = true
    console.error(`❌ [i18n Error] Chaves presentes em ${sourceName} mas ausentes em ${targetName}:`, missing)
  }
}

console.log(`🌐 Verificando paridade de traduções (${ptKeys.size} chaves PT, ${enKeys.size} chaves EN, ${esKeys.size} chaves ES)...`)

checkMissing('pt.ts', 'en.ts', ptKeys, enKeys)
checkMissing('pt.ts', 'es.ts', ptKeys, esKeys)
checkMissing('en.ts', 'pt.ts', enKeys, ptKeys)
checkMissing('en.ts', 'es.ts', enKeys, esKeys)
checkMissing('es.ts', 'pt.ts', esKeys, ptKeys)
checkMissing('es.ts', 'en.ts', esKeys, enKeys)

if (hasError) {
  console.error('\n🚨 VIOLAÇÃO I18N: Há discrepâncias entre os idiomas. Adicione as traduções ausentes!')
  process.exit(1)
} else {
  console.log('✅ 100% de paridade i18n confirmada (todas as chaves existem em PT, EN e ES)!')
}
