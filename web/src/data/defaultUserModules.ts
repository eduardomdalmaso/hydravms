import type { UserModulePermission, UserRole } from '../types/userTree'

export function createDefaultUserModules(role: UserRole, scopeTarget: string = 'GLOBAL'): UserModulePermission[] {
  const isMaster = role === 'admin_master'
  const isCompany = isMaster || role === 'company_admin'
  const isOperator = isCompany || role === 'operator'

  return [
    {
      id: 'admin_center',
      name: 'ADMIN CENTER // EMPRESAS & CLIENTES',
      isEnabled: isCompany,
      description: 'Gestão multi-tenant de empresas, clientes e servidores',
      scopeTarget: isMaster ? 'TODAS AS EMPRESAS (ROOT)' : scopeTarget,
      permissions: [
        { id: 'p_ac_companies', name: 'Cadastrar / Gerenciar Empresas e Clientes', isEnabled: isMaster },
        { id: 'p_ac_users', name: 'Gestão de Usuários do Tenant', isEnabled: isCompany },
        { id: 'p_ac_infra', name: 'Configuração de Servidores e Discos Storage', isEnabled: isMaster }
      ]
    },
    {
      id: 'video_streams',
      name: 'FLUXOS DE VÍDEO // TOPOLOGIA & CÂMERAS',
      isEnabled: isCompany,
      description: 'Ingestão RTSP/ONVIF e organização em pastas',
      scopeTarget: isMaster ? 'TODAS AS CÂMERAS' : scopeTarget,
      permissions: [
        { id: 'p_vs_create', name: 'Cadastrar Novos Fluxos de Vídeo', isEnabled: isCompany },
        { id: 'p_vs_delete', name: 'Excluir Câmeras do Sistema', isEnabled: isMaster },
        { id: 'p_vs_folders', name: 'Criar e Gerenciar Pastas de Câmeras', isEnabled: isCompany }
      ]
    },
    {
      id: 'mosaic_live',
      name: 'MOSAICO AO VIVO // MONITORAMENTO',
      isEnabled: true,
      description: 'Grade de câmeras, controle PTZ e comutação sub-stream',
      scopeTarget: isMaster ? 'ACESSO TOTAL' : scopeTarget,
      permissions: [
        { id: 'p_ml_view', name: 'Visualização de Streams Ao Vivo', isEnabled: true },
        { id: 'p_ml_ptz', name: 'Controle de Movimento PTZ e Joystick', isEnabled: isOperator },
        { id: 'p_ml_hq', name: 'Alternar para Modo 4K / Alta Resolução', isEnabled: isOperator }
      ]
    },
    {
      id: 'recordings_timeline',
      name: 'GRAVAÇÕES & TIMELINE // FORENSE',
      isEnabled: isOperator,
      description: 'Gaveta inferior de playback, busca por horário e clipes',
      scopeTarget: isMaster ? 'RETENÇÃO TOTAL' : scopeTarget,
      permissions: [
        { id: 'p_rec_playback', name: 'Reprodução na Barra Temporal', isEnabled: isOperator },
        { id: 'p_rec_export', name: 'Exportação MP4 com Hash SHA-256', isEnabled: isOperator },
        { id: 'p_rec_purge', name: 'Exclusão Manual de Trechos de Vídeo', isEnabled: isMaster }
      ]
    },
    {
      id: 'alarms_sensors',
      name: 'ALARMES & SENSORES // DISPAROS & ZONAS',
      isEnabled: isOperator,
      description: 'Arme/desarme, sirenes, regras de zona e eventos',
      scopeTarget: isMaster ? 'TODAS AS ZONAS' : scopeTarget,
      permissions: [
        { id: 'p_alm_view', name: 'Visualizar Alertas e Detecção', isEnabled: true },
        { id: 'p_alm_ack', name: 'Reconhecer e Finalizar Incidentes', isEnabled: isOperator },
        { id: 'p_alm_config', name: 'Configurar Regras de Zona e Relays', isEnabled: isCompany }
      ]
    },
    {
      id: 'layouts_rounds',
      name: 'LAYOUTS & RONDAS // CARROSSEIS',
      isEnabled: true,
      description: 'Configuração de telas 2x2, 3x3, 4x4 e rotações',
      scopeTarget: isMaster ? 'LAYOUTS GLOBAIS' : scopeTarget,
      permissions: [
        { id: 'p_lay_view', name: 'Utilizar Layouts Existentes', isEnabled: true },
        { id: 'p_lay_create', name: 'Criar / Editar Layouts Customizados', isEnabled: isOperator }
      ]
    }
  ]
}
