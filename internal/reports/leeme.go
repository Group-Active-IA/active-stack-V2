package reports

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// regenerateLEEME reconstruye openspec/changes/<change>/informes/LEEME.md
// leyendo TODOS los manifest.json presentes bajo ese árbol, en vez de
// agregarle una línea cada vez. Regenerar desde disco es la forma más simple
// de no duplicar ni desordenar el índice entre escritores concurrentes de
// change distintos directorios de identidad — no resuelve dos escritores
// tocando el MISMO LEEME.md al mismo tiempo (deuda: sin lock de archivo).
//
// "índice LEEME.md con versión, hashes, tamaños, procedencia y mediación
// automática": el encabezado deja explícito que es generado, no escrito a
// mano, para que nadie lo edite esperando que sobreviva a la próxima
// captura.
func regenerateLEEME(root, change string) error {
	informesDir := filepath.Join(root, "openspec", "changes", change, "informes")

	var manifests []manifest
	var relDirs []string
	err := filepath.WalkDir(informesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if filepath.Clean(path) == filepath.Clean(informesDir) {
				// El propio directorio todavía no existe (primera captura
				// de este change): nada que indexar, no es un error.
				return fs.SkipAll
			}
			return err
		}
		if d.IsDir() || d.Name() != "manifest.json" {
			return nil
		}
		dir := filepath.Dir(path)
		m, ok, errRead := readManifest(dir)
		if errRead != nil {
			return errRead
		}
		if !ok {
			return nil
		}
		rel, errRel := filepath.Rel(informesDir, dir)
		if errRel != nil {
			return errRel
		}
		manifests = append(manifests, m)
		relDirs = append(relDirs, rel)
		return nil
	})
	if err != nil {
		return fmt.Errorf("escanear %s para regenerar LEEME.md: %w", informesDir, err)
	}

	order := make([]int, len(manifests))
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(a, b int) bool { return relDirs[order[a]] < relDirs[order[b]] })

	var b strings.Builder
	fmt.Fprintf(&b, "# Índice: %s\n\n", change)
	b.WriteString("Generado automáticamente por internal/reports — mediación automática, no editar a mano: la próxima captura lo reescribe entero.\n\n")
	if len(order) == 0 {
		b.WriteString("Sin informes capturados todavía.\n")
	}
	for _, idx := range order {
		m := manifests[idx]
		dir := relDirs[idx]
		line := fmt.Sprintf("- `%s` — %s, versión %d, hash `%s`, %d bytes", dir, m.Status.String(), m.Version, m.ContentHash, m.SizeBytes)
		if m.DispatchID != "" {
			line += fmt.Sprintf(", dispatch `%s`", m.DispatchID)
		}
		if m.SHA != "" {
			line += fmt.Sprintf(", SHA `%s`", m.SHA)
		}
		if m.Provider != "" || m.Model != "" {
			line += fmt.Sprintf(", procedencia %s/%s", m.Provider, m.Model)
		}
		if m.Cause != "" {
			line += fmt.Sprintf(" — %s", m.Cause)
		}
		if m.RetainedPath != "" {
			line += fmt.Sprintf(" (original local, ignorado por git: `%s`)", m.RetainedPath)
		}
		b.WriteString(line + "\n")
	}

	return atomicWriteFile(filepath.Join(informesDir, "LEEME.md"), []byte(b.String()), 0o644)
}
