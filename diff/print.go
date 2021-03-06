package diff

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

// PrintMultiFileDiff prints a multi-file diff in unified diff format.
func PrintMultiFileDiff(ds []*FileDiff) ([]byte, error) {
	var buf bytes.Buffer
	for _, d := range ds {
		diff, err := PrintFileDiff(d)
		if err != nil {
			return nil, err
		}
		if _, err := buf.Write(diff); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// PrintFileDiff prints a FileDiff in unified diff format.
//
// TODO(sqs): handle escaping whitespace/etc. chars in filenames
func PrintFileDiff(d *FileDiff) ([]byte, error) {
	var buf bytes.Buffer

	for _, xheader := range d.Extended {
		if _, err := fmt.Fprintln(&buf, xheader); err != nil {
			return nil, err
		}
	}

	if err := printFileHeader(&buf, "--- ", d.OrigName, d.OrigTime); err != nil {
		return nil, err
	}
	if err := printFileHeader(&buf, "+++ ", d.NewName, d.NewTime); err != nil {
		return nil, err
	}

	ph, err := PrintHunks(d.Hunks)
	if err != nil {
		return nil, err
	}

	if _, err := buf.Write(ph); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func printFileHeader(w io.Writer, prefix string, filename string, timestamp *time.Time) error {
	if _, err := fmt.Fprint(w, prefix, filename); err != nil {
		return err
	}
	if timestamp != nil {
		if _, err := fmt.Fprint(w, "\t", timestamp.Format(diffTimeFormat)); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	return nil
}

// PrintHunks prints diff hunks in unified diff format.
func PrintHunks(hunks []*Hunk) ([]byte, error) {
	var buf bytes.Buffer
	for _, hunk := range hunks {
		_, err := fmt.Fprintf(&buf,
			"@@ -%d,%d +%d,%d @@", hunk.OrigStartLine, hunk.OrigLines, hunk.NewStartLine, hunk.NewLines,
		)
		if err != nil {
			return nil, err
		}
		if hunk.Section != "" {
			_, err := fmt.Fprint(&buf, " ", hunk.Section)
			if err != nil {
				return nil, err
			}
		}
		if _, err := fmt.Fprintln(&buf); err != nil {
			return nil, err
		}
		if len(hunk.Body) != 0 && hunk.Body[len(hunk.Body)-1] != '\n' {
			// Append message if hunk.Body doesn't end with a newline.
			hunk.Body = append(hunk.Body, noNewlineMessage...)
		}
		if _, err := buf.Write(hunk.Body); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
