package api

// TaskDescriptors for task management by API
var TaskDescriptors = Module{
	Name: "task",
	Protos: map[string]Proto{
		"ShowAllTasks": {
			Name:    "查看所有任务",
			handler: ShowAllTasks,
			Paras:   []ProtoPara{},
		},
	},
}
