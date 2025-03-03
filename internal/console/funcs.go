package console

import (
	"fmt"
	"strings"

	"github.com/vandi37/vanfs/pkg/ide"
)

func (c *Console) Help(s string) error {
	fmt.Println(`cd {path}: Change directory
tree {path}: View tree,
ls {path}: View list of files and directories,
mkdir {path}: Make directory,
rmdir {path}: Remove directory,
mkf {path}: Make file,
rm {path}: Remove file,
of {path}: Open file,
path: Current path,
clear: Clear console,
cat {path}: View file value,
help: Help`)
	return nil
}

func (c *Console) Tree(s string) error {
	res, err := c.fs.Tree(s)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func (c *Console) Ls(s string) error {
	res, err := c.fs.Ls(s)
	if err != nil {
		return err
	}
	fmt.Println("\033[38;2;189;38;93m", strings.Join(res, "\t"), "\033[0m")
	return nil
}

func (c *Console) Of(s string) error {
	file, err := c.fs.Of(s)
	if err != nil {
		return err
	}
	defer file.Close()
	err = ide.Run(file)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func (c *Console) Path(s string) error {
	fmt.Printf("\033[38;2;144;238;144m%s\033[0m\n", c.fs.Path)
	return nil
}

func (c *Console) Clear(s string) error {
	c.cleaner.Clear()
	return nil
}

func (c *Console) Cat(s string) error {
	res, err := c.fs.Cat(s)
	if err != nil {
		return err
	}
	fmt.Printf("\033[38;2;144;238;144m%s\033[0m\n", res)
	return nil
}
