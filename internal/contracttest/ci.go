package contracttest

var CI = []string{
	`(pending|queued|in_progress|cancelled|skipped|omitted|neutral|timed_out|another sha|different sha|zero runs|no runs|failed query|query failure)[^.\n]{0,75}(is|are|counts as|considered|treat as|means) (green|success)\b`,
	`(green|success)\b[^.\n]{0,45}(despite|with|even with)[^.\n]{0,35}(pending|queued|skipped|omitted|another sha|different sha|zero runs|no runs|failed query|query failure)`,
}
