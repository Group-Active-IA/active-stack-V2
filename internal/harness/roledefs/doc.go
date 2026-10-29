// Package roledefs installs twelve provider-neutral role definitions: five
// sequential workflow responsibilities and seven optional specialist criteria.
// Specialists retain the identifiers from internal/roles for compatibility;
// installing their definitions does not make them mandatory workflow attendees.
//
// Claude resolves flat role-*.md resources through its adapter AgentsDir.
// Adapters without native subagents return an empty directory and are skipped;
// the compact orchestrator contract supplies universal role guidance separately.
//
// Installation and removal preserve personal resources through ownership hashes,
// conflict handling, backups, atomic writes and rollback of partial adapter work.
// The manifest records installed content; no model or provider is selected here.
// A file byte-identical to a previously published asset (hashes embedded in
// shipped.go, checked against git history by tests) is also treated as ours.
// A conflict on a workflow role blocks the lean flow: install reports it and
// Health flags its content as "foreign"; a personal specialist is legitimate.
package roledefs
