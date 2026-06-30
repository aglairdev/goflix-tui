<div align="center">

[![GOFLIX](https://img.shields.io/badge/GOFLIX-C55A10?style=for-the-badge)](https://github.com/aglairdev/Goflix)

</div>

## Que isso?

Gerenciador de vídeos no terminal.

![Go](https://img.shields.io/badge/Go-333333?style=flat-square&logo=go&logoColor=white)

<div align="center">

| Goflix |
|:------:|
| <img src="https://github.com/user-attachments/assets/47f7fb85-d3f4-4d7b-8b25-7f27dc94f5f6" alt="Goflix Demo" width="600"> |

</div>

## Instalação

**Go (recomendado):**
```bash
go install github.com/aglairdev/goflix@latest
```
**Versão específica:**
```bash
go install github.com/aglairdev/goflix@v1.0.0
```
Requer: `go` `mpv` 

**Via release** (sem Go instalado):
Baixe o binário em [releases](https://github.com/aglairdev/goflix/releases)

Requer: `mpv` 

mova para `~/.local/bin` e dê permissão de execução:

```bash
cd go/bin/
mv goflix ~/.local/bin/
```
Certifique-se de que `~/.local/bin` está no seu PATH. Se não estiver, adicione ao seu `~/.bashrc` ou `~/.zshrc`:

```bash
export PATH="$HOME/.local/bin:$PATH"
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
| `-d` | modo debug (verbose no stderr) |
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
ꕤ Nova versão disponível: v1.0.1  (atual: v1.0.0)
─────────────────────────────────────────────────
u: atualizar agora    qualquer tecla: ignorar
```

Pressione `u` para atualizar automaticamente. Requer Go instalado.

Usuários que instalaram via release receberão o aviso, mas precisarão baixar o novo binário manualmente em [releases](https://github.com/aglairdev/goflix/releases).

## Dados

| Arquivo | Conteúdo |
|---|---|
| `~/.config/goflix/config` | diretórios mapeados |
| `~/.config/goflix/watched` | histórico de assistidos |
| `~/.config/goflix/settings` | preferências (tema) |

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
rm -r ~/.config/goflix
```

## Licença
[MIT](https://github.com/aglairdev/Goflix/blob/main/LICENSE)

<p align="center">ꕤ AGL</p>
