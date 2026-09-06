import type { LogEntry } from '../types/systemLogs'

export const mockSystemLogs: LogEntry[] = [
  {
    id: 'log_001', timestamp: '2026-09-05 22:31:14', category: 'SYSTEM', level: 'INFO',
    actor: 'SYSTEM // CLUSTER_DAEMON', tenantId: 'GLOBAL', tenantName: 'SISTEMA GLOBAL',
    action: 'CLUSTER_SYNC_OK', target: 'HYDRASTREAM_NODE_02',
    details: 'Heartbeat sincronizado com sucesso via gRPC mesh (< 2ms)', ipAddress: '192.168.1.11'
  },
  {
    id: 'log_002', timestamp: '2026-09-05 22:28:40', category: 'AUDIT', level: 'INFO',
    actor: 'marcos.admin@alpha.corp', tenantId: 'tenant_alpha', tenantName: 'Alpha Seguranca',
    action: 'USER_CREATED', target: 'operador_lucas@alpha.corp',
    details: 'Criou novo operador com perfil restrito de monitoramento', ipAddress: '10.0.1.45'
  },
  {
    id: 'log_003', timestamp: '2026-09-05 22:24:19', category: 'SYSTEM', level: 'WARNING',
    actor: 'SYSTEM // STORAGEGUARD', tenantId: 'GLOBAL', tenantName: 'SISTEMA GLOBAL',
    action: 'BUFFER_HIGH_WATERMARK', target: 'NVME_HOT_BUFFER_01',
    details: 'Buffer atingiu 82% de ocupacao. Dreno automatico disparado para HDs', ipAddress: '127.0.0.1'
  },
  {
    id: 'log_004', timestamp: '2026-09-05 22:19:02', category: 'AUDIT', level: 'WARNING',
    actor: 'operador_lucas@alpha.corp', tenantId: 'tenant_alpha', tenantName: 'Alpha Seguranca',
    action: 'VIDEO_EXPORT_CLIP', target: 'CAM_ENTRADA_PRINCIPAL',
    details: 'Exportou trecho forense MP4 (15min) com hash SHA-256 de autenticidade', ipAddress: '10.0.1.88'
  },
  {
    id: 'log_005', timestamp: '2026-09-05 22:15:33', category: 'AUDIT', level: 'INFO',
    actor: 'superadmin@hydra.io', tenantId: 'GLOBAL', tenantName: 'SISTEMA GLOBAL',
    action: 'STORAGE_POOL_REGISTER', target: 'POOL_WARM_ARCHIVE_01',
    details: 'Registrou novo volume Western Digital 18TB para gravacoes', ipAddress: '192.168.1.5'
  },
  {
    id: 'log_006', timestamp: '2026-09-05 22:10:50', category: 'AUDIT', level: 'CRITICAL',
    actor: 'roberta.admin@beta.com', tenantId: 'tenant_beta', tenantName: 'Beta Logistica',
    action: 'AUTH_FAILED_LOCKOUT', target: 'roberta.admin@beta.com',
    details: '3 tentativas invalidas de login consecutivas. IP bloqueado por 15min', ipAddress: '187.54.12.9'
  },
  {
    id: 'log_007', timestamp: '2026-09-05 22:04:11', category: 'AUDIT', level: 'INFO',
    actor: 'cliente_rodrigo@beta.com', tenantId: 'tenant_beta', tenantName: 'Beta Logistica',
    action: 'REPORT_DOWNLOAD', target: 'RELATORIO_ALARME_AGOSTO.PDF',
    details: 'Download de relatorio forense consolidado via link seguro MinIO', ipAddress: '187.54.12.14'
  },
  {
    id: 'log_008', timestamp: '2026-09-05 21:58:22', category: 'SYSTEM', level: 'CRITICAL',
    actor: 'SYSTEM // GPU_MONITOR', tenantId: 'GLOBAL', tenantName: 'SISTEMA GLOBAL',
    action: 'GPU_TEMP_ALERT', target: 'NVIDIA RTX 5090 (CUDA 0)',
    details: 'Temperatura da GPU atingiu 84C sob inferencia SAHI continua', ipAddress: '127.0.0.1'
  },
  {
    id: 'log_009', timestamp: '2026-09-05 21:45:00', category: 'AUDIT', level: 'INFO',
    actor: 'marcos.admin@alpha.corp', tenantId: 'tenant_alpha', tenantName: 'Alpha Seguranca',
    action: 'AI_RULE_UPDATE', target: 'REGRA_INVASAO_PERIMETRO',
    details: 'Alterou sensibilidade do modelo YOLO para alvos pequenos de 0.45 para 0.55', ipAddress: '10.0.1.45'
  },
  {
    id: 'log_010', timestamp: '2026-09-05 21:30:15', category: 'SYSTEM', level: 'INFO',
    actor: 'SYSTEM // ONVIF_ENGINE', tenantId: 'GLOBAL', tenantName: 'SISTEMA GLOBAL',
    action: 'DISCOVERY_PROBE_RUN', target: '239.255.255.250:3702',
    details: 'Sonda multicast ONVIF WS-Discovery identificou 4 novas cameras na rede local', ipAddress: '192.168.1.1'
  }
]
