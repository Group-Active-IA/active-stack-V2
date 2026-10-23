package reports

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

func sha256hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// atomicWriteFile escribe a un archivo temporal en el mismo directorio y
// renombra: un crash a mitad de escritura deja el archivo final intacto (o
// ausente), nunca truncado a medias. mismo patrón que requiere la
// recuperación de escrituras parciales (propose §2a, "cortes... y reinicio").
func atomicWriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".tmp-reports-*")
	if err != nil {
		return fmt.Errorf("crear archivo temporal en %s: %w", dir, err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // no-op si el rename ya lo movió
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("escribir %s: %w", tmpName, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("cerrar %s: %w", tmpName, err)
	}
	if err := os.Chmod(tmpName, perm); err != nil {
		return fmt.Errorf("chmod %s: %w", tmpName, err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("renombrar %s a %s: %w", tmpName, path, err)
	}
	return nil
}

// readManifest lee manifest.json en dir. Ausente => (zero, false, nil): no es
// un error, es "nada capturado todavía acá" o "escritura previa interrumpida
// antes de llegar al manifiesto" — en ambos casos, recuperable sin pisar
// nada porque no hay manifiesto confirmado que proteger. Un manifiesto
// presente pero no decodificable (JSON corrupto) se trata igual: recuperable,
// no es un fallo fatal. Un error de I/O real (permiso denegado) SÍ se
// propaga: ahí no sabemos si hay algo que proteger o no, y asumir que no lo
// hay sería arriesgar un pisado silencioso.
func readManifest(dir string) (manifest, bool, error) {
	b, err := os.ReadFile(filepath.Join(dir, "manifest.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return manifest{}, false, nil
		}
		return manifest{}, false, fmt.Errorf("leer manifest.json en %s: %w", dir, err)
	}
	m, err := decodeManifest(b)
	if err != nil {
		return manifest{}, false, nil
	}
	return m, true, nil
}

// writeFiles persiste m en dir. Orden fijo: contenido (raw.txt + informe.md,
// o el puntero retenido) primero, manifest.json AL FINAL — así "existe un
// manifest.json legible" implica "la escritura anterior terminó", que es
// exactamente lo que readManifest necesita para distinguir "nada capturado
// todavía" de "capturado y confirmado".
func writeFiles(root, dir string, m manifest, rawContent, text string) (manifest, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return m, fmt.Errorf("crear %s: %w", dir, err)
	}

	if m.Status == StatusRetained {
		retDir := filepath.Join(dir, ".retenidos")
		if err := os.MkdirAll(retDir, 0o700); err != nil {
			return m, fmt.Errorf("crear %s: %w", retDir, err)
		}
		// La carpeta se ignora sola: el repo del usuario no tiene por qué
		// tener la regla en su .gitignore, y sin esto `git add -A` sube el
		// original con el secreto adentro. Va ANTES del original.
		if err := atomicWriteFile(filepath.Join(retDir, ".gitignore"), []byte("*\n"), 0o644); err != nil {
			return m, err
		}
		retPath := filepath.Join(retDir, m.ContentHash[:12]+".txt")
		if err := atomicWriteFile(retPath, []byte(rawContent), 0o600); err != nil {
			return m, err
		}
		rel, err := filepath.Rel(root, retPath)
		if err != nil {
			return m, fmt.Errorf("ruta relativa de %s: %w", retPath, err)
		}
		m.RetainedPath = rel
	} else {
		rawPath := filepath.Join(dir, "raw.txt")
		if err := atomicWriteFile(rawPath, []byte(rawContent), 0o644); err != nil {
			return m, err
		}
		if err := atomicWriteFile(filepath.Join(dir, "informe.md"), []byte(text), 0o644); err != nil {
			return m, err
		}
		// Contraste byte/hash tras persistir (propose §2a): releer lo que
		// quedó en disco y compararlo contra el hash declarado, en vez de
		// confiar en que la escritura hizo lo que pedimos.
		back, err := os.ReadFile(rawPath)
		if err != nil {
			return m, fmt.Errorf("releer %s para contraste: %w", rawPath, err)
		}
		if sha256hex(string(back)) != m.RawHash {
			return m, fmt.Errorf("fallo de integridad: %s no coincide con el hash declarado tras persistir", rawPath)
		}
	}

	enc, err := m.encode()
	if err != nil {
		return m, fmt.Errorf("codificar manifest.json: %w", err)
	}
	if err := atomicWriteFile(filepath.Join(dir, "manifest.json"), enc, 0o644); err != nil {
		return m, err
	}
	return m, nil
}
