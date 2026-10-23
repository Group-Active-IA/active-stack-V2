package contracttest

var RF1 = []string{
	`(always|unconditionally)[^.\n]{0,40}(mark|report|declare)[^.\n]{0,25}auditor[^.\n]{0,25}(absent|missing)`,
	`(for this repository|current repository|this repo)[^.\n]{0,65}(require|requirement)|(apply|impose)[^.\n]{0,55}(requirement|ci)[^.\n]{0,40}current repository`,
	`senior (marks|labels)[^.\n]{0,30}(you|orchestrator)[^.\n]{0,15}(stop|hold|halt)`,
}
