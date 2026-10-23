package approval

import (
	"path"
	"regexp"
	"sort"
)

// ManifestEntry es una entrada determinista del manifiesto: ruta relativa al
// repo y hash de git del blob en P.
type ManifestEntry struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

// informeLinkPattern encuentra links markdown a informes: "(informes/...md)".
// Son los únicos archivos, además de la propuesta, que el contrato exige
// versionados e incluidos en el manifiesto — el resto del cuerpo de la
// propuesta no se audita byte a byte.
var informeLinkPattern = regexp.MustCompile(`\(((?:\.\./)*informes/[^)\s]+\.md)\)`)

// referencedInformes devuelve, en orden determinista, las rutas de informe
// que content referencia, resueltas relativas al directorio de proposalPath.
func referencedInformes(content, proposalPath string) []string {
	base := path.Dir(proposalPath)
	seen := map[string]bool{}
	var out []string
	for _, m := range informeLinkPattern.FindAllStringSubmatch(content, -1) {
		rel := m[1]
		full := path.Clean(path.Join(base, rel))
		if !seen[full] {
			seen[full] = true
			out = append(out, full)
		}
	}
	sort.Strings(out)
	return out
}

// buildManifest construye el manifiesto determinista de la propuesta y sus
// informes referenciados, hasheados en la revisión p. Devuelve el path que
// faltó como missing si algún informe referenciado no existe en p — el
// manifiesto no se completa a medias.
func buildManifest(dir, p, proposalPath string, informes []string) (entries []ManifestEntry, missing string, err error) {
	paths := append([]string{proposalPath}, informes...)
	sort.Strings(paths)
	for _, pth := range paths {
		hash, ok, errHash := blobHash(dir, p, pth)
		if errHash != nil {
			return nil, "", errHash
		}
		if !ok {
			return nil, pth, nil
		}
		entries = append(entries, ManifestEntry{Path: pth, Hash: hash})
	}
	return entries, "", nil
}
