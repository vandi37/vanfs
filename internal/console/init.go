package console

func (c *Console) Init() {
	c.funcs = map[string]func(s string) error{
		"cd":    c.fs.Cd,
		"mkdir": c.fs.Mkdir,
		"tree":  c.Tree,
		"ls":    c.Ls,
		"rmdir": c.fs.Rmdir,
		"mkf":   c.fs.Makefile,
		"rm":    c.fs.Rm,
		"of":    c.Of,
		"path":  c.Path,
		"clear": c.Clear,
		"cat":   c.Cat,
	}
}
