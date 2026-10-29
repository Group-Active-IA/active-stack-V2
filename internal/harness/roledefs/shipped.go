package roledefs

// shippedRoleHashes lists, per resource filename, the sha256 (hex) of every
// version of that role asset that active-stack already PUBLISHED (tarea 5.2
// de lean-agent-workflow). A file on disk byte-identical to one of these is
// ours even without an ownership manifest — a machine that installed before
// the manifest existed, or lost it — so Install treats it as "propio
// intacto": backup, then update. Anything else stays a conflict.
//
// Lista GENERADA, no leída de git en runtime: el binario no tiene por qué
// correr dentro de un checkout. Se genera recorriendo
//
//	git log --all --format=%H -- internal/harness/roledefs/assets
//
// y, por cada commit, el sha256 de cada assets/role-*.md
// (git show <commit>:internal/harness/roledefs/assets/<name> | sha256sum).
// TestShippedRoleHashes_MatchGitHistory recalcula lo mismo y falla si esta
// lista no coincide exactamente, así que un commit nuevo a los assets obliga
// a regenerarla en el mismo cambio.
//
// Commits fuente al regenerarla (2026-10-29): 2526eee, 5ce2c0f y el commit
// de la rama stack-actualizado que publica los cinco roles de flujo. Un
// role-auditor.md preexistente que no coincide con esa versión sigue siendo
// personal hasta que se demuestre lo contrario.
var shippedRoleHashes = map[string][]string{
	"role-analyst.md": {
		"47db39ee091573c80499b6b00d995437dc5eb78112e2449bdcb4f20ab2d99389", // 2526eee
		"7eea4cec3f6579b4bc0a572622efede34e02b58f9ad90c40503010b7146df159", // 5ce2c0f
		"68f1cb1b4a2d38a8290afecb4878edea7de9f02b91dc9ec5ba117b7483c87c49", // fab5966
	},
	"role-architect.md": {
		"a49c5bdd0360c869befadc7f6a9487c773575a2b08ffa6234deb3b297df79f1d", // 2526eee
		"7fae1cba18e19be4096bb334ab27b84aa394b7727256e28461c327de4dfa0778", // 5ce2c0f
		"1082a6ea7db0cbbad7efe5401c6ba674928b2112af21e1369da6076935bc5670", // fab5966
	},
	"role-auditor.md": {
		"3eabc55274cc43f53a0bd905f4adb2e6d055cd8b836040c28b550de34ceea1a4", // fab5966
	},
	"role-backend.md": {
		"80a7ad53cb727f68951735e4175179287d2f75ddd2c5b8e0c2f52da4eb28f9a1", // 2526eee
		"c3aea771ee24b5f4797f90332b9b2fac5e6c01474ee66d0d1330ab00b939e5de", // 5ce2c0f
		"fd965d74bb851e43bc12a75348cd1ff8a4541cd496c21c39b28e31bcf5d67353", // fab5966
	},
	"role-breaker.md": {
		"081f137064b09c1c122d112d3aea7e0c9abeefa2caa05d6a786e39320686c2b9", // fab5966
		"336ec5cba7898f2b388bfa8c1c8e5f00a4a2a1bb509d2943cb87f3512093e84c", // stack-actualizado: tests enfocados
	},
	"role-constructor.md": {
		"ff1026a02c5badc0a0a63630dbc1fcedf466f8dd196fe2348cd582250c08a368", // fab5966
		"fee534fe45ceb91075a363ee85b4b60a6eef980ea6368e263659a5516ad3bd20", // stack-actualizado: tests enfocados
	},
	"role-dba.md": {
		"9cadb7eb8fc1eed1dce6c8a04543978959a53589ba328be130663c3bac0b0f21", // 2526eee
		"def9b40b4f715ccf6882a569347bf6b938179ef5ab1268f9c3c55fbe6238cece", // 5ce2c0f
		"c2f5d002ba3653de83a85116e57e64cdfd083386fc336fb9e5cd28120f81b299", // fab5966
	},
	"role-devops.md": {
		"1067c9bb5c0c752f8e5df4347e33fa7e16c9afb7088815d161f3352b0ef17325", // 2526eee
		"2bf177695803482371cbf14e075ca4467438b986fdd14206ec663f8e0148a576", // 5ce2c0f
		"0521cbbddb38ac809b4c0f0123b8c4b292a48435ac2adcb8bc2bf18ab30f4971", // fab5966
	},
	"role-frontend.md": {
		"86a499ad1af5befed2f322f47435b8024d1aef10e9dc6cb0a869158b2da6f5e9", // 2526eee
		"27e393eb16a2deeebcfa7f7582250526b96a41f945a2bd3b524752164d3a6d85", // 5ce2c0f
		"60a18c25fcb0099afa5c4a08a5fb7461a87e6bf0ccd1c25874fbc13c7603ac5d", // fab5966
	},
	"role-plan-reviewer.md": {
		"5eb5b230a9efb446192e99c9520b890e993f96936578ee47a746e3fcdaf8f544", // fab5966
	},
	"role-planner.md": {
		"b16a407f130e8cebaaf872097e9137d4cc6a50752ce23fb8c659ee341695695d", // fab5966
	},
	"role-qa.md": {
		"a8d339300028d8633a0b201f2f8507e085b94054571382a2e7d57cec53c4ec0a", // 2526eee
		"67cc3d878ab0fcbfcb38ec51e69726c8bcb7f320cc2d78336ffbde1f9d027af7", // 5ce2c0f
		"8fe721ee7d68c1c71ed8ea15da6af9cba72b78a57d17429e18784596074d2f46", // fab5966
	},
}

// isShippedVersion reports whether content hashing to hash is byte-identical
// to a previously published version of resource name. The hash only counts
// for its own name: a published role-qa.md copied over role-dba.md is not
// ours to overwrite.
func isShippedVersion(name, hash string) bool {
	for _, h := range shippedRoleHashes[name] {
		if h == hash {
			return true
		}
	}
	return false
}
