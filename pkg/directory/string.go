package directory

import "strings"

func (d *Directory) String() string {
	return ".\n" + d.string("", true)
}

func (d *Directory) string(prefix string, isLast bool) string {
	var sb strings.Builder
	count := len(d.dirs)

	var filePrefix = "│   "
	if count == 0 {
		filePrefix = "   "
	}

	for filename := range d.files {
		sb.WriteString(prefix + filePrefix + filename + "\n")
	}

	if count == 0 {
		return sb.String()
	}

	marker := func() string {
		if isLast {
			return "└─"
		}
		return "├─"
	}

	i := 0

	for dirName, dir := range d.dirs {
		isLast = i == count-1
		i++
		if dir == nil {
			continue
		}

		sb.WriteString(prefix + marker() + dirName + "\n")

		newPrefix := prefix
		if isLast {
			newPrefix += "   "
		} else {
			newPrefix += "│  "
		}
		sb.WriteString(dir.string(newPrefix, isLast))
	}

	return sb.String()
}
