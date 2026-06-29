package main

// i18n.go: traduções do goflix
//
// Para adicionar um novo idioma:
//  1. Copie um bloco existente com uma nova chave (ex: "fr", "en")
//  2. Traduza os valores, não altere as chaves
//  3. Compile localmente:  go build -o goflix .
//     ou publique um PR e será incluído na próxima release
//
// Alternância de idioma: tecla  l  na tela inicial do app
// Idioma padrão: ptbr

var langs = map[string]map[string]string{
	"ptbr": {
		// Tela inicial
		"no_dirs":  "Nenhum diretório mapeado ainda.",
		"hint_add": "Pressione  n  para adicionar um diretório.",

		// Footers
		"footer_main":   "enter: selecionar  |  n: novo diretório  |  d: remover  |  l: idioma  |  q: sair",
		"footer_files_dirs": "enter: abrir  |  a: renomear  |  esc: voltar  |  q: sair",
		"footer_files":  "enter: abrir  |  v: marcar visto  |  r: remover visto  |  a: renomear  |  esc: voltar  |  q: sair",
		"footer_input":  "enter: confirmar  |  esc: cancelar",
		"footer_rename": "enter: confirmar  |  esc: cancelar",

		// Diretórios
		"dir_added":   "Diretório adicionado",
		"dir_exists":  "Diretório já existe na lista.",
		"dir_invalid": "Caminho inválido ou não é um diretório.",
		"dir_removed": "Diretório removido",
		"prompt_dir":  "Caminho do diretório:",

		// Vídeos
		"no_video":         "Nenhum vídeo encontrado.",
		"dur_unknown":      "?",
		"progress":         "de",
		"continue":         "Continuar",
		"watched_label":    "Visto",
		"marked_watched":   "Marcado como visto",
		"unmarked_watched": "Removido de visto",

		// Renomeação
		"rename_label": "Renomear",
		"renamed":      "Renomeado",

		// Atualização
		"update_available": "Nova versão disponível",
		"update_current":   "atual",
		"update_prompt":    "u: atualizar agora    qualquer tecla: ignorar",
		"update_updating":  "Atualizando...",
		"update_error":     "Erro ao atualizar.",

		// Misc
		"loading":      "Carregando...",
		"bye":          "até logo.",
		"lang_changed": "Idioma alterado →",
	},

	"es": {
		// Pantalla inicial
		"no_dirs":  "Ningún directorio mapeado aún.",
		"hint_add": "Presiona  n  para agregar un directorio.",

		// Footers
		"footer_main":   "enter: seleccionar  |  n: nuevo directorio  |  d: eliminar  |  l: idioma  |  q: salir",
		"footer_files_dirs": "enter: abrir  |  a: renombrar  |  esc: volver  |  q: salir",
		"footer_files":  "enter: abrir  |  v: marcar visto  |  r: desmarcar visto  |  a: renombrar  |  esc: volver  |  q: salir",
		"footer_input":  "enter: confirmar  |  esc: cancelar",
		"footer_rename": "enter: confirmar  |  esc: cancelar",

		// Directorios
		"dir_added":   "Directorio agregado",
		"dir_exists":  "El directorio ya está en la lista.",
		"dir_invalid": "Ruta inválida o no es un directorio.",
		"dir_removed": "Directorio eliminado",
		"prompt_dir":  "Ruta del directorio:",

		// Videos
		"no_video":         "No se encontraron videos.",
		"dur_unknown":      "?",
		"progress":         "de",
		"continue":         "Continuar",
		"watched_label":    "Visto",
		"marked_watched":   "Marcado como visto",
		"unmarked_watched": "Desmarcado de visto",

		// Renombrado
		"rename_label": "Renombrar",
		"renamed":      "Renombrado",

		// Actualización
		"update_available": "Nueva versión disponible",
		"update_current":   "actual",
		"update_prompt":    "u: actualizar ahora    cualquier tecla: ignorar",
		"update_updating":  "Actualizando...",
		"update_error":     "Error al actualizar.",

		// Misc
		"loading":      "Cargando...",
		"bye":          "hasta luego.",
		"lang_changed": "Idioma cambiado →",
	},
}

var currentLang = "ptbr"

// langLabel retorna o nome legível do idioma atual para exibição no flash
var langLabel = map[string]string{
	"ptbr": "Português",
	"es":   "Español",
}

// t retorna a string traduzida para o idioma atual
// Faz fallback para ptbr se a chave não existir no idioma escolhido
func t(key string) string {
	if v, ok := langs[currentLang][key]; ok {
		return v
	}
	if v, ok := langs["ptbr"][key]; ok {
		return v
	}
	return "[" + key + "]"
}

// toggleLang alterna entre os idiomas disponíveis em ordem alfabética
func toggleLang() {
	keys := make([]string, 0, len(langs))
	for k := range langs {
		keys = append(keys, k)
	}
	for i := 0; i < len(keys)-1; i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	for i, k := range keys {
		if k == currentLang {
			currentLang = keys[(i+1)%len(keys)]
			return
		}
	}
}
