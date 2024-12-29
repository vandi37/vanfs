package ide

import (
	"io"
)

func (m *model) loadFile() error {
	content, err := io.ReadAll(m.file)
	if err != nil {
		return err
	}
	m.textarea.InsertString(string(content))

	return nil
}

func (m *model) saveFile() error {
	_, err := m.file.WriteAt([]byte(m.textarea.Value()), 0)
	return err
}
