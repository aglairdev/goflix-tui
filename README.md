# goflix ꕤ

Acervo de vídeos no terminal.

## Demo

<img width="1919" height="1079" alt="Demo" src="https://github.com/user-attachments/assets/db0e9c06-28dd-4de8-937f-507a5f1ce825" />


## Instalação

**Go (recomendado):**
```bash
go install github.com/aglairdev/goflix@latest
```
Requer: `go` `mpv` `ffprobe` (opcional)

**Versão específica:**
```bash
go install github.com/aglairdev/goflix@v1.0.0
```
Requer: `go` `mpv` `ffprobe` (opcional)

**Via release** (sem Go instalado):
Baixe o binário em [Releases](https://github.com/aglairdev/goflix/releases), mova para `~/.local/bin` e dê permissão de execução:
```bash
chmod +x goflix
mv goflix ~/.local/bin/
```
Requer: `mpv` `ffprobe` (opcional)

Certifique-se de que `~/.local/bin` está no seu PATH. Se não estiver, adicione ao seu `~/.bashrc` ou `~/.zshrc`:
```bash
export PATH="$HOME/.local/bin:$PATH"
```

## Uso

```bash
goflix
```

### Atalhos

**Tela inicial:**

| Tecla | Ação |
|---|---|
| `enter` | abrir diretório |
| `n` | adicionar diretório |
| `d` | remover diretório |
| `l` | alternar idioma (pt/es) |
| `q` | sair |

**Dentro de um diretório:**

| Tecla | Ação |
|---|---|
| `enter` | reproduzir vídeo |
| `v` | marcar como visto |
| `r` | desmarcar visto + resetar progresso |
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

Usuários que instalaram via release receberão o aviso, mas precisarão baixar o novo binário manualmente em [Releases](https://github.com/aglairdev/goflix/releases).

## Dados

| Arquivo | Conteúdo |
|---|---|
| `~/.config/goflix/config` | diretórios mapeados |
| `~/.config/goflix/watched` | histórico de assistidos |

## Personalização

### Cor 

Edite a constante `ColorAccent` no topo de `main.go` (valor hex):

```go
ColorAccent = "#FF5FA7" // padrão: rosa
ColorAccent = "#00BFFF" // azul
ColorAccent = "#A8FF3E" // verde
```

Qualquer cor hex de 6 dígitos funciona. Recompile após alterar:
```bash
go build -o goflix .
```

### Idioma

O idioma pode ser alternado com `l` diretamente no app, sem reiniciar.

## Tradução

Edite `i18n.go` e adicione um novo bloco com a chave do idioma:
```go
"en": {
    "no_dirs":  "No directories mapped yet.",
    "hint_add": "Press  n  to add a directory.",
    // ... translate all keys using the "ptbr" block as reference
},
```
Adicione também o nome do idioma em `langLabel`:
```go
"en": "English",
```
Use o bloco `"ptbr"` como referência - todas as chaves precisam estar presentes. Abra um PR com o novo idioma.

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
