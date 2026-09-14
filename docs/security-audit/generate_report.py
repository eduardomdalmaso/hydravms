#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Script de Geração do Relatório de Auditoria de Segurança - HydraVMS
Gera o PDF visual completo em pt-BR em docs/security-audit/relatorio-auditoria-seguranca.pdf
"""

import os
import sys
from datetime import datetime
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
from reportlab.lib.pagesizes import A4
from reportlab.lib import colors
from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle
from reportlab.lib.units import cm, mm
from reportlab.platypus import (
    SimpleDocTemplate, Paragraph, Spacer, Table, TableStyle, Image, KeepTogether, PageBreak, HRFlowable
)
from reportlab.pdfgen import canvas

# Paleta de Cores da Auditoria
COLOR_CRITICA = colors.HexColor('#B91C1C')
COLOR_ALTA = colors.HexColor('#EA580C')
COLOR_MEDIA = colors.HexColor('#D97706')
COLOR_BAIXA = colors.HexColor('#2563EB')
COLOR_FORTE = colors.HexColor('#059669')
COLOR_DARK = colors.HexColor('#0F172A')
COLOR_SURFACE = colors.HexColor('#1E293B')
COLOR_BG_LIGHT = colors.HexColor('#F8FAFC')
COLOR_BORDER = colors.HexColor('#E2E8F0')
COLOR_TEXT = colors.HexColor('#334155')

class NumberedCanvas(canvas.Canvas):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self._saved_page_states = []

    def showPage(self):
        self._saved_page_states.append(dict(self.__dict__))
        self._startPage()

    def save(self):
        num_pages = len(self._saved_page_states)
        for state in self._saved_page_states:
            self.__dict__.update(state)
            self.draw_header_footer(num_pages)
            super().showPage()
        super().save()

    def draw_header_footer(self, page_count):
        self.saveState()
        if self._pageNumber > 1:
            # Header
            self.setFont("Helvetica-Bold", 8)
            self.setFillColor(COLOR_DARK)
            self.drawString(20 * mm, 285 * mm, "HYDRAVMS // RELATÓRIO DE AUDITORIA DE SEGURANÇA")
            self.setFont("Helvetica", 8)
            self.setFillColor(colors.HexColor('#64748B'))
            self.drawRightString(190 * mm, 285 * mm, "CONFIDENCIAL • CONTROLE INTERNO")
            self.setStrokeColor(COLOR_BORDER)
            self.setLineWidth(0.5)
            self.line(20 * mm, 282 * mm, 190 * mm, 282 * mm)

        # Footer
        self.setFont("Helvetica", 8)
        self.setFillColor(colors.HexColor('#64748B'))
        self.drawString(20 * mm, 12 * mm, "Hydra Vision Ecosystem • Segurança da Informação")
        self.drawRightString(190 * mm, 12 * mm, f"Página {self._pageNumber} de {page_count}")
        self.setStrokeColor(COLOR_BORDER)
        self.setLineWidth(0.5)
        self.line(20 * mm, 15 * mm, 190 * mm, 15 * mm)
        self.restoreState()

def generate_charts(output_dir):
    os.makedirs(output_dir, exist_ok=True)
    
    # 1. Gráfico de Rosca por Severidade
    fig, ax = plt.subplots(figsize=(4.2, 3.2), subplot_kw=dict(aspect="equal"))
    labels = ['Crítica (3)', 'Alta (5)', 'Média (4)', 'Baixa (1)']
    sizes = [3, 5, 4, 1]
    chart_colors = ['#B91C1C', '#EA580C', '#D97706', '#2563EB']
    
    wedges, texts, autotexts = ax.pie(
        sizes, labels=labels, autopct='%1.0f%%', startangle=140,
        colors=chart_colors, pctdistance=0.75,
        textprops=dict(color="#0F172A", fontsize=8, fontweight='bold'),
        wedgeprops=dict(width=0.45, edgecolor='white', linewidth=2)
    )
    for at in autotexts:
        at.set_color('white')
        at.set_fontsize(9)
    ax.set_title("Achados por Severidade", fontsize=11, fontweight='bold', pad=10, color='#0F172A')
    plt.tight_layout()
    donut_path = os.path.join(output_dir, "chart_donut_severity.png")
    plt.savefig(donut_path, dpi=200, transparent=True)
    plt.close()

    # 2. Gráfico de Barras por Categoria
    fig, ax = plt.subplots(figsize=(4.8, 3.2))
    categories = [
        '1. Banco sem\nTranca',
        '2. Permissão\nNavegador',
        '3. IDOR',
        '4. Chaves\nExpostas',
        '5. Inputs/\nXSS/Files'
    ]
    counts = [4, 2, 4, 3, 3]
    bar_colors = ['#B91C1C', '#EA580C', '#EA580C', '#B91C1C', '#D97706']
    
    bars = ax.bar(categories, counts, color=bar_colors, width=0.55, edgecolor='none', zorder=3)
    ax.set_title("Achados por Categoria", fontsize=11, fontweight='bold', pad=10, color='#0F172A')
    ax.set_ylabel("Quantidade de Falhas", fontsize=8, color='#334155')
    ax.tick_params(axis='x', labelsize=7.5)
    ax.tick_params(axis='y', labelsize=8)
    ax.set_ylim(0, 5.5)
    ax.grid(axis='y', linestyle='--', alpha=0.4, zorder=0)
    for spine in ['top', 'right', 'left', 'bottom']:
        ax.spines[spine].set_color('#CBD5E1')
    
    for bar in bars:
        yval = bar.get_height()
        ax.text(bar.get_x() + bar.get_width()/2.0, yval + 0.15, int(yval), ha='center', va='bottom', fontsize=9, fontweight='bold', color='#0F172A')

    plt.tight_layout()
    bar_path = os.path.join(output_dir, "chart_bar_categories.png")
    plt.savefig(bar_path, dpi=200, transparent=True)
    plt.close()

    return donut_path, bar_path

def build_pdf():
    pdf_path = os.path.join("docs", "security-audit", "relatorio-auditoria-seguranca.pdf")
    charts_dir = os.path.join("docs", "security-audit", "assets")
    donut_img, bar_img = generate_charts(charts_dir)

    doc = SimpleDocTemplate(
        pdf_path,
        pagesize=A4,
        leftMargin=18 * mm,
        rightMargin=18 * mm,
        topMargin=20 * mm,
        bottomMargin=20 * mm
    )

    styles = getSampleStyleSheet()
    
    # Custom Typography Styles
    title_style = ParagraphStyle(
        'CoverTitle',
        parent=styles['Heading1'],
        fontName='Helvetica-Bold',
        fontSize=24,
        leading=28,
        textColor=COLOR_DARK,
        spaceAfter=6
    )
    subtitle_style = ParagraphStyle(
        'CoverSubtitle',
        parent=styles['Normal'],
        fontName='Helvetica',
        fontSize=12,
        leading=16,
        textColor=colors.HexColor('#64748B'),
        spaceAfter=15
    )
    h1_style = ParagraphStyle(
        'Heading1Custom',
        parent=styles['Heading1'],
        fontName='Helvetica-Bold',
        fontSize=14,
        leading=18,
        textColor=COLOR_DARK,
        spaceBefore=12,
        spaceAfter=8,
        keepWithNext=True
    )
    h2_style = ParagraphStyle(
        'Heading2Custom',
        parent=styles['Heading2'],
        fontName='Helvetica-Bold',
        fontSize=11,
        leading=15,
        textColor=COLOR_SURFACE,
        spaceBefore=8,
        spaceAfter=4,
        keepWithNext=True
    )
    body_style = ParagraphStyle(
        'BodyCustom',
        parent=styles['Normal'],
        fontName='Helvetica',
        fontSize=8.5,
        leading=12,
        textColor=COLOR_TEXT,
        spaceAfter=5
    )
    code_style = ParagraphStyle(
        'CodeSnippet',
        parent=styles['Normal'],
        fontName='Courier',
        fontSize=7.5,
        leading=9.5,
        textColor=colors.HexColor('#0F172A'),
        backColor=colors.HexColor('#F1F5F9'),
        borderColor=colors.HexColor('#CBD5E1'),
        borderWidth=0.5,
        borderPadding=4,
        spaceBefore=3,
        spaceAfter=4
    )
    issue_code_style = ParagraphStyle(
        'IssueMarkdown',
        parent=styles['Normal'],
        fontName='Courier',
        fontSize=7,
        leading=8.5,
        textColor=colors.HexColor('#1E293B'),
        backColor=colors.HexColor('#F8FAFC'),
        borderColor=colors.HexColor('#E2E8F0'),
        borderWidth=0.5,
        borderPadding=5,
        spaceBefore=4,
        spaceAfter=6
    )
    table_cell_style = ParagraphStyle(
        'TableCell',
        parent=styles['Normal'],
        fontName='Helvetica',
        fontSize=7.5,
        leading=10,
        textColor=COLOR_TEXT
    )
    table_cell_bold = ParagraphStyle(
        'TableCellBold',
        parent=styles['Normal'],
        fontName='Helvetica-Bold',
        fontSize=7.5,
        leading=10,
        textColor=COLOR_DARK
    )

    story = []

    # =========================================================================
    # 1. CAPA / CABEÇALHO EXECUTIVO
    # =========================================================================
    badge_data = [[
        Paragraph("<font color='#B91C1C'><b>RELATÓRIO TÉCNICO DE SEGURANÇA E CONFORMIDADE</b></font>", body_style),
        Paragraph("<font color='#64748B'><b>VERSÃO 1.0 • SETEMBRO 2026</b></font>", body_style)
    ]]
    t_badge = Table(badge_data, colWidths=[100*mm, 74*mm])
    t_badge.setStyle(TableStyle([
        ('VALIGN', (0,0), (-1,-1), 'MIDDLE'),
        ('ALIGN', (1,0), (1,0), 'RIGHT'),
        ('BOTTOMPADDING', (0,0), (-1,-1), 0),
        ('TOPPADDING', (0,0), (-1,-1), 0),
    ]))
    story.append(t_badge)
    story.append(Spacer(1, 4*mm))
    story.append(HRFlowable(width="100%", thickness=2, color=COLOR_CRITICA, spaceAfter=8*mm))

    story.append(Paragraph("Relatório de Auditoria de Segurança — HydraVMS", title_style))
    story.append(Paragraph("Avaliação Estrita de Código-Fonte, Controle de Acesso, Isolamento Multi-Tenant e Vulnerabilidades Críticas", subtitle_style))

    # Tabela Metadados da Auditoria
    meta_table_data = [
        [
            Paragraph("<b>Projeto Auditado:</b> HydraVMS (Control Plane & Web UI)", table_cell_style),
            Paragraph("<b>Data da Auditoria:</b> 14 de Setembro de 2026", table_cell_style)
        ],
        [
            Paragraph("<b>Linguagens & Runtimes:</b> Go 1.22+ (Backend) / TypeScript & Vue 3 (Frontend)", table_cell_style),
            Paragraph("<b>Banco de Dados:</b> PostgreSQL 15+ (pgxpool) & MinIO S3", table_cell_style)
        ],
        [
            Paragraph("<b>Mensageria / Eventos:</b> NATS JetStream 2.10", table_cell_style),
            Paragraph("<b>Mecanismo de Auth:</b> JWT HMAC-SHA256 & Context Injection", table_cell_style)
        ],
    ]
    t_meta = Table(meta_table_data, colWidths=[87*mm, 87*mm])
    t_meta.setStyle(TableStyle([
        ('BACKGROUND', (0,0), (-1,-1), COLOR_BG_LIGHT),
        ('BOX', (0,0), (-1,-1), 0.5, COLOR_BORDER),
        ('INNERGRID', (0,0), (-1,-1), 0.5, COLOR_BORDER),
        ('TOPPADDING', (0,0), (-1,-1), 3*mm),
        ('BOTTOMPADDING', (0,0), (-1,-1), 3*mm),
        ('LEFTPADDING', (0,0), (-1,-1), 4*mm),
        ('RIGHTPADDING', (0,0), (-1,-1), 4*mm),
    ]))
    story.append(t_meta)
    story.append(Spacer(1, 6*mm))

    # Nota Metodológica
    story.append(Paragraph("Nota Metodológica & Mapeamento da Stack", h2_style))
    story.append(Paragraph(
        "A auditoria foi conduzida através de análise estática e revisão de fluxo de controle linha por linha em 100% dos handlers HTTP, "
        "repositórios relacionais PostgreSQL, middlewares de autenticação, barramento de eventos NATS e componentes Vue 3 do HydraVMS. "
        "As 5 categorias de segurança solicitadas foram mapeadas diretamente para as construções da stack:<br/>"
        "• <b>1. Banco sem Tranca:</b> Verificação da obrigatoriedade do <font name='Courier'>tenant_id</font> extraído do token JWT no contexto nas queries SQL com pgx.<br/>"
        "• <b>2. Permissão no Navegador:</b> Confronto entre o gate frontend (<font name='Courier'>isAdmin</font>) e ausência de verificação de role no backend Go.<br/>"
        "• <b>3. IDOR:</b> Rastreamento de parâmetros de rota (<font name='Courier'>{id}</font>) executando <font name='Courier'>DELETE/GET/PUT</font> sem validação de posse.<br/>"
        "• <b>4. Chaves Expostas:</b> Auditoria de segredos estáticos (JWT_SECRET fallback, senhas padrão de admin no boot e configs de S3/DB).<br/>"
        "• <b>5. Inputs / XSS / Injeção:</b> Auditoria de <font name='Courier'>http.ServeFile</font> com chaves de entrada, SSRF em sondas de rede e templates.",
        body_style
    ))
    story.append(Spacer(1, 6*mm))

    # =========================================================================
    # 2. RESUMO EXECUTIVO & GRÁFICOS
    # =========================================================================
    story.append(Paragraph("Resumo Executivo", h1_style))
    story.append(Paragraph(
        "Foram identificados <b>13 achados de segurança</b> no código auditado, dos quais <b>3 são de severidade Crítica</b> e <b>5 de severidade Alta</b>. "
        "O sistema possui uma fundação arquitetural bem desenhada (arquitetura hexagonal, consultas parametrizadas contra SQL Injection, "
        "isolamento de WebSockets por tenant), porém sofre de vulnerabilidades de bypass de autorização e segredos padrão que comprometem o isolamento multi-tenant.",
        body_style
    ))
    story.append(Spacer(1, 3*mm))

    # Tabela com Gráficos lado a lado
    chart_table_data = [[
        Image(donut_img, width=82*mm, height=58*mm),
        Image(bar_img, width=88*mm, height=58*mm)
    ]]
    t_charts = Table(chart_table_data, colWidths=[85*mm, 89*mm])
    t_charts.setStyle(TableStyle([
        ('ALIGN', (0,0), (-1,-1), 'CENTER'),
        ('VALIGN', (0,0), (-1,-1), 'MIDDLE'),
        ('LEFTPADDING', (0,0), (-1,-1), 0),
        ('RIGHTPADDING', (0,0), (-1,-1), 0),
        ('TOPPADDING', (0,0), (-1,-1), 0),
        ('BOTTOMPADDING', (0,0), (-1,-1), 0),
    ]))
    story.append(t_charts)
    story.append(Spacer(1, 4*mm))

    # =========================================================================
    # 3. PONTOS FORTES E PONTOS FRACOS
    # =========================================================================
    story.append(Paragraph("Pontos Fortes e Riscos Centrais", h1_style))
    
    strengths_weaknesses_data = [
        [
            Paragraph("<font color='#059669'><b>PONTOS FORTES (BOAS PRÁTICAS OBSERVADAS)</b></font>", table_cell_bold),
            Paragraph("<font color='#B91C1C'><b>PONTOS FRACOS (VULNERABILIDADES CENTRAIS)</b></font>", table_cell_bold)
        ],
        [
            Paragraph(
                "• <b>Proteção Total contra SQL Injection:</b> Todos os repositórios (<font name='Courier'>camera_repository.go</font>, <font name='Courier'>folder_repository.go</font>) utilizam consultas 100% parametrizadas com <font name='Courier'>$1, $2...</font> via pgx/v5.<br/>"
                "• <b>Autenticação Token-Only Obrigatória:</b> <font name='Courier'>AuthMiddleware</font> cobre todos os endpoints e valida JWT em todas as rotas.<br/>"
                "• <b>Frontend Imune a XSS Clássico:</b> O Vue 3 não utiliza <font name='Courier'>v-html</font> nem inserção de HTML bruto.<br/>"
                "• <b>WebSockets Isolados:</b> O <font name='Courier'>ws.Hub</font> separa canais por tenant UUID.",
                table_cell_style
            ),
            Paragraph(
                "• <b>Bypass de Chave Secreta JWT:</b> Segredo de assinatura hardcoded estático caso variável de ambiente esteja vazia.<br/>"
                "• <b>Falta de RBAC no Backend:</b> Endpoints administrativos não validam a role do usuário no servidor.<br/>"
                "• <b>Vazamento Multi-Tenant em Câmeras/Usuários:</b> Consultas sem filtro de tenant ou com cláusula <font name='Courier'>OR tenant_id = root</font>.<br/>"
                "• <b>Arbitrary File Read em Gravações:</b> <font name='Courier'>http.ServeFile</font> permite carregar arquivos do SO.",
                table_cell_style
            )
        ]
    ]
    t_sw = Table(strengths_weaknesses_data, colWidths=[87*mm, 87*mm])
    t_sw.setStyle(TableStyle([
        ('BACKGROUND', (0,0), (0,0), colors.HexColor('#ECFDF5')),
        ('BACKGROUND', (1,0), (1,0), colors.HexColor('#FEF2F2')),
        ('BACKGROUND', (0,1), (0,1), COLOR_BG_LIGHT),
        ('BACKGROUND', (1,1), (1,1), COLOR_BG_LIGHT),
        ('BOX', (0,0), (-1,-1), 0.5, COLOR_BORDER),
        ('INNERGRID', (0,0), (-1,-1), 0.5, COLOR_BORDER),
        ('TOPPADDING', (0,0), (-1,-1), 3*mm),
        ('BOTTOMPADDING', (0,0), (-1,-1), 3*mm),
        ('LEFTPADDING', (0,0), (-1,-1), 4*mm),
        ('RIGHTPADDING', (0,0), (-1,-1), 4*mm),
    ]))
    story.append(t_sw)

    story.append(PageBreak())

    # =========================================================================
    # 4. TABELA DE ACHADOS DETALHADOS POR CATEGORIA
    # =========================================================================
    story.append(Paragraph("Tabela de Achados Detalhados por Categoria", h1_style))
    
    findings_data = [
        [
            Paragraph("<b>SEV.</b>", table_cell_bold),
            Paragraph("<b>CAT.</b>", table_cell_bold),
            Paragraph("<b>ARQUIVO : LINHA</b>", table_cell_bold),
            Paragraph("<b>DESCRIÇÃO DA FALHA & IMPACTO</b>", table_cell_bold)
        ],
        # Cat 1
        [
            Paragraph("<font color='#B91C1C'><b>CRÍTICA</b></font>", table_cell_style),
            Paragraph("1. Banco", table_cell_style),
            Paragraph("<font name='Courier'>user_handler.go:29-48</font>", table_cell_style),
            Paragraph("<b>Listagem global de usuários sem filtro de tenant.</b> Endpoint <font name='Courier'>GET /api/v1/users</font> vaza todos os e-mails, nomes e escopos de todas as empresas.", table_cell_style)
        ],
        [
            Paragraph("<font color='#EA580C'><b>ALTA</b></font>", table_cell_style),
            Paragraph("1. Banco", table_cell_style),
            Paragraph("<font name='Courier'>audit_handler.go:30<br/>audit_log_repo.go:89</font>", table_cell_style),
            Paragraph("<b>Logs de auditoria forense sem tenant scoping obrigatório.</b> TenantID lido da query string do cliente. Vaza histórico de ações de todos os clientes.", table_cell_style)
        ],
        [
            Paragraph("<font color='#EA580C'><b>ALTA</b></font>", table_cell_style),
            Paragraph("1. Banco", table_cell_style),
            Paragraph("<font name='Courier'>camera_repo.go:66,102<br/>folder_repo.go:50,77</font>", table_cell_style),
            Paragraph("<b>Cláusula de bypass do Tenant Raiz.</b> Cláusula <font name='Courier'>OR tenant_id = root</font> permite que inquilinos secundários visualizem e excluam câmeras da matriz.", table_cell_style)
        ],
        [
            Paragraph("<font color='#D97706'><b>MÉDIA</b></font>", table_cell_style),
            Paragraph("1. Banco", table_cell_style),
            Paragraph("<font name='Courier'>cluster_node_handler.go:218</font>", table_cell_style),
            Paragraph("<b>Hardcode do Tenant Raiz na gestão de nós.</b> Handler força UUID raiz fixo ignorando o tenant autenticado no token.", table_cell_style)
        ],
        # Cat 2
        [
            Paragraph("<font color='#B91C1C'><b>CRÍTICA</b></font>", table_cell_style),
            Paragraph("2. Permissão", table_cell_style),
            Paragraph("<font name='Courier'>router.go:60-100<br/>App.vue:59,90</font>", table_cell_style),
            Paragraph("<b>Operações administrativas desprotegidas no Backend.</b> Frontend restringe Admin Center com <font name='Courier'>v-if='isAdmin'</font>, mas o servidor Go não valida se a role é admin nas rotas de Storage, Nós, Usuários e Logs.", table_cell_style)
        ],
        [
            Paragraph("<font color='#D97706'><b>MÉDIA</b></font>", table_cell_style),
            Paragraph("2. Permissão", table_cell_style),
            Paragraph("<font name='Courier'>audit_handler.go:92-147</font>", table_cell_style),
            Paragraph("<b>Forjamento de logs de auditoria por operadores.</b> Rota <font name='Courier'>POST /api/v1/audit/logs</font> permite inserção arbitrária de logs sem checar papel admin/sistema.", table_cell_style)
        ],
        # Cat 3
        [
            Paragraph("<font color='#EA580C'><b>ALTA</b></font>", table_cell_style),
            Paragraph("3. IDOR", table_cell_style),
            Paragraph("<font name='Courier'>cluster_node_handler.go:308</font>", table_cell_style),
            Paragraph("<b>Exclusão de Nós de Cluster por IDOR.</b> <font name='Courier'>DELETE /api/v1/cluster/nodes/{id}</font> deleta o nó sem validar tenant ou permissão do nó.", table_cell_style)
        ],
        [
            Paragraph("<font color='#EA580C'><b>ALTA</b></font>", table_cell_style),
            Paragraph("3. IDOR", table_cell_style),
            Paragraph("<font name='Courier'>storage_pool_handler.go:85</font>", table_cell_style),
            Paragraph("<b>Exclusão de Storage Pools por IDOR.</b> <font name='Courier'>DELETE /api/v1/storage/pools/{id}</font> remove pools de armazenamento sem controle de posse.", table_cell_style)
        ],
        [
            Paragraph("<font color='#EA580C'><b>ALTA</b></font>", table_cell_style),
            Paragraph("3. IDOR", table_cell_style),
            Paragraph("<font name='Courier'>storage_pool_handler.go:120</font>", table_cell_style),
            Paragraph("<b>Geração de Presigned URL arbitrária no MinIO S3.</b> Permite obter link assinado para qualquer bucket e chave de objeto S3 de qualquer inquilino.", table_cell_style)
        ],
        # Cat 4
        [
            Paragraph("<font color='#B91C1C'><b>CRÍTICA</b></font>", table_cell_style),
            Paragraph("4. Chaves", table_cell_style),
            Paragraph("<font name='Courier'>auth_service.go:30-38</font>", table_cell_style),
            Paragraph("<b>Chave Secreta JWT Estática em Fallback.</b> Uso de string fixa no código fonte quando JWT_SECRET não é setado. Permite forjar tokens de superadmin.", table_cell_style)
        ],
        [
            Paragraph("<font color='#EA580C'><b>ALTA</b></font>", table_cell_style),
            Paragraph("4. Chaves", table_cell_style),
            Paragraph("<font name='Courier'>user_repository.go:42-57</font>", table_cell_style),
            Paragraph("<b>Reset Automático de Senha de Admin no Startup.</b> Boot do servidor executa upsert do usuário admin com a senha fraca 'admin' via bcrypt.", table_cell_style)
        ],
        # Cat 5
        [
            Paragraph("<font color='#B91C1C'><b>CRÍTICA</b></font>", table_cell_style),
            Paragraph("5. Inputs", table_cell_style),
            Paragraph("<font name='Courier'>camera_handler.go:252-308</font>", table_cell_style),
            Paragraph("<b>Arbitrary File Read / Path Traversal em Gravações.</b> <font name='Courier'>http.ServeFile</font> serve o caminho fornecido em <font name='Courier'>s3_key</font> sem validação de diretório permitido.", table_cell_style)
        ],
        [
            Paragraph("<font color='#D97706'><b>MÉDIA</b></font>", table_cell_style),
            Paragraph("5. Inputs", table_cell_style),
            Paragraph("<font name='Courier'>cluster_node_handler.go:344</font>", table_cell_style),
            Paragraph("<b>SSRF em Sonda de Rede (Probe).</b> Endpoint aceita IP e porta arbitrários e dispara requisições HTTP internas sem sanitização.", table_cell_style)
        ],
    ]

    t_findings = Table(findings_data, colWidths=[18*mm, 20*mm, 42*mm, 94*mm])
    t_findings.setStyle(TableStyle([
        ('BACKGROUND', (0,0), (-1,0), COLOR_DARK),
        ('TEXTCOLOR', (0,0), (-1,0), colors.white),
        ('BOX', (0,0), (-1,-1), 0.5, COLOR_BORDER),
        ('INNERGRID', (0,0), (-1,-1), 0.5, COLOR_BORDER),
        ('TOPPADDING', (0,0), (-1,-1), 2*mm),
        ('BOTTOMPADDING', (0,0), (-1,-1), 2*mm),
        ('LEFTPADDING', (0,0), (-1,-1), 2.5*mm),
        ('RIGHTPADDING', (0,0), (-1,-1), 2.5*mm),
        ('VALIGN', (0,0), (-1,-1), 'TOP'),
    ]))
    for i in range(1, len(findings_data)):
        if i % 2 == 0:
            t_findings.setStyle(TableStyle([('BACKGROUND', (0,i), (-1,i), COLOR_BG_LIGHT)]))

    story.append(t_findings)
    story.append(Spacer(1, 6*mm))

    # =========================================================================
    # 5. RECOMENDAÇÕES PRIORIZADAS (P1, P2, P3)
    # =========================================================================
    story.append(Paragraph("Plano de Remediação Priorizado", h1_style))
    
    recom_data = [
        [
            Paragraph("<font color='#B91C1C'><b>PRIORIDADE 1 (P1) — IMEDIATO (24H)</b></font>", table_cell_bold),
            Paragraph(
                "1. <b>Exigir JWT_SECRET obrigatório no startup:</b> Abortar inicialização (<font name='Courier'>log.Fatalf</font>) se a variável for vazia ou menor que 32 caracteres.<br/>"
                "2. <b>Corrigir isolamento de Usuários:</b> Adicionar <font name='Courier'>WHERE tenant_id = $1</font> em <font name='Courier'>HandleUsers</font>.<br/>"
                "3. <b>Eliminar Arbitrary File Read:</b> Validar <font name='Courier'>filepath.Clean</font> e garantir que o arquivo esteja restrito a <font name='Courier'>storage/</font> antes de <font name='Courier'>http.ServeFile</font>.<br/>"
                "4. <b>Remover sobrescrita de senha 'admin':</b> Criar o admin padrão apenas se a tabela estiver vazia.",
                table_cell_style
            )
        ],
        [
            Paragraph("<font color='#EA580C'><b>PRIORIDADE 2 (P2) — ALTO (72H)</b></font>", table_cell_bold),
            Paragraph(
                "1. <b>Implementar Middleware RequireRole('admin'):</b> Proteger rotas de usuários, storage, cluster e logs no <font name='Courier'>router.go</font>.<br/>"
                "2. <b>Sanitizar IDORs em Storage e Cluster:</b> Exigir validação de tenant em exclusões e assinaturas de URLs S3.<br/>"
                "3. <b>Remover bypass 'OR tenant_id = root':</b> Isolar estritamente recursos do tenant autenticado.",
                table_cell_style
            )
        ],
        [
            Paragraph("<font color='#D97706'><b>PRIORIDADE 3 (P3) — MÉDIO (1 SEMANA)</b></font>", table_cell_bold),
            Paragraph(
                "1. <b>Mitigar SSRF em Probe:</b> Criar allowlist para IP/rede local ou restringir chamadas a nós previamente cadastrados.<br/>"
                "2. <b>Padronizar Respostas de Erro JSON:</b> Substituir formatações manuais com <font name='Courier'>fmt.Sprintf</font> por encoders seguros.",
                table_cell_style
            )
        ],
    ]
    t_recom = Table(recom_data, colWidths=[55*mm, 119*mm])
    t_recom.setStyle(TableStyle([
        ('BOX', (0,0), (-1,-1), 0.5, COLOR_BORDER),
        ('INNERGRID', (0,0), (-1,-1), 0.5, COLOR_BORDER),
        ('TOPPADDING', (0,0), (-1,-1), 2.5*mm),
        ('BOTTOMPADDING', (0,0), (-1,-1), 2.5*mm),
        ('LEFTPADDING', (0,0), (-1,-1), 3*mm),
        ('RIGHTPADDING', (0,0), (-1,-1), 3*mm),
        ('BACKGROUND', (0,0), (0,0), colors.HexColor('#FEF2F2')),
        ('BACKGROUND', (0,1), (0,1), colors.HexColor('#FFF7ED')),
        ('BACKGROUND', (0,2), (0,2), colors.HexColor('#FEFCE8')),
    ]))
    story.append(t_recom)

    story.append(PageBreak())

    # =========================================================================
    # 6. SEÇÃO DE ISSUES PARA O GITHUB (MARKDOWN COMPLETO)
    # =========================================================================
    story.append(Paragraph("Issues Prontas para o GitHub (Copiar e Colar)", h1_style))
    story.append(Paragraph(
        "Abaixo estão os templates completos formatados em Markdown para criação imediata no repositório GitHub do projeto:",
        body_style
    ))
    story.append(Spacer(1, 3*mm))

    issues_text = [
        # ISSUE 1
        """--- ISSUE 1 ---
# [Segurança] Vazamento de dados multi-tenant na rota GET /api/v1/users
**Labels:** security, critical, bug, backend

### Descrição do Problema
O handler `HandleUsers` em `internal/adapters/primary/http/user_handler.go` executa uma query SQL sem filtrar pelo `tenant_id` da sessão autenticada. Qualquer usuário autenticado de qualquer inquilino consegue listar todos os usuários, e-mails, nomes completos, papéis e escopos de outras empresas.

### Evidência
- **Arquivo:** `internal/adapters/primary/http/user_handler.go:29-48`
```go
query := `
    SELECT 
        u.id::text, 
        u.email, 
        COALESCE(u.name, 'Administrador'), 
        u.role, 
        u.is_active, 
        COALESCE(t.name, 'Empresa Alfa'), 
        u.created_at, 
        COALESCE(u.last_login_at, u.created_at)
    FROM users u
    LEFT JOIN tenants t ON t.id = u.tenant_id
    ORDER BY u.created_at DESC
`
```

### Impacto
Quebra total do isolamento multi-tenant e violação da LGPD/GDPR, permitindo enumeração de clientes e colaboradores.

### Sugestão de Correção
1. Extrair o `tenantID` do contexto via `middleware.GetTenantID(r.Context())`.
2. Adicionar cláusula `WHERE u.tenant_id = $1` na query SQL.

### Critérios de Aceite
- [ ] Usuário do Tenant A só recebe usuários pertencentes ao Tenant A.
- [ ] Teste unitário e de integração validando isolamento entre inquilinos.
--- FIM ISSUE 1 ---""",

        # ISSUE 2
        """--- ISSUE 2 ---
# [Segurança] Ausência de verificação de papel de Administrador nas rotas sensíveis do Backend
**Labels:** security, critical, rbac, backend

### Descrição do Problema
O frontend protege as páginas de administração verificando `v-if="isAdmin"` no cliente. No entanto, o backend (`router.go` e handlers de usuários, storage, nós e auditoria) não valida se o usuário possui `role == 'admin'`, permitindo que tokens de operadores e visualizadores executem ações de escrita e configurações globais.

### Evidência
- **Arquivos:** `internal/adapters/primary/http/router.go:60-100`, `user_handler.go`, `storage_pool_handler.go`
```go
// router.go aplica apenas AuthMiddleware sem verificar papéis
handler = middleware.AuthMiddleware(validator)(handler)
```

### Impacto
Escalação de privilégios horizontal e vertical: operadores podem excluir pools de storage, cadastrar/deletar nós e acessar logs forenses.

### Sugestão de Correção
Criar o middleware `RequireRole(domain.RoleAdmin)` e aplicá-lo em `/api/v1/users`, `/api/v1/storage/*`, `/api/v1/cluster/*` e `/api/v1/audit/*`.

### Critérios de Aceite
- [ ] Requisições de usuários com `role != 'admin'` retornam `403 Forbidden` nas rotas administrativas.
--- FIM ISSUE 2 ---""",

        # ISSUE 3
        """--- ISSUE 3 ---
# [Segurança] Chave estática padrão para assinatura JWT e ausência de bloqueio no boot
**Labels:** security, critical, auth, backend

### Descrição do Problema
Se a variável de ambiente `JWT_SECRET` não for informada, o serviço utiliza um segredo estático público (`hydravms-super-secret-production-key-2026-auth-token-guard`). Além disso, o startup não valida o comprimento e complexidade da chave.

### Evidência
- **Arquivo:** `internal/application/auth_service.go:30-38`
```go
secret := os.Getenv("JWT_SECRET")
if secret == "" {
    secret = "hydravms-super-secret-production-key-2026-auth-token-guard"
}
```

### Impacto
Qualquer atacante que conheça o código aberto pode forjar tokens JWT arbitrários com qualquer UUID de tenant e papel de administrador.

### Sugestão de Correção
Remover o fallback estático e emitir `log.Fatalf` se `JWT_SECRET` estiver vazio ou possuir menos de 32 bytes.

### Critérios de Aceite
- [ ] Servidor recusa iniciar sem `JWT_SECRET` válido.
- [ ] Tokens forjados com a chave antiga são rejeitados.
--- FIM ISSUE 3 ---""",

        # ISSUE 4
        """--- ISSUE 4 ---
# [Segurança] Arbitrary File Read / Path Traversal em endpoint de streaming de gravações
**Labels:** security, critical, backend, lfi

### Descrição do Problema
O endpoint de stream de vídeo de gravações (`GET /api/v1/cameras/{id}/recordings/stream`) executa `http.ServeFile(w, r, chosen.S3Key)` diretamente no valor gravado em `s3_key` sem sanitizar se o caminho aponta para arquivos do sistema operacional fora do diretório de gravações.

### Evidência
- **Arquivo:** `internal/adapters/primary/http/camera_handler.go:297-308`
```go
if chosen.S3Key == "" {
    writeError(w, http.StatusNotFound, "Recording file path not found")
    return
}
http.ServeFile(w, r, chosen.S3Key)
```

### Impacto
Leitura arbitrária de arquivos confidenciais do servidor (como `/etc/passwd`, chaves SSH, arquivos `.env` e bancos de dados).

### Sugestão de Correção
1. Validar que `s3_key` esteja dentro do caminho raiz permitido usando `filepath.Clean` e `strings.HasPrefix`.
2. Rejeitar caminhos absolutos fora de `storage/` ou servir via MinIO S3 Stream em vez de `http.ServeFile`.

### Critérios de Aceite
- [ ] Tentativas de acessar caminhos fora de `storage/` retornam `403 Forbidden`.
--- FIM ISSUE 4 ---""",

        # ISSUE 5
        """--- ISSUE 5 ---
# [Segurança] IDOR na exclusão de Nós de Cluster e Pools de Storage
**Labels:** security, high, idor, backend

### Descrição do Problema
Os endpoints `DELETE /api/v1/cluster/nodes/{id}` e `DELETE /api/v1/storage/pools/{id}` executam a remoção direta do objeto pelo ID sem verificar se o recurso pertence ao tenant do usuário autenticado.

### Evidência
- **Arquivo:** `internal/adapters/primary/http/cluster_node_handler.go:305-320`
```go
idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/cluster/nodes/")
id, err := uuid.Parse(idStr)
h.repo.Delete(r.Context(), id) // Sem validação de tenant
```

### Impacto
Inquilinos maliciosos podem deletar nós de GPU/Ingestão e pools de armazenamento de outros inquilinos ou do cluster principal.

### Sugestão de Correção
Passar `tenantID` para o método `Delete(ctx, tenantID, id)` e incluir `WHERE id = $1 AND tenant_id = $2` na query SQL.

### Critérios de Aceite
- [ ] Tentativa de excluir nó ou pool de outro tenant retorna `404 Not Found` ou `403 Forbidden`.
--- FIM ISSUE 5 ---"""
    ]

    for issue in issues_text:
        story.append(Paragraph(issue.replace("\n", "<br/>").replace(" ", "&nbsp;"), issue_code_style))
        story.append(Spacer(1, 2*mm))

    # Construir PDF com NumberedCanvas
    doc.build(story, canvasmaker=NumberedCanvas)
    print(f"✅ PDF gerado com sucesso em: {pdf_path}")

if __name__ == "__main__":
    build_pdf()
