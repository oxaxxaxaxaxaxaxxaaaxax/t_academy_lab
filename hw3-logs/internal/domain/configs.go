package domain

type Config struct {
	PathToFile string
	Format     string
	Output     string
	From       string
	To         string
}

type LogCtx struct {
	FilePath []string
}
