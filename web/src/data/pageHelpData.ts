import type { PageHelpGuide } from '../types/helpTypes'

export const pageHelpCatalog: Record<string, PageHelpGuide> = {
  video_streams: {
    pageId: 'video_streams', title: 'GESTÃO DE FLUXOS & CÂMERAS',
    subtitle: 'INGESTÃO TCP RFC 2326, ONVIF WS-DISCOVERY & VIRTUAL LOOP', badge: 'INGESTÃO',
    sections: [
      {
        id: 'overview', title: 'Visão Geral e Descoberta', tag: 'CONECTIVIDADE',
        description: 'Centraliza a conexão com todas as câmeras IP, encoders e arquivos em loop com teste de socket ao vivo.',
        steps: [
          'Radar ONVIF escuta e lista câmeras na rede local automaticamente.',
          'Clique em [+ NOVO FLUXO] para cadastrar câmeras RTSP, ONVIF, RTMP ou LOOP.',
          'Organize câmeras em pastas por condomínio, cliente ou setor.'
        ],
        animationType: 'streams_scan',
        example: { label: 'CASO DE USO // CÂMERA DE PORTARIA', scenario: 'IP 192.168.1.100 RTSP 554 com usuário e senha.', result: 'Fluxo validado com snapshot ao vivo e codificação H.264/H.265.' }
      },
      {
        id: 'import_export', title: 'Exportação & Importação Segura', tag: 'ZERO-TRUST',
        description: 'Exporte listas completas em Planilhas Excel (.csv) com URLs mascaradas para proteger credenciais.',
        steps: [
          'Clique no botão Exportar para gerar planilha CSV protegida.',
          'As senhas físicas são ocultadas pela URL proxy do HydraStream.',
          'Importe arquivos CSV, TXT ou JSON em lote com 1-clique.'
        ],
        animationType: 'streams_scan',
        example: { label: 'EXEMPLO CSV // IMPORTAÇÃO', scenario: 'Planilha com 50 câmeras gerada pelo integrador.', result: '50 fluxos criados instantaneamente na pasta do cliente.' }
      }
    ]
  },
  layouts: {
    pageId: 'layouts', title: 'GESTÃO DE LAYOUTS DE CÂMERAS',
    subtitle: 'MATRIZES DE VÍDEO (1x1 A 4x4) & DISTRIBUIÇÃO MULTI-TENANT', badge: 'MOSAICOS',
    sections: [
      {
        id: 'overview', title: 'Criação de Grades e Slots', tag: 'MATRIZ DE VÍDEO',
        description: 'Crie mosaicos personalizados com até 16 câmeras simultâneas otimizadas com comutação de sub-stream.',
        steps: [
          'Crie uma pasta com o nome do cliente ou setor da empresa.',
          'Clique em [+ NOVO LAYOUT] e escolha a grade (1x1, 2x2, 3x3, 4x4).',
          'Arraste as câmeras para cada slot e trave com o Cadeado.'
        ],
        animationType: 'grid_slots',
        example: { label: 'EXEMPLO // MONITORAMENTO 2x2', scenario: 'Grade com Portaria, Doca, Entrada e Estacionamento.', result: 'Layout exibido instantaneamente na aba de vigilância do cliente.' }
      }
    ]
  },
  users: {
    pageId: 'users', title: 'USUÁRIOS & PERMISSÕES (RBAC)',
    subtitle: 'CONTROLE DE ACESSO HIERÁRQUICO & ISOLAMENTO DE CLIENTES', badge: 'SEGURANÇA',
    sections: [
      {
        id: 'overview', title: 'Hierarquia de Acesso e Escopo', tag: 'RBAC',
        description: 'Defina quem pode acessar cada recurso. Cada operador visualiza estritamente as câmeras autorizadas.',
        steps: [
          'Admin Master: Acesso total a todas as empresas e configurações.',
          'Gestor de Empresa: Administra apenas os operadores do seu tenant.',
          'Operador: Acesso restrito apenas ao Mosaico ao Vivo e Playback.'
        ],
        animationType: 'rbac_tree',
        example: { label: 'EXEMPLO // OPERADOR NOTURNO', scenario: 'Usuário restrito à empresa SHOPPING ALPHA.', result: 'Não enxerga câmeras de outros clientes nem altera configs.' }
      }
    ]
  },
  branding: {
    pageId: 'branding', title: 'LOGOMARCA & IDENTIDADE VISUAL',
    subtitle: 'CUSTOMIZAÇÃO WHITE-LABEL // LOGOTIPOS, CORES E NOMES', badge: 'WHITE-LABEL',
    sections: [
      {
        id: 'overview', title: 'Personalização da Marca', tag: 'WHITE-LABEL',
        description: 'Substitua a identidade visual do sistema pela marca da sua empresa de segurança ou integrador.',
        steps: [
          'Altere o Nome do Sistema e do Painel Administrativo.',
          'Envie o ícone da barra superior e o Favicon da aba do navegador.',
          'Escolha a cor primária de destaque e visualize no simulador.'
        ],
        animationType: 'brand_palette',
        example: { label: 'EXEMPLO // VISION SECURITY', scenario: 'Upload do logo SVG da Vision Security e cor azul (#3b82f6).', result: 'Aba do navegador, header e tela de login exibem a marca própria.' }
      }
    ]
  },
  workflows: {
    pageId: 'workflows', title: 'WORKFLOWS & AUTOMAÇÃO',
    subtitle: 'DISPARO DE ALARMES, INTEGRAÇÃO COM TELEGRAM & WEBHOOKS', badge: 'AUTOMAÇÃO',
    sections: [
      {
        id: 'overview', title: 'Gatilhos e Ações Automáticas', tag: 'EVENT PIPELINE',
        description: 'Crie fluxos que reagem a detecções de IA e enviam fotos para o Telegram ou Webhooks em tempo real.',
        steps: [
          'Selecione o gatilho de evento (ex: Intrusão Noturna, Linha Virtual).',
          'Configure a ação de envio instantâneo com foto para o Telegram.',
          'Defina a janela anti-spam para evitar notificações repetidas.'
        ],
        animationType: 'flow_pipes',
        example: { label: 'EXEMPLO // ALARME NOTURNO', scenario: 'Pessoa na área restrita após as 22h00.', result: 'Foto anotada despachada para o Telegram em < 200ms.' }
      }
    ]
  },
  storage: {
    pageId: 'storage', title: 'ARMAZENAMENTO & DISCOS',
    subtitle: 'NVME HOT BUFFER, AUTO-PURGE & STORAGEGUARD', badge: 'STORAGE',
    sections: [
      {
        id: 'overview', title: 'Gestão Inteligente de Armazenamento', tag: 'ANTI-CRASH',
        description: 'Gerencie pools de gravação contínua e alarmes com proteção automática contra estouro de disco.',
        steps: [
          'Monitore o espaço livre dos discos NVMe e HDs mecânicos.',
          'StorageGuard purga automaticamente gravações antigas aos 80% e 90%.',
          'Gravações com trava forense (is_pinned) nunca são excluídas.'
        ],
        animationType: 'storage_pool',
        example: { label: 'EXEMPLO // RETENÇÃO 30 DIAS', scenario: 'Buffer de 2TB gravando 16 canais Full-HD em H.265.', result: 'Dreno contínuo e preservação automática dos clipes de alarme.' }
      }
    ]
  },
  performance: {
    pageId: 'performance', title: 'PERFORMANCE & NÓS DO CLUSTER',
    subtitle: 'INGESTÃO ZERO-COPY /DEV/SHM & IA GPU NA RTX 5090', badge: 'TELEMETRIA',
    sections: [
      {
        id: 'overview', title: 'Cluster e Aceleração por Hardware', tag: 'HARDWARE',
        description: 'Visualize a telemetria dos nós de ingestão HydraStream e processamento de IA na RTX 5090.',
        steps: [
          'Monitore o throughput do Ring Buffer de memória /dev/shm.',
          'Acompanhe o consumo de VRAM e FPS de inferência YOLO.',
          'Conecte novos servidores Edge ao cluster com descoberta automática.'
        ],
        animationType: 'cluster_nodes',
        example: { label: 'EXEMPLO // RTX 5090', scenario: '32 streams simultâneos com análise SAHI.', result: 'Inferência a mais de 18.500 FPS com 4ms de latência.' }
      }
    ]
  }
}
