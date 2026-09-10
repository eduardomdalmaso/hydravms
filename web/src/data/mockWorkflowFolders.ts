import type { WorkflowFolderNode, EnterpriseWorkflowItem, PaletteBlockDefinition } from '../types/workflowTree'

export const paletteBlocks: PaletteBlockDefinition[] = [
  { kind: 'AI_DETECT', type: 'TRIGGER', title: 'DETECCAO IA', summary: 'Pessoas, veiculos ou EPI', defaultConfig: { targetClass: 'PESSOA', minConfidence: 85 } },
  { kind: 'ALARM_SENSOR', type: 'TRIGGER', title: 'SENSOR / ALARME', summary: 'Barreira perimetral', defaultConfig: { zone: 'ZONA_01', eventType: 'VIOLACAO' } },
  { kind: 'BRANCH_OR', type: 'CONDITION', title: 'JUNCAO EM Y (OU)', summary: 'Interliga 2 entradas em 1 saida', defaultConfig: { logic: 'OR' } },
  { kind: 'LOGIC_AND', type: 'CONDITION', title: 'PORTA E (AND)', summary: 'Exige 2 eventos simultaneos', defaultConfig: { windowSeconds: 5 } },
  { kind: 'SPLIT_Y', type: 'CONDITION', title: 'DIVISAO EM Y', summary: '1 evento bifurca para 2 saidas', defaultConfig: { mode: 'PARALELO' } },
  { kind: 'SCHEDULE', type: 'CONDITION', title: 'AGENDAMENTO', summary: 'Janela de horario', defaultConfig: { timeWindow: 'NOTURNO // 22H AS 06H' } },
  { kind: 'ANTI_SPAM', type: 'CONDITION', title: 'ANTI-SPAM', summary: 'Silenciamento continuo', defaultConfig: { cooldownSeconds: 30 } },
  { kind: 'TELEGRAM', type: 'ACTION', title: 'TELEGRAM BOT', summary: 'Texto e snapshot', defaultConfig: { chatId: '-1001928374', sendSnapshot: true } },
  { kind: 'WEBHOOK', type: 'ACTION', title: 'WEBHOOK HTTP', summary: 'POST para SIEM externo', defaultConfig: { endpointUrl: 'https://siem.corp.com/alerts' } },
  { kind: 'RELAY_SIREN', type: 'ACTION', title: 'ACIONAR SIRENE', summary: 'Ativar rele fisico', defaultConfig: { relayId: 'RELE_01', durationSec: 15 } },
  { kind: 'TERMINAL_STOP', type: 'ACTION', title: 'PARE / FIM', summary: 'Finalizar execucao do pipeline', defaultConfig: {} }
]

export const initialRootWorkflows: EnterpriseWorkflowItem[] = []
export const initialWorkflowFolders: WorkflowFolderNode[] = []
