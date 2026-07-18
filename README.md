<div align="center">

<video src="https://github.com/user-attachments/assets/a7fd04a3-8afb-440a-8b58-89a30be09e73" controls width="400"></video>

[![Release](https://img.shields.io/github/v/release/aglairdev/goflix-tui?style=for-the-badge&color=CBA6F7&label=release)](https://github.com/aglairdev/goflix-tui/releases)
![OS](https://img.shields.io/badge/OS-Linux-CBA6F7?style=for-the-badge&logo=linux&logoColor=white)
[![License](https://img.shields.io/github/license/aglairdev/goflix-tui?style=for-the-badge&color=CBA6F7)](LICENSE)
![Go](https://img.shields.io/badge/%3C%2F%3E-Go-CBA6F7?style=for-the-badge)

Gerenciador de mpv TUI.

</div>

## Instalação

**Go (recomendado):**

```bash
go install github.com/aglairdev/goflix@latest
```

Requer: `go` `mpv` 

**Via release**:

Baixe o binário em [releases](https://github.com/aglairdev/goflix-tui/releases).

Requer: `mpv` 

Mova para `~/go/bin` e dê permissão de execução:

```bash
cd ~/go/bin/
sudo chmod +x goflix
```

Certifique-se de que `~/go/bin` está no seu PATH. Se não estiver, adicione ao seu `~/.bashrc` ou `~/.zshrc`:

```bash
export PATH="$HOME/go/bin:$PATH"
```

Ou fish (`~/.config/fish/config.fish`):

```bash
set -Ux PATH $HOME/go/bin $PATH
```

> [!NOTE]
> ffprobe (`ffmpeg`) é necessário apenas se deseja ter a duração dos vídeos no menu.

## Uso

```bash
goflix
```

### Flags

| Flag | Descrição |
|------|-----------|
| `-v` | exibe versão |
| `-d` | modo debug (`~/.config/goflix/debug.log`) |
| `-h` | mostra ajuda |

### Atalhos

**Tela inicial:**

| Tecla | Ação |
|---|---|
| `enter` | abrir diretório |
| `n` | adicionar diretório |
| `d` | remover diretório |
| `t` | alternar tema |
| `l` | alternar idioma (pt/es) |
| `q` | sair |

**Dentro de um diretório:**

| Tecla | Ação |
|---|---|
| `enter` | reproduzir vídeo / abrir subpasta |
| `v` | marcar como visto |
| `r` | desmarcar visto + resetar progresso |
| `a` | renomear arquivo/diretório |
| `esc` | voltar |
| `q` | sair |

## Atualização

O app verifica atualizações ao iniciar. Se houver uma versão nova:

```
ꕤ Nova versão disponível: v1.1.5  (atual: v1.1.4)
─────────────────────────────────────────────────
u: atualizar agora    qualquer tecla: ignorar
```

Pressione `u` para atualizar automaticamente. Requer Go instalado.

> [!NOTE]
> Atualizações automáticas funcionam a partir da versão [1.1.8](https://github.com/aglairdev/goflix-tui/compare/v1.1.7...v1.1.8)

Usuários que instalaram via release receberão o aviso, mas precisarão baixar o novo binário manualmente em [releases](https://github.com/aglairdev/goflix-tui/releases).

## Dados

| Arquivo | Conteúdo |
|---|---|
| `~/.config/goflix/config` | diretórios mapeados |
| `~/.config/goflix/watched` | histórico de assistidos |
| `~/.config/goflix/settings` | preferências (tema) |
| `~/.config/goflix/debug.log` | logs do modo debug |

## Personalização

### Tema

O app tem 5 temas integrados alternados com `t` na tela inicial:

| Tema | Cor |
|------|-----|
| catppuccin (padrão) | `#CBA6F7` |
| cyberpunk | `#00FF9C` |
| gruvbox | `#FE8019` |
| nord | `#88C0D0` |
| netflix | `#E50914` |

A escolha é persistida e restaurada ao iniciar.

### Idioma

O idioma pode ser alternado com `l` diretamente no app, sem reiniciar.
- pt-br/es

## Tradução

Adicione um novo bloco em `i18n.go` copiando o `ptbr` e traduzindo os valores.
Registre o nome em `langLabel`. Abra um PR.

## Build 

```bash
git clone https://github.com/aglairdev/goflix
cd goflix
go build -o goflix .
```

## Remoção

```bash
rm ~/go/bin/goflix
rm -r ~/.config/goflix #histórico de assistidos, diretórios mapeados
```

<p align="center">ꕤ AGL</p>
