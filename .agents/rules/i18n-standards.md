# Regras de Internacionalização (i18n Standards)

## 1. Paridade Obrigatória de Idiomas (PT, EN e ES)
Toda e qualquer nova chave de texto, rótulo de tela, mensagem de erro, placeholder ou tooltip adicionada no frontend DEVE obrigatoriamente ser registrada e traduzida simultaneamente em todos os arquivos de idioma:
- `web/src/locales/pt.ts` (Português)
- `web/src/locales/en.ts` (Inglês)
- `web/src/locales/es.ts` (Espanhol)

## 2. Proibição de Strings Hardcoded
É expressamente proibido inserir textos estáticos em português ou inglês diretamente dentro do template Vue (`<template>`). Todo texto exibido ao usuário deve usar o composable `useI18n()` e a função `t('chave_de_texto')`.

## 3. Verificação Automatizada
Após qualquer adição ou modificação no frontend, a paridade de chaves deve ser validada com:
```bash
npm --prefix web run check:i18n
```
Se houver qualquer discrepância ou chave ausente em qualquer idioma, a build ou commit deve ser bloqueado imediatamente até a correção.
