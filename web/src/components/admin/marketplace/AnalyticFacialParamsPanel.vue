<script setup lang="ts">
defineProps<{
  minHeight: number; maxHeight: number | string; minFaces: number; blurThreshold: number
  recognitionThreshold: number; detectorSensitivity: 'baixo' | 'normal' | 'alto'; detectGenderAge: boolean
}>()
const emit = defineEmits<{
  (e: 'update:minHeight', val: number): void; (e: 'update:maxHeight', val: number | string): void
  (e: 'update:minFaces', val: number): void; (e: 'update:blurThreshold', val: number): void
  (e: 'update:recognitionThreshold', val: number): void; (e: 'update:detectorSensitivity', val: 'baixo' | 'normal' | 'alto'): void
  (e: 'update:detectGenderAge', val: boolean): void
}>()
</script>

<template>
  <div class="vms-flex-col" style="gap: 0.45rem;">
    <!-- Sensibilidade do Detector (3 Níveis Estilo Vezha) -->
    <div class="vms-flex-col" style="gap: 0.3rem; background: rgba(0,0,0,0.25); padding: 0.45rem; border-radius: 6px; border: 1px solid var(--vms-border);">
      <div class="vms-flex-between" style="align-items: center;">
        <span class="vms-text-mono vms-text-2xs vms-text-dim">SENSIBILIDADE DO DETECTOR</span>
        <span class="vms-info-tip" title="Valor mais baixo: mais detecções e maior risco de falsos alarmes. Valor mais alto: menos detecções e menor risco de falsos alarmes.">?</span>
      </div>
      <div class="vms-flex-row" style="gap: 0.35rem;">
        <button v-for="lvl in (['baixo', 'normal', 'alto'] as const)" :key="lvl" type="button" class="vms-btn vms-btn-sm" :class="detectorSensitivity === lvl ? 'vms-btn-primary' : 'vms-btn-secondary'" style="flex: 1; padding: 3px; font-size: 9.5px;" @click="emit('update:detectorSensitivity', lvl)">
          {{ lvl.toUpperCase() }}
        </button>
      </div>
    </div>

    <!-- Filtros de Altura da Face -->
    <div class="vms-flex-col" style="gap: 0.3rem; background: rgba(0,0,0,0.25); padding: 0.45rem; border-radius: 6px; border: 1px solid var(--vms-border);">
      <div class="vms-flex-between" style="align-items: center;">
        <span class="vms-text-mono vms-text-2xs vms-text-dim">DIMENSÃO DA FACE (BOUNDING BOX)</span>
        <span class="vms-info-tip" title="Altura mínima filtra ruídos distantes. Altura máxima em branco (ou 0) aceita qualquer aproximação sem impor teto.">?</span>
      </div>
      <div class="vms-flex-row" style="gap: 0.45rem;">
        <div class="vms-flex-col" style="flex: 1; gap: 2px;">
          <label class="vms-text-mono vms-text-2xs" style="color: #fff;">MÍNIMA (PX):</label>
          <input :value="minHeight" type="number" min="10" max="1080" class="vms-auth-input" style="padding: 3px 5px; font-size: 11px;" @input="emit('update:minHeight', Number(($event.target as HTMLInputElement).value))" />
        </div>
        <div class="vms-flex-col" style="flex: 1; gap: 2px;">
          <label class="vms-text-mono vms-text-2xs" style="color: #fff;">MÁXIMA (PX):</label>
          <input :value="maxHeight === 0 ? '' : maxHeight" type="number" min="0" max="2160" placeholder="Sem limite (0)" class="vms-auth-input" style="padding: 3px 5px; font-size: 11px;" @input="emit('update:maxHeight', ($event.target as HTMLInputElement).value ? Number(($event.target as HTMLInputElement).value) : 0)" />
        </div>
      </div>
    </div>

    <!-- Critérios de Rastreamento e Nitidez -->
    <div class="vms-flex-col" style="gap: 0.3rem; background: rgba(0,0,0,0.25); padding: 0.45rem; border-radius: 6px; border: 1px solid var(--vms-border);">
      <div class="vms-flex-between" style="align-items: center;">
        <span class="vms-text-mono vms-text-2xs vms-text-dim">CRITÉRIOS DE CAPTURA & TRACKING</span>
        <span class="vms-info-tip" title="25 frames equivalem a ~1s a 25 FPS (ou ~1.6s a 15 FPS). O sistema processa no teto ou na saída do alvo.">?</span>
      </div>
      <div class="vms-flex-row" style="gap: 0.45rem;">
        <div class="vms-flex-col" style="flex: 1; gap: 2px;">
          <label class="vms-text-mono vms-text-2xs" style="color: #fff;">MÍNIMO DE FACES:</label>
          <input :value="minFaces" type="number" min="1" max="120" class="vms-auth-input" style="padding: 3px 5px; font-size: 11px;" @input="emit('update:minFaces', Number(($event.target as HTMLInputElement).value))" />
        </div>
        <div class="vms-flex-col" style="flex: 1; gap: 2px;">
          <label class="vms-text-mono vms-text-2xs" style="color: #fff;">FILTRO ANTI-BLUR:</label>
          <input :value="blurThreshold" type="number" min="10" max="300" class="vms-auth-input" style="padding: 3px 5px; font-size: 11px;" @input="emit('update:blurThreshold', Number(($event.target as HTMLInputElement).value))" />
        </div>
      </div>
    </div>

    <!-- Limiar de Precisão Biométrica (3 Níveis) -->
    <div class="vms-flex-col" style="gap: 0.3rem; background: rgba(0,0,0,0.25); padding: 0.45rem; border-radius: 6px; border: 1px solid var(--vms-border);">
      <div class="vms-flex-between" style="align-items: center;">
        <span class="vms-text-mono vms-text-2xs" style="color: #fff;">LIMIAR DE PRECISÃO BIOMÉTRICA:</span>
        <span class="vms-info-tip" title="Baixo (65%): Monitoramento geral e alta tolerância. Médio (75%): Padrão de fábrica equilibrado. Alto (85%): Alta segurança e controle de acesso restrito.">?</span>
      </div>
      <div class="vms-flex-row" style="gap: 0.35rem;">
        <button v-for="p in ([{ val: 65, label: 'BAIXO (65%)' }, { val: 75, label: 'MÉDIO (75%)' }, { val: 85, label: 'ALTO (85%)' }] as const)" :key="p.val" type="button" class="vms-btn vms-btn-sm" :class="recognitionThreshold === p.val ? 'vms-btn-primary' : 'vms-btn-secondary'" style="flex: 1; padding: 3px; font-size: 9.5px;" @click="emit('update:recognitionThreshold', p.val)">
          {{ p.label }}
        </button>
      </div>
    </div>

    <!-- Detecção de Idade e Gênero (Toggle Estilo Vezha) -->
    <div class="vms-flex-between" style="align-items: center; background: rgba(0,0,0,0.25); padding: 0.45rem 0.55rem; border-radius: 6px; border: 1px solid var(--vms-border);">
      <label class="vms-flex-row vms-text-xs vms-font-semibold" style="gap: 0.45rem; align-items: center; cursor: pointer; color: #fff;">
        <input :checked="detectGenderAge" type="checkbox" style="accent-color: var(--vms-neu-accent-orange);" @change="emit('update:detectGenderAge', ($event.target as HTMLInputElement).checked)" />
        <span>Detecção de Idade e Gênero</span>
      </label>
      <span class="vms-info-tip" title="Classifica faixa etária e gênero no evento de biometria.">?</span>
    </div>
  </div>
</template>

<style scoped>
.vms-info-tip { display: inline-flex; align-items: center; justify-content: center; width: 14px; height: 14px; border-radius: 50%; background: rgba(255, 255, 255, 0.08); color: #8b94a0; font-size: 10px; font-family: var(--vms-font-jetbrains); font-weight: 700; cursor: help; transition: all 0.2s ease; user-select: none; }
.vms-info-tip:hover { background: var(--vms-neu-accent-orange); color: #ffffff; transform: scale(1.1); }
</style>
