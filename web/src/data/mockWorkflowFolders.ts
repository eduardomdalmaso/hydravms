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

export const initialRootWorkflows: EnterpriseWorkflowItem[] = [
  {
    id: 'wf_root_01',
    name: 'Alerta Perimetro // Deteccao & Fim',
    companyScope: 'EMPRESA // ALPHA SEGURANCA',
    is_locked: false,
    status: 'ATIVO',
    cooldownSeconds: 30,
    nodes: [
      { id: 'n_01', type: 'TRIGGER', kind: 'AI_DETECT', title: 'DETECCAO IA', summary: 'Fluxo Camera', position: { x: 80, y: 140 }, config: { targetClass: 'TODAS' } },
      { id: 'n_02', type: 'ACTION', kind: 'TERMINAL_STOP', title: 'PARE / FIM', summary: 'Finalizar', position: { x: 360, y: 140 }, config: {} }
    ],
    edges: [],
    createdAt: '2026-03-01'
  }
]

export const initialWorkflowFolders: WorkflowFolderNode[] = [
  {
    id: 'wff_01',
    name: 'EMPRESA // ALPHA SEGURANCA',
    clientType: 'company',
    isExpanded: true,
    workflows: [
      {
        id: 'wf_01',
        name: 'Invasao de Docas // Webhook SIEM',
        companyScope: 'EMPRESA // ALPHA SEGURANCA',
        is_locked: false,
        status: 'ATIVO',
        cooldownSeconds: 20,
        nodes: [
          { id: 'n_05', type: 'TRIGGER', kind: 'ALARM_SENSOR', title: 'SENSOR / ALARME', summary: 'Zona: DOCAS_BARREIRA', position: { x: 80, y: 100 }, config: { zone: 'DOCAS_BARREIRA' } },
          { id: 'n_06', type: 'ACTION', kind: 'WEBHOOK', title: 'WEBHOOK HTTP', summary: 'POST /alerts', position: { x: 380, y: 100 }, config: { endpointUrl: 'https://siem.corp.com/alerts' } }
        ],
        edges: [{ id: 'e_05_06', source: 'n_05', target: 'n_06', animated: true }],
        createdAt: '2026-03-02'
      },
      {
        id: 'wf_02',
        name: 'Sirene Automatica // Barreira Norte',
        companyScope: 'EMPRESA // ALPHA SEGURANCA',
        is_locked: true,
        status: 'PAUSADO',
        cooldownSeconds: 60,
        nodes: [
          { id: 'n_07', type: 'TRIGGER', kind: 'ALARM_SENSOR', title: 'SENSOR / ALARME', summary: 'Zona: BARREIRA_NORTE', position: { x: 80, y: 100 }, config: { zone: 'BARREIRA_NORTE' } },
          { id: 'n_08', type: 'ACTION', kind: 'RELAY_SIREN', title: 'ACIONAR SIRENE', summary: 'Rele 01 // Duracao 15s', position: { x: 380, y: 100 }, config: { relayId: 'RELE_01', durationSec: 15 } }
        ],
        edges: [{ id: 'e_07_08', source: 'n_07', target: 'n_08', animated: true }],
        createdAt: '2026-03-03'
      }
    ]
  },
  {
    id: 'wff_02',
    name: 'CLIENTE // CONDOMINIO VISTA VERDE',
    clientType: 'final_client',
    isExpanded: false,
    workflows: [
      {
        id: 'wf_03',
        name: 'Pop-up Operador // Veiculo Suspeito',
        companyScope: 'CLIENTE // CONDOMINIO VISTA VERDE',
        is_locked: false,
        status: 'ATIVO',
        cooldownSeconds: 15,
        nodes: [
          { id: 'n_09', type: 'TRIGGER', kind: 'AI_DETECT', title: 'DETECCAO IA', summary: 'Classe: VEICULO // Portaria', position: { x: 80, y: 100 }, config: { targetClass: 'VEICULO' } },
          { id: 'n_10', type: 'ACTION', kind: 'WEBSOCKET', title: 'WEBSOCKET OPERADOR', summary: 'Pop-up HUD com som', position: { x: 380, y: 100 }, config: { priority: 'ALTA' } }
        ],
        edges: [{ id: 'e_09_10', source: 'n_09', target: 'n_10', animated: true }],
        createdAt: '2026-03-04'
      }
    ]
  }
]
